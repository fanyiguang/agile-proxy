package http

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/server/base"
	"agile-proxy/modules/transport"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

type Http struct {
	base.Server
	basicToken string
}

func (h *Http) Run() (err error) {
	h.init()
	err = h.listen()
	return
}

func (h *Http) Close() (err error) {
	if h.Listen != nil {
		err = h.Listen.Close()
	}
	return
}

func (h *Http) listen() (err error) {
	server := &http.Server{
		Addr: net.JoinHostPort(h.Host, h.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.proxy(w, r)
		}),
	}

	h.Listen, err = net.Listen("tcp", server.Addr)
	if err != nil {
		err = errors.Wrap(err, "net.Listen")
		return
	}

	errCh := make(chan error)
	Go.Go(func() {
		err := server.Serve(h.Listen)
		if err != nil {
			select {
			case errCh <- err:
			case <-time.After(time.Second * 5):
				log.WarnF("server: %v server.ListenAndServe failed-1: %v", h.Name(), err)
			}
		}
	})

	select {
	case err = <-errCh:
		log.WarnF("server: %v server.ListenAndServe failed-2: %v", h.Name(), err)
	case <-time.After(time.Second * 5):
		log.InfoF("server: %v init successful, listen: %v", h.Name(), server.Addr)
	}

	return
}

func (h *Http) proxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		err := h.handleConnect(w, r)
		if err != nil {
			log.WarnF("%+v", err)
		}
	} else {
		err := h.handleNormal(w, r)
		if err != nil {
			log.WarnF("%+v", err)
		}
	}
}

func (h *Http) handleConnect(w http.ResponseWriter, r *http.Request) (err error) {
	err = h.authentication(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		err = errors.New("Hijacking not supported")
		return
	}

	conn, _, _err := hijacker.Hijack()
	if _err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		err = errors.Wrap(_err, "hijacker.Hijack failed")
		return
	}

	host, port, err := h.GetHostAndPort(r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return h.transport(conn, host, port)
}

func (h *Http) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if h.Transmitter != nil {
		err = h.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (h *Http) authentication(r *http.Request) (err error) {
	if h.Username == "" && h.Password == "" { // no auth
		return
	}

	token := r.Header.Get("Proxy-Authorization")
	if token == "" {
		err = errors.New("Proxy-Authorization is nil")
		return
	}

	// 直接将请求的token和事先用正确账户密码加密好的字符串对比
	// 提高效率，不用每次都解密请求的token
	if h.basicToken != token {
		err = errors.New("username or password Incorrect")
	}

	return
}

func (h *Http) handleNormal(w http.ResponseWriter, r *http.Request) (err error) {
	http.Error(w, "normal proxy not supported", http.StatusServiceUnavailable)
	log.Warn("normal proxy not supported")
	return
}

func (h *Http) GetHostAndPort(host string) (newHost, newPort []byte, err error) {
	var _host, _port string
	index := strings.Index(host, ":")
	if index == -1 {
		_host = host
		_port = "80"
	} else {
		_host = host[:index]
		_port = host[index+1:]
	}

	newHost = common.StrToBytes(_host)
	newPort = common.StrToBytes(_port)
	return
}

func (h *Http) init() {
	if h.Username != "" && h.Password != "" {
		h.basicToken = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(h.Username+":"+h.Password)))
	}
}

func New(jsonConfig json.RawMessage) (obj *Http, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Http{
		Server: base.Server{
			Net: plugin.Net{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			DoneCh: make(chan struct{}),
		},
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
