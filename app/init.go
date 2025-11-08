package app

import (
	"log/slog"
	"os"
	"path"
)

const (
	envMode         = "MODE"
	modeDevelopment = "development"
)

var (
	mode        string
	cmdName     string
	hostname, _ = os.Hostname()
)

func Init(cmd string) {
	cmdName = cmd
	mode = os.Getenv(envMode)
	if mode == "" {
		mode = modeDevelopment
	}
	workdir, _ := os.Getwd()
	initConfig(path.Join(workdir, "configs"))
	initLog()
	slog.Info("APP initialized")
}

func InitWithConfigPath(cmd string, configPath string) {
	cmdName = cmd
	mode = os.Getenv(envMode)
	if mode == "" {
		mode = modeDevelopment
	}
	initConfig(configPath)
	initLog()
	slog.Info("APP initialized")
}
