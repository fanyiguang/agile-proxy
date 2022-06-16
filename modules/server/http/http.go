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
	sysHttp "net/http"
	"strings"
	"time"
)

type http struct {
	base.Server
	basicToken string
}

func (h *http) Run() (err error) {
	h.init()
	err = h.listen()
	return
}

func (h *http) Close() (err error) {
	if h.Listen != nil {
		err = h.Listen.Close()
	}
	return
}

func (h *http) listen() (err error) {
	server := &sysHttp.Server{
		Addr: net.JoinHostPort(h.Host, h.Port),
		Handler: sysHttp.HandlerFunc(func(w sysHttp.ResponseWriter, r *sysHttp.Request) {
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

func (h *http) proxy(w sysHttp.ResponseWriter, r *sysHttp.Request) {
	if r.Method == sysHttp.MethodConnect {
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

func (h *http) handleConnect(w sysHttp.ResponseWriter, r *sysHttp.Request) (err error) {
	err = h.authentication(r)
	if err != nil {
		sysHttp.Error(w, err.Error(), sysHttp.StatusUnauthorized)
		return
	}

	w.WriteHeader(sysHttp.StatusOK)
	hijacker, ok := w.(sysHttp.Hijacker)
	if !ok {
		sysHttp.Error(w, "Hijacking not supported", sysHttp.StatusInternalServerError)
		err = errors.New("Hijacking not supported")
		return
	}

	conn, _, _err := hijacker.Hijack()
	if _err != nil {
		sysHttp.Error(w, err.Error(), sysHttp.StatusServiceUnavailable)
		err = errors.Wrap(_err, "hijacker.Hijack failed")
		return
	}

	host, port, err := h.GetHostAndPort(r.Host)
	if err != nil {
		sysHttp.Error(w, err.Error(), sysHttp.StatusInternalServerError)
		return err
	}

	return h.transport(conn, host, port)
}

func (h *http) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if h.Transmitter != nil {
		err = h.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (h *http) authentication(r *sysHttp.Request) (err error) {
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

func (h *http) handleNormal(w sysHttp.ResponseWriter, r *sysHttp.Request) (err error) {
	sysHttp.Error(w, "normal proxy not supported", sysHttp.StatusServiceUnavailable)
	log.Warn("normal proxy not supported")
	return
}

func (h *http) GetHostAndPort(host string) (newHost, newPort []byte, err error) {
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

func (h *http) init() {
	if h.Username != "" && h.Password != "" {
		h.basicToken = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(h.Username+":"+h.Password)))
	}
}

func New(jsonConfig json.RawMessage) (obj *http, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &http{
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
