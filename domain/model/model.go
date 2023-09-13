package model

type Config struct {
	PVEApiUrl     string
	PVEApiTokenID string
	PVEApiSecret  string
	Port          int
}

type MachineStatus string

const (
	Unknown   MachineStatus = "Unknown"   // Status Unknown
	Healthy   MachineStatus = "Healthy"   // Power ON and can access
	Unhealthy MachineStatus = "Unhealthy" //Power ON, but can not access
	Inactive  MachineStatus = "Inactive"  // Power OFF, but can not access
)
