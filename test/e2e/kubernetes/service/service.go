package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Service represents a kubernetes service
type Service struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

// Metadata holds information like name, namespace, and labels
type Metadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// Spec holds information like clusterIP and port
type Spec struct {
	ClusterIP string `json:"clusterIP"`
	Ports     []Port `json:"ports"`
	Type      string `json:"type"`
}

// Port represents a service port definition
type Port struct {
	NodePort   int    `json:"nodePort"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	TargetPort int    `json:"targetPort"`
}

// Status holds the load balancer definition
type Status struct {
	LoadBalancer LoadBalancer `json:"loadBalancer"`
}

// LoadBalancer holds the ingress definitions
type LoadBalancer struct {
	Ingress []map[string]string `json:"ingress"`
}

// Get returns the service definition specified in a given namespace
func Get(name, namespace string) (*Service, error) {
	out, err := exec.Command("kubectl", "get", "svc", "-o", "json", "-n", namespace, name).CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get svc':%s\n", string(out))
		return nil, err
	}
	s := Service{}
	err = json.Unmarshal(out, &s)
	if err != nil {
		log.Printf("Error unmarshalling service json:%s\n", err)
		return nil, err
	}
	return &s, nil
}

// Delete will delete a service in a given namespace
func (s *Service) Delete() error {
	out, err := exec.Command("kubectl", "delete", "svc", "-n", s.Metadata.Namespace, s.Metadata.Name).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to delete service %s in namespace %s:%s\n", s.Metadata.Namespace, s.Metadata.Name, string(out))
		return err
	}
	return nil
}

// GetNodePort will return the node port for a given pod
func (s *Service) GetNodePort(port int) int {
	for _, p := range s.Spec.Ports {
		if p.Port == port {
			return p.NodePort
		}
	}
	return 0
}

// WaitForExternalIP waits for an external ip to be provisioned
func (s *Service) WaitForExternalIP(wait, sleep int) (*Service, error) {
	svcCh := make(chan *Service)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(wait))
	defer cancel()
	go func() {
		var svc *Service
		var err error
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded while waiting for External IP to be provisioned")
			default:
				svc, err = Get(s.Metadata.Name, s.Metadata.Namespace)
				if err != nil {
					errCh <- err
				}
				if svc.Status.LoadBalancer.Ingress != nil {
					svcCh <- svc
				}
				time.Sleep(time.Second * time.Duration(sleep))
			}
		}
	}()
	for {
		select {
		case err := <-errCh:
			return nil, err
		case svc := <-svcCh:
			return svc, nil
		}
	}
}
