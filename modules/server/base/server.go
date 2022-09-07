package base

import (
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/route"
	"agile-proxy/modules/satellite"
	"net"
)

type Server struct {
	assembly.Net
	assembly.Identity
	assembly.Pipeline
	model.Satellites
	DoneCh    chan struct{}
	Listen    net.Listener
	Route     route.Route
	RouteName string
}

func (s *Server) Init() {
	if len(s.RouteName) > 0 {
		s.Route = route.GetRoute(s.RouteName)
	}

	for _, _satellite := range s.Satellites.Satellites {
		_msg := satellite.GetSatellite(_satellite.Name)
		if _msg != nil {
			msgPipeline, level := _msg.Subscribe(s.Name(), s.Pipeline.PipeCh, _satellite.Level)
			s.Subscribe(_satellite.Name, msgPipeline, level)
		} else {
			log.WarnF("%v server get msg failed pipeline name: %v", s.Name(), _satellite.Name)
		}
	}
}
