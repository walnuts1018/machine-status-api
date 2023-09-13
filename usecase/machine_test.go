package usecase

import (
	"testing"

	"github.com/walnuts1018/machine-status-api/domain/model"
	"github.com/walnuts1018/machine-status-api/mock/mockGpio"
	"github.com/walnuts1018/machine-status-api/mock/mockProxmox"
)

func newClient() *MachineUsecase {
	proxmoxMockClient := mockProxmox.NewClient(true)
	gpioMockClient := mockGpio.NewClient()
	NewClient(proxmoxMockClient, gpioMockClient)
	return NewClient(proxmoxMockClient, gpioMockClient)
}

func TestMachineUsecase(t *testing.T) {
	usecase := newClient()
	type args struct {
		machineName string
	}
	tests := []struct {
		name    string
		c       MachineUsecase
		args    args
		wantErr bool
	}{
		{
			name: "start alice",
			c:    *usecase,
			args: args{
				machineName: "alice",
			},
			wantErr: false,
		},
		{
			name: "start machine1",
			c:    *usecase,
			args: args{
				machineName: "machine1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// start
			if err := tt.c.StartMachine(tt.args.machineName); (err != nil) != tt.wantErr {
				t.Errorf("MachineUsecase.StartMachine() error = %v, wantErr %v", err, tt.wantErr)
			}
			status, err := tt.c.GetMachineStatus(tt.args.machineName)
			if err != nil {
				t.Errorf("MachineUsecase.GetMachineStatus() error = %v", err)
			}
			if status != model.Healthy {
				t.Errorf("MachineUsecase.GetMachineStatus() status = %v, want %v", status, model.Healthy)
			}

			// stop
			if err := tt.c.StopMachine(tt.args.machineName); (err != nil) != tt.wantErr {
				t.Errorf("MachineUsecase.StartMachine() error = %v, wantErr %v", err, tt.wantErr)
			}
			status, err = tt.c.GetMachineStatus(tt.args.machineName)
			if err != nil {
				t.Errorf("MachineUsecase.GetMachineStatus() error = %v", err)
			}
			if status != model.Inactive {
				t.Errorf("MachineUsecase.GetMachineStatus() status = %v, want %v", status, model.Inactive)
			}
		})
	}
}
