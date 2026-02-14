package main

import (
	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/scheduler"
)

func main() {
	cfg := config.Get()
	logger.InitLogger(cfg)

	log := logger.GetLogger()
	log.Info("Starting custom kubernetes scheduler....")

	nodesListener := scheduler.NewNodesListener(cfg)
	nodesListener.StartNodesListener()

	podsListener := scheduler.NewPodsListener(cfg)
	podsListener.StartPodsListener()

	scheduler := scheduler.NewScheduler(nodesListener, podsListener)
	scheduler.StartScheduler()
}