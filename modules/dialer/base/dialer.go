package base

import "agile-proxy/modules/plugin"

type Dialer struct {
	plugin.IdentInfo
	plugin.OutputMsg
	IFace string
}
