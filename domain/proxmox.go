package domain

import "github.com/walnuts1018/machine-status-api/domain/model"

type ProxmoxClient interface {
	FindMachineNameByID(vmID int) (string, error)
	StartMachine(nodeName string, vmName string) error
	StopMachine(nodeName string, vmName string) error
	GetMachineStatus(nodeName string, vmName string) (model.MachineStatus, error)
	IsPVEServerActive() bool
}
