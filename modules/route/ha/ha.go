package ha

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client"
	"agile-proxy/modules/route/base"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

type ha struct {
	baseTransport base.Transport
	clients       []client.Client
	clientNames   string
}

func (h *ha) Transport(cConn net.Conn, host, port []byte) (err error) {
	if h.clients != nil {
		host, err = h.baseTransport.Dns.GetHost(host)
		if err != nil {
			return
		}

		sConn, err := h.getConn(host, port)
		if err != nil {
			return err
		}

		defer sConn.Close()
		h.baseTransport.AsyncSendMsg(h.baseTransport.Name(), -1, fmt.Sprintf("%v handshark success", common.BytesToStr(host)))
		h.baseTransport.Copy(sConn, cConn)
	} else {
		err = errors.New("client is nil")
	}
	return
}

func (h *ha) HttpTransport(w http.ResponseWriter, r *http.Request) (err error) {
	var newHost []byte
	newHost, err = h.baseTransport.Dns.GetHost(common.StrToBytes(r.Host))
	if err != nil {
		return
	}

	r.Host = common.BytesToStr(newHost)
	resp, err := h.getResp(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	h.baseTransport.HttpCopy(w, resp)
	return
}

func (h *ha) Run() (err error) {
	h.init()
	return
}

func (h *ha) Close() (err error) {
	return
}

func (h *ha) getConn(host, port []byte) (conn net.Conn, err error) {
	connCh := make(chan net.Conn)
	for _, c := range h.clients {
		_client := c
		Go.Go(func() {
			var sConn net.Conn
			sConn, err = _client.Dial("tcp", host, port)
			if err != nil {
				log.DebugF("ha client dial failed: %v %v %v %v", err, host, port, _client.Name())
				return
			}

			select {
			case connCh <- sConn:
			default:
				_ = sConn.Close()
			}
		})
	}

	select {
	case conn = <-connCh:
	case <-time.After(time.Second * 15):
		err = errors.New("get conn timeout")
	}
	return
}

func (h *ha) getResp(r *http.Request) (resp *http.Response, err error) {
	respCh := make(chan *http.Response)
	for _, c := range h.clients {
		_client := c
		Go.Go(func() {
			var _resp *http.Response
			_resp, err = _client.GetRoundTripper().RoundTrip(r)
			if err != nil {
				log.DebugF("ha client roundTrip failed: %v %v %v", err, r.Host, _client.Name())
				return
			}

			select {
			case respCh <- _resp:
			default:
				_ = _resp.Body.Close()
			}
		})
	}

	select {
	case resp = <-respCh:
	case <-time.After(time.Second * 15):
		err = errors.New("get resp timeout")
	}
	return
}

func (h *ha) init() {
	h.baseTransport.Init()
	if h.clientNames != "" {
		clientNames := strings.Split(h.clientNames, ",")
		for _, clientName := range clientNames {
			_client := client.GetClient(clientName)
			if _client != nil {
				h.clients = append(h.clients, _client)
			}
		}
	}
}

func New(jsonConfig json.RawMessage) (obj *ha, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		marshalJSON, _ := jsonConfig.MarshalJSON()
		err = errors.Wrap(err, common.BytesToStr(marshalJSON))
		return
	}

	obj = &ha{
		baseTransport: base.Transport{
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			Dns:        assembly.CreateDns(config.DnsInfo.Server, config.DnsInfo.LocalDns),
			BufferPool: common.CreateByteBufferSyncPool(1024 * 32),
			Satellites: config.Satellites,
		},
		clientNames: config.ClientNames,
	}

	return
}
