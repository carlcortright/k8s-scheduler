package scheduler

import (
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"

	"go.uber.org/zap"
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

		pods := s.podsListener.GetPods()
		for _, pod := range pods {
			log.Info("Pod: ", zap.String("pod", pod.Name))
		}

		// TODO: implement scheduler logic
	}
}