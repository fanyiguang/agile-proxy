package https

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/server/base"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

type https struct {
	base.Server
	assembly.Tls
	basicToken string
	//readTimeout  int
	//writeTimeout int
}

func (h *https) Run() (err error) {
	h.init()
	err = h.listen()
	return
}

func (h *https) Close() (err error) {
	if h.Listen != nil {
		err = h.Listen.Close()
	}
	return
}

func (h *https) listen() (err error) {
	tlsConfig, err := h.CreateServerTlsConfig()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:      net.JoinHostPort(h.Host, h.Port),
		TLSConfig: tlsConfig,
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
		err := server.ServeTLS(h.Listen, "", "")
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

func (h *https) proxy(w http.ResponseWriter, r *http.Request) {
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

func (h *https) handleConnect(w http.ResponseWriter, r *http.Request) (err error) {
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

func (h *https) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if h.Route != nil {
		err = h.Route.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Route is nil")
	}
	return
}

func (h *https) authentication(r *http.Request) (err error) {
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

func (h *https) handleNormal(w http.ResponseWriter, r *http.Request) (err error) {
	err = h.authentication(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if connection := r.Header.Get("Proxy-Connection"); connection != "" {
		r.Header.Set("Connection", connection)
		r.Header.Del("Proxy-Connection")
	}
	r.Header.Del("Proxy-Authorization")

	err = h.normalTransport(w, r)
	if err != nil {
		http.Error(w, "system failed", http.StatusInternalServerError)
	}
	return
}

func (h *https) normalTransport(w http.ResponseWriter, r *http.Request) (err error) {
	if h.Route != nil {
		err = h.Route.HttpTransport(w, r)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (h *https) GetHostAndPort(host string) (newHost, newPort []byte, err error) {
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

func (h *https) init() {
	h.Server.Init()
	if h.Username != "" && h.Password != "" {
		h.basicToken = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(h.Username+":"+h.Password)))
	}
}

func New(jsonConfig json.RawMessage) (obj *https, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &https{
		Server: base.Server{
			Net:        assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			DoneCh:     make(chan struct{}),
			RouteName:  config.RouteName,
			Satellites: config.Satellites,
		},
		Tls: assembly.CreateTls(config.CrtPath, config.KeyPath, "", ""),
	}

	return
}
