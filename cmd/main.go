package main

import (
	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
)

func main() {
	cfg := config.Get()
	logger.InitLogger(cfg)

	log := logger.GetLogger()
	log.Info("Starting custom kubernetes scheduler....")

	


}