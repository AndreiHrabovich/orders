package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Cfg struct {
	Env string
}

func main() {
	var cfg Cfg
	logger, _ := zap.NewProduction()

	run(&cfg, logger)
}

func run(cfg *Cfg, logger *zap.Logger) {
	//  export APP_ENV=dev
	logger.Sugar().Log(zap.InfoLevel, "Starting the order service...")
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		logger.Sugar().Error("Failed to proccess the APP environment variable")
		os.Exit(1)
	}

	logger.Sugar().Log(zap.InfoLevel, cfg.Env)
}
