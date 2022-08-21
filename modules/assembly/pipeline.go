package assembly

import (
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"time"
)

// Pipeline RealTimeSubObjs的订阅都是需要实时响应的，它可能不希望收到无用的广播消息，
// 请精准投递消息给它。一些实时性要求不高的订阅频道请放在SubObjs中。
type Pipeline struct {
	Level           string
	PipeCh          chan model.ModuleMessage
	SubObjs         map[string]chan model.ModuleMessage // name:write chan
	RealTimeSubObjs map[string]chan model.ModuleMessage // name:write chan
}

func (p *Pipeline) Subscribe(name string, writeCh chan model.ModuleMessage, level string) (chan model.ModuleMessage, string) {
	switch level {
	case RealTime:
		p.RealTimeSubObjs[name] = writeCh
	case Normal:
		fallthrough
	default:
		p.SubObjs[name] = writeCh
	}
	return p.PipeCh, p.Level
}

func (p *Pipeline) AsyncSendMsgByName(msgName string, moduleName string, action int, content string) {
	// 异步对外发送消息，减少对主流程的影响
	// 对外保持0信任原则，设置超时时间如果
	// 外部阻塞也不会导致协程泄漏。
	if subObj, ok := p.RealTimeSubObjs[msgName]; ok {
		go func() {
			select {
			case subObj <- model.ModuleMessage{
				Message: model.Message{
					Content: content,
					Action:  action,
				},
				Name: moduleName,
			}:
			case <-time.After(time.Second * 3):
				log.InfoF("pipeline message lock: %v %v %v", content, moduleName, msgName)
			}
		}()
	}

	if subObj, ok := p.SubObjs[msgName]; ok {
		go func() {
			select {
			case subObj <- model.ModuleMessage{
				Message: model.Message{
					Content: content,
					Action:  action,
				},
				Name: moduleName,
			}:
			case <-time.After(time.Second * 3):
				log.InfoF("pipeline message lock: %v %v %v", content, moduleName, msgName)
			}
		}()
	}
}

func (p *Pipeline) AsyncSendMsg(moduleName string, action int, content string) {
	// 异步对外发送消息，减少对主流程的影响
	// 对外保持0信任原则，设置超时时间如果
	// 外部阻塞也不会导致协程泄漏。
	for msgName, subObj := range p.SubObjs {
		_subObj := subObj
		go func() {
			select {
			case _subObj <- model.ModuleMessage{
				Message: model.Message{
					Content: content,
					Action:  action,
				},
				Name: moduleName,
			}:
			case <-time.After(time.Second * 3):
				log.InfoF("pipeline message lock: %v %v %v", content, moduleName, msgName)
			}
		}()
	}
}

func (p *Pipeline) SendMsgByName(msgName string, moduleName string, action int, content string) {
	if subObj, ok := p.RealTimeSubObjs[msgName]; ok {
		select {
		case subObj <- model.ModuleMessage{
			Message: model.Message{
				Content: content,
				Action:  action,
			},
			Name: moduleName,
		}:
		case <-time.After(time.Second * 3):
			log.InfoF("pipeline message lock: %v %v %v", content, moduleName, msgName)
		}
	}

	if subObj, ok := p.SubObjs[msgName]; ok {
		select {
		case subObj <- model.ModuleMessage{
			Message: model.Message{
				Content: content,
				Action:  action,
			},
			Name: moduleName,
		}:
		case <-time.After(time.Second * 3):
			log.InfoF("pipeline message lock: %v %v %v", content, moduleName, msgName)
		}
	}
}

func (p *Pipeline) GetPipeCh() <-chan model.ModuleMessage {
	return p.PipeCh
}

func CreatePipeline() Pipeline {
	return Pipeline{
		PipeCh:          make(chan model.ModuleMessage),
		SubObjs:         make(map[string]chan model.ModuleMessage),
		RealTimeSubObjs: make(map[string]chan model.ModuleMessage),
	}
}
