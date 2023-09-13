package main

import (
	"context"
	"fmt"
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
	slog.Info("Loaded config")

	proxmoxClient := proxmox.NewClient(config)
	gpioClient := gpio.NewClient()
	machineUsecase := usecase.NewClient(proxmoxClient, gpioClient)
	handler := handler.NewHandler(machineUsecase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		slog.Info("start Task Loop")
		usecase.Run(ctx)
	}()

	slog.Info("start handler")
	err = handler.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}
