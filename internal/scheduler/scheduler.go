package scheduler

import (
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"

	"go.uber.org/zap"
)

const bindRetries = 6

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
		nodes := s.nodesListener.GetNodes()
		if len(nodes) == 0 {
			continue
		}

		usedNodes := make(map[string]struct{})
		for _, p := range pods {
			if p.SchedulerName == s.cfg.SchedulerName && p.NodeName != "" {
				usedNodes[p.NodeName] = struct{}{}
			}
		}

		var pending []k8s.PodInfo
		for _, p := range pods {
			if p.SchedulerName == s.cfg.SchedulerName && p.Phase == "Pending" && p.NodeName == "" {
				pending = append(pending, p)
			}
		}

		for _, pod := range pending {
			var chosen string
			for _, n := range nodes {
				if _, used := usedNodes[n]; !used {
					chosen = n
					break
				}
			}
			if chosen == "" {
				log.Warn("No free node for pod", zap.String("pod", pod.Namespace+"/"+pod.Name))
				continue
			}
			if err := s.bindPodToNodeWithRetry(pod, chosen); err != nil {
				log.Error("Failed to bind pod after retries", zap.String("pod", pod.Namespace+"/"+pod.Name), zap.String("node", chosen), zap.Error(err))
				continue
			}
			usedNodes[chosen] = struct{}{}
			log.Info("Bound pod to node", zap.String("pod", pod.Namespace+"/"+pod.Name), zap.String("node", chosen))
		}
	}
}

// bind node to pod with exponential backoff for retries
func (s *Scheduler) bindPodToNodeWithRetry(pod k8s.PodInfo, nodeName string) error {
	log := logger.GetLogger()
	var lastErr error
	backoff := time.Second
	for attempt := 0; attempt <= bindRetries; attempt++ {
		lastErr = s.k8sClient.BindPodToNode(pod.Namespace, pod.Name, nodeName)
		if lastErr == nil {
			return nil
		}
		if attempt < bindRetries {
			log.Warn("Bind failed, retrying", zap.String("pod", pod.Namespace+"/"+pod.Name), zap.Int("attempt", attempt+1), zap.Duration("backoff", backoff), zap.Error(lastErr))
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return lastErr
}