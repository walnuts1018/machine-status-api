package usecase

import (
	"fmt"

	"github.com/walnuts1018/machine-status-api/domain"
	"github.com/walnuts1018/machine-status-api/domain/model"
)

type MachineUsecase struct {
	proxmoxClient domain.ProxmoxClient
	gpioClient    domain.GPIOClient
}

func NewClient(proxmoxClient domain.ProxmoxClient, gpioClient domain.GPIOClient) *MachineUsecase {
	return &MachineUsecase{
		proxmoxClient: proxmoxClient,
		gpioClient:    gpioClient,
	}
}

func (c MachineUsecase) StartMachine(machineName string) error {
	err := c.gpioClient.StartAlice()
	if err != nil {
		return fmt.Errorf("failed to start alice: %w", err)
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

func (c MachineUsecase) StopMachine(machineName string) error {
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

func (c MachineUsecase) GetMachineStatus(machineName string) (model.MachineStatus, error) {
	if machineName == "alice" {
		return c.getAliceStatus()
	}

	return c.proxmoxClient.GetMachineStatus("alice", machineName)
}

func (c MachineUsecase) getAliceStatus() (model.MachineStatus, error) {
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
