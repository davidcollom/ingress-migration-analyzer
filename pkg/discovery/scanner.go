package discovery

import (
	"context"
	"fmt"
	"time"

	"ingress-migration-analyzer/internal/models"
)

// Scanner handles discovery of Ingress resources
type Scanner struct {
	client              *Client
	namespace           string
	ingressClassName    string
	controllerNamespace string
}

// NewScanner creates a new scanner instance
func NewScanner(client *Client, namespace, ingressClassName, controllerNamespace string) *Scanner {
	return &Scanner{
		client:              client,
		namespace:           namespace,
		ingressClassName:    ingressClassName,
		controllerNamespace: controllerNamespace,
	}
}

// ScanCluster scans the cluster for ingress-nginx resources
func (s *Scanner) ScanCluster(ctx context.Context) (*models.ScanResult, error) {

	fmt.Println("🔍 Scanning for Global Ingress-Nginx Config")
	nginxConfig, err := s.DiscoverGlobalConfigMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to discover global configmap: %w", err)
	}

	fmt.Println("🔍 Scanning cluster for Ingress resources...")

	// Get all Ingress resources
	ingresses, err := s.listIngresses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	fmt.Printf("📊 Found %d total Ingress resources\n", len(ingresses))

	// Filter for nginx ingresses
	nginxIngresses := s.filterNginxIngresses(ingresses)
	fmt.Printf("🎯 Found %d ingress-nginx resources\n", len(nginxIngresses))

	// Convert to our model
	ingressResources := s.ingressToModel(nginxIngresses)

	result := &models.ScanResult{
		ClusterVersion: s.client.ClusterVersion,
		ContextName:    s.client.Context,
		TotalIngresses: len(ingresses),
		NginxConfig:    nginxConfig,
		NginxIngresses: ingressResources,
		ScanTime:       time.Now(),
	}

	return result, nil
}
