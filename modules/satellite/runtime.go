package satellite

var (
	satellites = make(map[string]Satellite)
)

func GetSatellite(name string) (t Satellite) {
	return satellites[name]
}

func GetAllSatellite() map[string]Satellite {
	return satellites
}

func CloseAllSatellite() {
	for _, satellite := range satellites {
		if satellite != nil {
			_ = satellite.Close()
		}
	}
}
