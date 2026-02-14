package scheduler

import (
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/logger"
)

type Scheduler struct {
	nodesListener *NodesListener
	podsListener *PodsListener
}

func NewScheduler(nodesListener *NodesListener, podsListener *PodsListener) *Scheduler {
	return &Scheduler{
		nodesListener: nodesListener,
		podsListener: podsListener,
	}
}	

func (s *Scheduler) StartScheduler() {
	log := logger.GetLogger()
	log.Info("Starting scheduler....")

	for {
		time.Sleep(1 * time.Second)
		log.Info("Scheduler running....")

		// TODO: implement scheduler logic
	}
}