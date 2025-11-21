package discovery

import (
	"context"
	"fmt"
	"ingress-migration-analyzer/internal/models"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// listIngresses gets all Ingress resources from the cluster
func (s *Scanner) listIngresses(ctx context.Context) ([]networkingv1.Ingress, error) {
	var allIngresses []networkingv1.Ingress

	if s.namespace != "" {
		// Scan specific namespace
		ingresses, err := s.listIngressesInNamespace(ctx, s.namespace)
		if err != nil {
			return nil, err
		}
		allIngresses = ingresses
	} else {
		// Scan all namespaces
		namespaces, err := s.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}

		for _, ns := range namespaces.Items {
			ingresses, err := s.listIngressesInNamespace(ctx, ns.Name)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to list ingresses in namespace %s: %v\n", ns.Name, err)
				continue
			}
			allIngresses = append(allIngresses, ingresses...)
		}
	}

	return allIngresses, nil
}

// listIngressesInNamespace lists ingresses in a specific namespace
func (s *Scanner) listIngressesInNamespace(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	ingressList, err := s.client.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses in namespace %s: %w", namespace, err)
	}

	return ingressList.Items, nil
}

// filterNginxIngresses filters ingresses that use nginx
func (s *Scanner) filterNginxIngresses(ingresses []networkingv1.Ingress) []networkingv1.Ingress {
	var nginxIngresses []networkingv1.Ingress

	for _, ingress := range ingresses {
		if s.isNginxIngress(ingress) {
			nginxIngresses = append(nginxIngresses, ingress)
		}
	}

	return nginxIngresses
}

// isNginxIngress determines if an Ingress uses nginx
func (s *Scanner) isNginxIngress(ingress networkingv1.Ingress) bool {
	// Check IngressClassName
	if ingress.Spec.IngressClassName != nil && *ingress.Spec.IngressClassName == s.ingressClassName {
		return true
	}

	// Check legacy annotation
	if class, exists := ingress.Annotations["kubernetes.io/ingress.class"]; exists && class == s.ingressClassName {
		return true
	}

	// Check for any nginx-specific annotations
	for key := range ingress.Annotations {
		if strings.HasPrefix(key, "nginx.ingress.kubernetes.io/") {
			return true
		}
	}

	return false
}

// getIngressClass extracts the ingress class name
func (s *Scanner) getIngressClass(ingress networkingv1.Ingress) string {
	if ingress.Spec.IngressClassName != nil {
		return *ingress.Spec.IngressClassName
	}

	// Fall back to annotation
	if class, exists := ingress.Annotations["kubernetes.io/ingress.class"]; exists {
		return class
	}

	return ""
}

// extractHosts extracts all hostnames from an Ingress
func (s *Scanner) extractHosts(ingress networkingv1.Ingress) []string {
	var hosts []string
	seen := make(map[string]bool)

	for _, rule := range ingress.Spec.Rules {
		if rule.Host != "" && !seen[rule.Host] {
			hosts = append(hosts, rule.Host)
			seen[rule.Host] = true
		}
	}

	return hosts
}

// extractPaths extracts all paths from an Ingress
func (s *Scanner) extractPaths(ingress networkingv1.Ingress) []string {
	var paths []string
	seen := make(map[string]bool)

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				pathStr := path.Path
				if pathStr == "" {
					pathStr = "/"
				}
				if !seen[pathStr] {
					paths = append(paths, pathStr)
					seen[pathStr] = true
				}
			}
		}
	}

	return paths
}

// convertToModel converts Kubernetes Ingress to our internal model
func (s *Scanner) ingressToModel(ingresses []networkingv1.Ingress) []models.IngressResource {
	var resources []models.IngressResource

	for _, ingress := range ingresses {
		resource := models.IngressResource{
			Name:        ingress.Name,
			Namespace:   ingress.Namespace,
			ClassName:   s.getIngressClass(ingress),
			Annotations: s.copyMap(ingress.Annotations),
			Labels:      s.copyMap(ingress.Labels),
			Hosts:       s.extractHosts(ingress),
			Paths:       s.extractPaths(ingress),
			CreatedAt:   ingress.CreationTimestamp.Time,
		}
		resources = append(resources, resource)
	}

	return resources
}
