package discovery

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var configMapName = "ingress-nginx-controller"

// DiscoverConfigMap discovers the ingress-nginx configmap
func (s *Scanner) DiscoverGlobalConfigMap(ctx context.Context) (map[string]string, error) {
	var configMapData map[string]string

	fmt.Printf("🔍 Discovering ingress-nginx ConfigMap in %s/%s...\n", s.controllerNamespace, configMapName)
	configMap := s.getNginxConfigMap(ctx, s.controllerNamespace)
	if configMap != nil {
		fmt.Printf("✅ Found ingress-nginx ConfigMap: %s/%s\n", s.controllerNamespace, configMap.Name)
		return configMap.Data, nil
	}

	return configMapData, fmt.Errorf("configmap %s/%s not found", s.controllerNamespace, configMapName)
}

func (s *Scanner) getNginxConfigMap(ctx context.Context, namespace string) *v1.ConfigMap {
	configMap, err := s.client.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	return configMap
}
