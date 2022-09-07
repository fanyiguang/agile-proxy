package base

import "agile-proxy/modules/assembly"

type Satellite struct {
	assembly.Identity
	assembly.Pipeline
	Level int
}
