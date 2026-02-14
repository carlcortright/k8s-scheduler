package scheduler

import (
	"sync"
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"
)

type PodsListener struct {
	pods []string
	lastUpdated time.Time

	mutex *sync.Mutex

	cfg *config.Config
	k8sClient *k8s.K8sClient
}

func NewPodsListener(cfg *config.Config, k8sClient *k8s.K8sClient) *PodsListener {
	return &PodsListener{
		pods: nil,
		mutex: &sync.Mutex{},
		cfg: cfg,
		k8sClient: k8sClient,
	}
}

func (l *PodsListener) StartPodsListener() {
	log := logger.GetLogger()

	log.Info("Starting pods listener....")

	go func() {
		for {
			time.Sleep(l.cfg.PollingInterval)
			l.mutex.Lock()


			l.lastUpdated = time.Now()
			l.mutex.Unlock()
		}
	}()
}