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

func SetCMDName(name string) {
	cmdName = name
}

func Init(workdir string) {
	mode = os.Getenv(envMode)
	if mode == "" {
		mode = modeDevelopment
	}
	initConfig(path.Join(workdir, "configs"))
	initLog()
	slog.Info("APP initialized")
}
