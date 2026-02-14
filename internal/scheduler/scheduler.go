package scheduler

import (
	"sync"
	"time"

	"github.com/carlcortright/k8s-scheduler/internal/config"
	"github.com/carlcortright/k8s-scheduler/internal/logger"
	"github.com/carlcortright/k8s-scheduler/internal/clients/k8s"

	"go.uber.org/zap"
)

const bindRetries = 6

type Scheduler struct {
	cfg           *config.Config
	k8sClient     *k8s.K8sClient
	nodesListener *NodesListener
	podsListener  *PodsListener

	mu      sync.Mutex
	podsMap map[string]k8s.PodInfo // pod-name -> PodInfo, kept in sync with poller + bind/evict so that we don't have to block on pod or node api calls
}

func NewScheduler(cfg *config.Config, k8sClient *k8s.K8sClient, nodesListener *NodesListener, podsListener *PodsListener) *Scheduler {
	return &Scheduler{
		cfg:           cfg,
		k8sClient:     k8sClient,
		nodesListener: nodesListener,
		podsListener:  podsListener,
		podsMap:       make(map[string]k8s.PodInfo),
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

		// Refresh internal map from poller (source of truth). Key = pod name (single namespace).
		s.mu.Lock()
		s.podsMap = make(map[string]k8s.PodInfo, len(pods))
		for _, p := range pods {
			s.podsMap[p.Name] = p
		}
		podsCopy := make([]k8s.PodInfo, 0, len(s.podsMap))
		for _, p := range s.podsMap {
			podsCopy = append(podsCopy, p)
		}
		s.mu.Unlock()

		usedNodes := make(map[string]struct{})
		for _, p := range podsCopy {
			if p.SchedulerName == s.cfg.SchedulerName && p.NodeName != "" {
				usedNodes[p.NodeName] = struct{}{}
			}
		}

		var pending []k8s.PodInfo
		for _, p := range podsCopy {
			if p.SchedulerName == s.cfg.SchedulerName && p.Phase == "Pending" && p.NodeName == "" {
				pending = append(pending, p)
			}
		}

		nodeToPod := make(map[string]k8s.PodInfo)
		for _, p := range podsCopy {
			if p.SchedulerName == s.cfg.SchedulerName && p.NodeName != "" {
				nodeToPod[p.NodeName] = p
			}
		}

		groups := make(map[string][]k8s.PodInfo)
		for _, p := range pending {
			key := p.PodGroup
			if key == "" {
				key = p.Name
			}
			groups[key] = append(groups[key], p)
		}

		var unplaced []k8s.PodInfo
		for _, groupPods := range groups {
			var freeNodes []string
			for _, n := range nodes {
				if _, used := usedNodes[n]; !used {
					freeNodes = append(freeNodes, n)
				}
			}
			if len(freeNodes) < len(groupPods) {
				unplaced = append(unplaced, groupPods...)
				continue
			}
			allOk := true
			for i := range groupPods {
				if err := s.bindPodToNodeWithRetry(groupPods[i], freeNodes[i]); err != nil {
					log.Error("Failed to bind pod after retries", zap.String("pod", groupPods[i].Namespace+"/"+groupPods[i].Name), zap.String("node", freeNodes[i]), zap.Error(err))
					allOk = false
					break
				}
			}
			if allOk {
				for i := range groupPods {
					s.recordBind(groupPods[i].Name, freeNodes[i])
					usedNodes[freeNodes[i]] = struct{}{}
					nodeToPod[freeNodes[i]] = groupPods[i]
					log.Info("Bound pod to node", zap.String("pod", groupPods[i].Namespace+"/"+groupPods[i].Name), zap.String("node", freeNodes[i]))
				}
			} else {
				unplaced = append(unplaced, groupPods...)
			}
		}

		// Preemption: for each unplaced pending pod, if it has higher priority than a scheduled pod, evict and bind.
		for _, pod := range unplaced {
			var victimNode string
			var victim k8s.PodInfo
			for _, n := range nodes {
				if cur, ok := nodeToPod[n]; ok && cur.Priority < pod.Priority {
					victimNode = n
					victim = cur
					break
				}
			}
			if victimNode == "" {
				continue
			}
			if err := s.evictPodWithRetry(victim); err != nil {
				log.Error("Preemption evict failed", zap.String("victim", victim.Namespace+"/"+victim.Name), zap.Error(err))
				continue
			}
			s.recordEvict(victim.Name)
			delete(usedNodes, victimNode)
			delete(nodeToPod, victimNode)
			if err := s.bindPodToNodeWithRetry(pod, victimNode); err != nil {
				log.Error("Preemption bind failed after evict", zap.String("pod", pod.Namespace+"/"+pod.Name), zap.String("node", victimNode), zap.Error(err))
				continue
			}
			s.recordBind(pod.Name, victimNode)
			usedNodes[victimNode] = struct{}{}
			nodeToPod[victimNode] = pod
			log.Info("Preempted", zap.String("victim", victim.Namespace+"/"+victim.Name), zap.String("node", victimNode), zap.String("replacement", pod.Namespace+"/"+pod.Name))
		}
	}
}

// recordBind updates the internal pods map after a successful bind.
func (s *Scheduler) recordBind(podName, nodeName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.podsMap[podName]; ok {
		p.NodeName = nodeName
		s.podsMap[podName] = p
	}
}

// recordEvict removes the pod from the internal map after a successful eviction.
func (s *Scheduler) recordEvict(podName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.podsMap, podName)
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

// evictPodWithRetry evicts a pod with exponential backoff for up to bindRetries retries.
func (s *Scheduler) evictPodWithRetry(pod k8s.PodInfo) error {
	log := logger.GetLogger()
	var lastErr error
	backoff := time.Second
	for attempt := 0; attempt <= bindRetries; attempt++ {
		lastErr = s.k8sClient.EvictPod(pod.Namespace, pod.Name)
		if lastErr == nil {
			return nil
		}
		if attempt < bindRetries {
			log.Warn("Evict failed, retrying", zap.String("pod", pod.Namespace+"/"+pod.Name), zap.Int("attempt", attempt+1), zap.Duration("backoff", backoff), zap.Error(lastErr))
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return lastErr
}