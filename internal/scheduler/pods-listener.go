package scheduler

import (
	"sync"
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"go.uber.org/zap"
)

type PodsListener struct {
	pods []string
	lastUpdated time.Time

	mutex *sync.Mutex

	cfg *config.Config
}

func NewPodsListener(cfg *config.Config) *PodsListener {
	return &PodsListener{
		nodes: nil,
		mutex: &sync.Mutex{},
		cfg: cfg,
	}
}

func (l *PodsListener) StartPodsListener() {
	log := logger.GetLogger()

	log.Info("Starting pods listener....")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			l.mutex.Lock()
			log.Info("Pods listener updated at", zap.String("time", time.Now().Format(time.RFC3339)))
			l.lastUpdated = time.Now()
			l.mutex.Unlock()
		}
	}()
}