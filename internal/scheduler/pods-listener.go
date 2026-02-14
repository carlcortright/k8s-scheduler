package scheduler

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"
)

type PodsListener struct {
	pods []k8s.PodInfo
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

			pods, err := l.k8sClient.GetPods()
			if err != nil {
				log.Error("Failed to get pods", zap.Error(err))
				continue
			}
			
			l.pods = pods

			l.lastUpdated = time.Now()
			l.mutex.Unlock()
		}
	}()
}

func (l *PodsListener) GetPods() []k8s.PodInfo {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.pods
}