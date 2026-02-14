package scheduler

import (
	"sync"
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"

	"go.uber.org/zap"
)

type NodesListener struct {
	nodes []string
	mutex *sync.Mutex
	lastUpdated time.Time

	cfg *config.Config
	k8sClient *k8s.K8sClient	
}

func NewNodesListener(cfg *config.Config, k8sClient *k8s.K8sClient) *NodesListener {
	return &NodesListener{
		nodes: nil,
		mutex: &sync.Mutex{},
		cfg: cfg,
		k8sClient: k8sClient,
	}
}

func (l *NodesListener) StartNodesListener() {
	log := logger.GetLogger()
	log.Info("Starting nodes listener....")

	go func() {
		for {
			time.Sleep(l.cfg.PollingInterval)
			l.mutex.Lock()
			log.Info("Nodes listener updated at", zap.String("time", time.Now().Format(time.RFC3339)))
			l.lastUpdated = time.Now()
			l.mutex.Unlock()
		}
	}()
}