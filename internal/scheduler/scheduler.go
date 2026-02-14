package scheduler

import (
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/config"
)

type Scheduler struct {
	cfg *config.Config
	k8sClient *k8s.K8sClient
	nodesListener *NodesListener
	podsListener *PodsListener
}

func NewScheduler(cfg *config.Config, k8sClient *k8s.K8sClient, nodesListener *NodesListener, podsListener *PodsListener) *Scheduler {
	return &Scheduler{
		cfg: cfg,
		k8sClient: k8sClient,
		nodesListener: nodesListener,
		podsListener: podsListener,
	}
}	

func (s *Scheduler) StartScheduler() {
	log := logger.GetLogger()
	log.Info("Starting scheduler....")

	for {
		time.Sleep(s.cfg.PollingInterval)
		log.Info("Scheduler running....")

		// TODO: implement scheduler logic
	}
}