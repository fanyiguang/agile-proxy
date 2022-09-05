package base

import (
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/msg"
	"agile-proxy/modules/route"
	"net"
)

type Server struct {
	assembly.Net
	assembly.Identity
	assembly.Pipeline
	model.PipelineInfos
	DoneCh    chan struct{}
	Listen    net.Listener
	Route     route.Route
	RouteName string
}

func (s *Server) Init() {
	if len(s.RouteName) > 0 {
		s.Route = route.GetRoute(s.RouteName)
	}

	for _, pipelineInfo := range s.PipelineInfo {
		_msg := msg.GetMsg(pipelineInfo.Name)
		if _msg != nil {
			msgPipeline, level := _msg.Subscribe(s.Name(), s.Pipeline.PipeCh, pipelineInfo.Level)
			s.Subscribe(pipelineInfo.Name, msgPipeline, level)
		} else {
			log.WarnF("%v server get msg failed pipeline name: %v", s.Name(), pipelineInfo.Name)
		}
	}
}
