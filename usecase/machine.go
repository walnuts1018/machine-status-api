package usecase

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/walnuts1018/machine-status-api/domain"
	"github.com/walnuts1018/machine-status-api/domain/model"
)

type MachineUsecase struct {
	proxmoxClient domain.ProxmoxClient
	gpioClient    domain.GPIOClient
	machines      map[string]int
	mutex         *sync.Mutex
}

func NewClient(proxmoxClient domain.ProxmoxClient, gpioClient domain.GPIOClient) *MachineUsecase {
	return &MachineUsecase{
		proxmoxClient: proxmoxClient,
		gpioClient:    gpioClient,
	}
}

func (c *MachineUsecase) StartMachine(machineName string) error {
	err := c.gpioClient.StartAlice()
	if err != nil {
		return fmt.Errorf("failed to start alice: %w", err)
	}

	timeout := time.After(5 * time.Minute)
LOOP:
	for {
		select {
		case <-timeout:
			return fmt.Errorf("failed to start alice: timeout")
		default:
			aliceStatus, err := c.getAliceStatus()
			if err != nil {
				return fmt.Errorf("failed to get alice status: %w", err)
			}
			if aliceStatus == model.Healthy {
				break LOOP
			}
		}
		time.Sleep(1 * time.Second)
	}

	if machineName == "alice" {
		return nil
	}

	err = c.proxmoxClient.StartMachine("alice", machineName)
	if err != nil {
		return fmt.Errorf("failed to start machine: %w", err)
	}
	return nil
}

func (c *MachineUsecase) StartMachineAutomated(machineName string) error {
	c.mutex.Lock()
	c.machines[machineName]++
	if c.machines[machineName] == 1 {
		err := c.StartMachine(machineName)
		if err != nil {
			return fmt.Errorf("failed to start machine: %w", err)
		}
	}
	slog.Info("machine started by automated action", "dependecy", c.machines[machineName])
	c.mutex.Unlock()
	return nil
}

func (c *MachineUsecase) StopMachine(machineName string) error {
	if machineName == "alice" {
		err := c.gpioClient.StopAlice()
		if err != nil {
			return fmt.Errorf("failed to stop alice: %w", err)
		}
		return nil
	}

	err := c.proxmoxClient.StopMachine("alice", machineName)
	if err != nil {
		return fmt.Errorf("failed to stop machine: %w", err)
	}
	return nil
}

func (c *MachineUsecase) StopMachineAutomated(machineName string) error {
	c.mutex.Lock()
	if c.machines[machineName] != 0 {
		c.machines[machineName]--
		if c.machines[machineName] == 0 {
			err := c.StopMachine(machineName)
			if err != nil {
				return fmt.Errorf("failed to stop machine: %w", err)
			}
		}
	}
	slog.Info("machine stopped by automated action", "dependecy", c.machines[machineName])
	c.mutex.Unlock()
	return nil
}

func (c *MachineUsecase) GetMachineStatus(machineName string) (model.MachineStatus, error) {
	aliceStatus, err := c.getAliceStatus()
	if err != nil {
		return aliceStatus, fmt.Errorf("failed to get alice status: %w", err)
	}

	if machineName == "alice" {
		return aliceStatus, nil
	}

	if aliceStatus != model.Healthy {
		return model.Inactive, nil
	}

	return c.proxmoxClient.GetMachineStatus("alice", machineName)
}

func (c *MachineUsecase) getAliceStatus() (model.MachineStatus, error) {
	pwon, err := c.gpioClient.IsPwOn()
	if err != nil {
		return model.Unknown, fmt.Errorf("failed to get read alice pw led voltage:%v", err)
	}

	if pwon {
		if c.proxmoxClient.IsPVEServerActive() {
			return model.Healthy, nil
		} else {
			return model.Unhealthy, nil
		}
	} else {
		return model.Inactive, nil
	}
}
