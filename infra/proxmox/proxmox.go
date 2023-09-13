package proxmox

import (
	"fmt"
	"log/slog"
	"time"

	proxmox "github.com/luthermonson/go-proxmox"
	"github.com/walnuts1018/machine-status-api/domain"
	"github.com/walnuts1018/machine-status-api/domain/model"
)

const (
	interval        = 5 * time.Second
	defaultTimeout  = 1 * time.Minute
	shutdownTimeout = 5 * time.Minute
)

type client struct {
	pveClient *proxmox.Client
}

func NewClient(config *model.Config) domain.ProxmoxClient {
	return client{pveClient: proxmox.NewClient(config.PVEApiUrl, proxmox.WithAPIToken(config.PVEApiTokenID, config.PVEApiSecret))}
}

func (c client) listNode() ([]*proxmox.Node, error) {
	nodeStatuses, err := c.pveClient.Nodes()
	if err != nil {
		return nil, err
	}

	nodes := make([]*proxmox.Node, len(nodeStatuses))
	for _, nodeStatus := range nodeStatuses {
		node, err := c.pveClient.Node(nodeStatus.Name)
		if err != nil {
			slog.Warn("failed to get node", "NodeName", nodeStatus.Name, "error", err)
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (c client) FindMachineNameByID(vmID int) (string, error) {
	nodes, err := c.listNode()
	if err != nil {
		return "", err
	}

	for _, node := range nodes {
		v, err := node.VirtualMachine(vmID)
		if err != nil {
			continue
		}
		return v.Name, nil
	}
	return "", fmt.Errorf("failed to find machine name by id: %d", vmID)
}

func (c client) StartMachine(nodeName string, vmName string) error {
	node, err := c.pveClient.Node(nodeName)
	if err != nil {
		return err
	}

	vms, err := node.VirtualMachines()
	if err != nil {
		return err
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			isrunning := vm.IsRunning()
			if isrunning {
				slog.Warn("machine is already running", "machineName", vmName)
				return nil
			}

			task, err := vm.Start()
			if err != nil {
				return err
			}
			err = task.Wait(interval, defaultTimeout)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("failed to find vm by name: %s", vmName)
}

func (c client) StopMachine(nodeName string, vmName string) error {
	node, err := c.pveClient.Node(nodeName)
	if err != nil {
		return err
	}

	vms, err := node.VirtualMachines()
	if err != nil {
		return err
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			isrunning := vm.IsRunning()
			if !isrunning {
				slog.Warn("machine is already stopped", "machineName", vmName)
				return nil
			}

			task, err := vm.Shutdown()
			if err != nil {
				return err
			}

			err = task.Wait(interval, shutdownTimeout)
			if err != nil {
				slog.Warn("failed to shutdown, try force stop", "error", err)
				task, err := vm.Stop()
				if err != nil {
					return err
				}

				err = task.Wait(interval, defaultTimeout)
				if err != nil {
					return err
				}
				return nil
			}
			return nil
		}
	}
	return fmt.Errorf("failed to find vm by name: %s", vmName)
}

func (c client) GetMachineStatus(nodeName string, vmName string) (model.MachineStatus, error) {
	node, err := c.pveClient.Node(nodeName)
	if err != nil {
		return model.Unknown, err
	}

	vms, err := node.VirtualMachines()
	if err != nil {
		return model.Unknown, err
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			vm.Ping()
			if vm.IsRunning() {
				return model.Healthy, nil
			} else {
				return model.Inactive, nil
			}
		}
	}
	return model.Unknown, fmt.Errorf("failed to find vm by name: %s", vmName)
}

func (c client) IsPVEServerActive() bool {
	_, err := c.pveClient.Version()
	if err != nil {
		return false
	}
	return true
}
