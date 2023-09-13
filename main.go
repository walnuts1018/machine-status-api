package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/walnuts1018/machine-status-api/handler"
	"github.com/walnuts1018/machine-status-api/infra/config"
	"github.com/walnuts1018/machine-status-api/infra/gpio"
	"github.com/walnuts1018/machine-status-api/infra/proxmox"
	"github.com/walnuts1018/machine-status-api/usecase"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to create config", "error", err)
		os.Exit(1)
	}

	proxmoxClient := proxmox.NewClient(config)
	gpioClient := gpio.NewClient()
	machineUsecase := usecase.NewClient(proxmoxClient, gpioClient)
	handler := handler.NewHandler(machineUsecase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		usecase.Run(ctx)
	}()

	err = handler.Run()
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}
