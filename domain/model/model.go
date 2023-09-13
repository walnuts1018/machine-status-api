package model

type Config struct {
	PVEApiUrl     string
	PVEApiTokenID string
	PVEApiSecret  string
	Port          int
}

type MachineStatus int

const (
	Unknown   MachineStatus = iota // Status Unknown
	Healthy                        // Power ON and can access
	Unhealthy                      //Power ON, but can not access
	Inactive                       // Power OFF, but can not access
)

func (s MachineStatus) String() string {
	switch s {
	case 1:
		return "Healthy"
	case 0:
		return "Unknown"
	case -1:
		return "Unhealthy"
	case -2:
		return "Inactive"
	default:
		return "Error"
	}
}
