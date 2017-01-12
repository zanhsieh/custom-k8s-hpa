package k8sclient

import (
	"fmt"

	apiv1 "k8s.io/client-go/pkg/api/v1"
)

var _ = K8sClient(&MockK8sClient{})

// MockK8sClient implements K8sClientInterface
type MockK8sClient struct {
	NumOfNodes    int
	NumOfCores    int
	NumOfReplicas int
	ConfigMap     *apiv1.ConfigMap
}

// FetchConfigMap mocks fetching the requested configmap from the Apiserver
func (k *MockK8sClient) FetchConfigMap(namespace, configmap string) (*apiv1.ConfigMap, error) {
	if k.ConfigMap.ObjectMeta.ResourceVersion == "" {
		return nil, fmt.Errorf("config map not exist")
	}
	return k.ConfigMap, nil
}

// CreateConfigMap mocks creating a configmap with given namespace, name and params
func (k *MockK8sClient) CreateConfigMap(namespace, configmap string, params map[string]string) (*apiv1.ConfigMap, error) {
	return nil, nil
}

// UpdateConfigMap mocks updating a configmap with given namespace, name and params
func (k *MockK8sClient) UpdateConfigMap(namespace, configmap string, params map[string]string) (*apiv1.ConfigMap, error) {
	return nil, nil
}

// GetClusterStatus mocks counting schedulable nodes and cores in the cluster
func (k *MockK8sClient) GetClusterStatus() (*ClusterStatus, error) {
	return &ClusterStatus{int32(k.NumOfNodes), int32(k.NumOfNodes), int32(k.NumOfCores), int32(k.NumOfCores)}, nil
}

// GetNamespace mocks returning the namespace of target resource.
func (k *MockK8sClient) GetNamespace() string {
	return ""
}

// UpdateReplicas mocks updating the number of replicas for the resource and return the previous replicas count
func (k *MockK8sClient) UpdateReplicas(expReplicas int32) (int32, error) {
	prevReplicas := int32(k.NumOfReplicas)
	k.NumOfReplicas = int(expReplicas)
	return prevReplicas, nil
}
