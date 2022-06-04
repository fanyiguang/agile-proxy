package plugin

type IdentInfo struct {
	ModuleName string
	ModuleType string
}

func (i IdentInfo) Name() string {
	return i.ModuleName
}

func (i IdentInfo) Type() string {
	return i.ModuleType
}
