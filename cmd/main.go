package main

import (
	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/scheduler"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"
)

func main() {
	cfg := config.Get()
	logger.InitLogger(cfg)

	log := logger.GetLogger()
	log.Info("Starting custom kubernetes scheduler....")

	k8sClient := k8s.NewK8sClient(cfg)

	nodesListener := scheduler.NewNodesListener(cfg, k8sClient)
	nodesListener.StartNodesListener()

	podsListener := scheduler.NewPodsListener(cfg, k8sClient)
	podsListener.StartPodsListener()

	scheduler := scheduler.NewScheduler(cfg, nodesListener, podsListener)
	scheduler.StartScheduler()
}