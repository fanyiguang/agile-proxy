package assembly

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

func CreateIdentity(name, _type string) Identity {
	return Identity{
		ModuleName: name,
		ModuleType: _type,
	}
}
