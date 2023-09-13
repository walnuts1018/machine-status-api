package mockProxmox

import (
	"fmt"

	"github.com/walnuts1018/machine-status-api/domain"
	"github.com/walnuts1018/machine-status-api/domain/model"
)

type mockMachine struct {
	name   string
	id     int
	status model.MachineStatus
}

type mockPveClient struct {
	isPVEServerActive bool
	machines          []*mockMachine
}

func NewClient(active bool) domain.ProxmoxClient {
	return &mockPveClient{
		isPVEServerActive: active,
		machines: []*mockMachine{
			{
				name:   "machine1",
				id:     100,
				status: model.Inactive,
			},
			{
				name:   "machine2",
				id:     101,
				status: model.Inactive,
			},
			{
				name:   "machine3",
				id:     102,
				status: model.Inactive,
			},
		},
	}
}

func (c *mockPveClient) FindMachineNameByID(vmID int) (string, error) {
	for _, machine := range c.machines {
		if machine.id == vmID {
			return machine.name, nil
		}
	}
	return "", fmt.Errorf("failed to find machine name by id: %d", vmID)
}

func (c *mockPveClient) StartMachine(nodeName string, vmName string) error {
	for _, machine := range c.machines {
		if machine.name == vmName {
			machine.status = model.Healthy
			return nil
		}
	}
	return fmt.Errorf("failed to start machine: %s", vmName)
}
func (c *mockPveClient) StopMachine(nodeName string, vmName string) error {
	for _, machine := range c.machines {
		if machine.name == vmName {
			machine.status = model.Inactive
			return nil
		}
	}
	return fmt.Errorf("failed to stop machine: %s", vmName)
}
func (c *mockPveClient) GetMachineStatus(nodeName string, vmName string) (model.MachineStatus, error) {
	for _, machine := range c.machines {
		if machine.name == vmName {
			return machine.status, nil
		}
	}
	return model.Unknown, fmt.Errorf("failed to get machine status: %s", vmName)
}
func (c *mockPveClient) IsPVEServerActive() bool {
	return c.isPVEServerActive
}
