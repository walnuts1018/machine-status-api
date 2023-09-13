package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/walnuts1018/machine-status-api/domain/model"
)

func NewConfig() (*model.Config, error) {
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	err := godotenv.Load("./.config/.env")
	if err != nil {
		slog.Warn("failed to load .env, use default", err)
	}

	pveApiUrl, ok := os.LookupEnv("PVE_API_URL")
	if !ok {
		return &model.Config{}, fmt.Errorf("failed to get PVE_API_URL")
	}

	pveApiTokenID, ok := os.LookupEnv("PVE_API_TOKEN_ID")
	if !ok {
		return &model.Config{}, fmt.Errorf("failed to get PVE_API_TOKEN_ID")
	}

	pveApiSecret, ok := os.LookupEnv("PVE_API_SECRET")
	if !ok {
		return &model.Config{}, fmt.Errorf("failed to get PVE_API_SECRET")
	}

	return &model.Config{
		PVEApiUrl:     pveApiUrl,
		PVEApiTokenID: pveApiTokenID,
		PVEApiSecret:  pveApiSecret,
		Port:          *port,
	}, nil
}
