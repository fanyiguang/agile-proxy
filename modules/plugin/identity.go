package plugin

type Identity struct {
	ModuleName string
	ModuleType string
}

func (i Identity) Name() string {
	return i.ModuleName
}

func (i Identity) Type() string {
	return i.ModuleType
}
