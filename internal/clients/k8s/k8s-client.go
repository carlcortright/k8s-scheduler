package k8s

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/carlcortright/k8s-scheduler/internal/config"
)

type K8sClient struct {
	cfg     *config.Config
	BaseURL string
	client  *http.Client
}

func NewK8sClient(cfg *config.Config) *K8sClient {
	baseURL := strings.TrimSuffix(cfg.K8sAPIServerURL, "/")
	var transport http.RoundTripper = http.DefaultTransport
	if strings.HasPrefix(baseURL, "https://") {
		transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	if cfg.K8sAuthTokenPath != "" {
		token, err := os.ReadFile(cfg.K8sAuthTokenPath)
		if err == nil {
			transport = &authTransport{base: transport, token: strings.TrimSpace(string(token))}
		}
	}
	return &K8sClient{
		cfg:     cfg,
		BaseURL: baseURL,
		client:  &http.Client{Transport: transport},
	}
}

type authTransport struct {
	base  http.RoundTripper
	token string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req2)
}


type nodeList struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"items"`
}

type podList struct {
	Items []podItem `json:"items"`
}

const PodGroupAnnotationKey = "pod-group"

type podItem struct {
	Metadata struct {
		Name        string            `json:"name"`
		Namespace   string            `json:"namespace"`
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		SchedulerName string `json:"schedulerName"`
		NodeName      string `json:"nodeName"`
	} `json:"spec"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}

type PodInfo struct {
	Namespace     string
	Name          string
	NodeName      string
	SchedulerName string
	Phase         string
	PodGroup      string 
}

func (c *K8sClient) GetNodes() ([]string, error) {
	url := c.BaseURL + "/api/v1/nodes"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get nodes: %s", resp.Status)
	}
	var list nodeList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(list.Items))
	for i := range list.Items {
		names = append(names, list.Items[i].Metadata.Name)
	}
	return names, nil
}

// GetPods returns pods from all namespaces. Caller can filter by schedulerName and phase.
func (c *K8sClient) GetPods() ([]PodInfo, error) {
	url := c.BaseURL + "/api/v1/pods"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get pods: %s", resp.Status)
	}
	var list podList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	out := make([]PodInfo, 0, len(list.Items))
	for i := range list.Items {
		out = append(out, podItemToPodInfo(&list.Items[i]))
	}
	return out, nil
}

func podItemToPodInfo(p *podItem) PodInfo {
	group := ""
	if p.Metadata.Annotations != nil {
		group = p.Metadata.Annotations[PodGroupAnnotationKey]
	}
	return PodInfo{
		Namespace:     p.Metadata.Namespace,
		Name:          p.Metadata.Name,
		NodeName:      p.Spec.NodeName,
		SchedulerName: p.Spec.SchedulerName,
		Phase:         p.Status.Phase,
		PodGroup:      group,
	}
}

// BindPodToNode creates a Binding so the pod is scheduled onto the given node.
// Namespace is required (binding is a namespaced subresource).
func (c *K8sClient) BindPodToNode(namespace, podName, nodeName string) error {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/binding", c.BaseURL, namespace, podName)
	body := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Binding",
		"target":     map[string]string{"name": nodeName},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bind pod %s/%s to %s: %s %s", namespace, podName, nodeName, resp.Status, string(b))
	}
	return nil
}