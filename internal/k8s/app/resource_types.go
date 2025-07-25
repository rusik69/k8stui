package app

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceType represents different Kubernetes resource types
type ResourceType string

const (
	ResourceTypeNamespace   ResourceType = "namespace"
	ResourceTypePod         ResourceType = "pod"
	ResourceTypeDeployment  ResourceType = "deployment"
	ResourceTypeService     ResourceType = "service"
	ResourceTypeConfigMap   ResourceType = "configmap"
	ResourceTypeSecret      ResourceType = "secret"
	ResourceTypeIngress     ResourceType = "ingress"
	ResourceTypeNode        ResourceType = "node"
	ResourceTypeReplicaSet  ResourceType = "replicaset"
	ResourceTypeStatefulSet ResourceType = "statefulset"
	ResourceTypeDaemonSet   ResourceType = "daemonset"
	ResourceTypeJob         ResourceType = "job"
	ResourceTypeCronJob     ResourceType = "cronjob"
	ResourceTypePVC         ResourceType = "pvc"
	ResourceTypePV          ResourceType = "pv"
	ResourceTypeNetworkPolicy ResourceType = "networkpolicy"
	ResourceTypeServiceAccount ResourceType = "serviceaccount"
	ResourceTypeRole        ResourceType = "role"
	ResourceTypeRoleBinding ResourceType = "rolebinding"
	ResourceTypeClusterRole ResourceType = "clusterrole"
	ResourceTypeClusterRoleBinding ResourceType = "clusterrolebinding"
	ResourceTypeEndpoint    ResourceType = "endpoint"
	ResourceTypeHPA         ResourceType = "hpa"
	ResourceTypeLimitRange  ResourceType = "limitrange"
	ResourceTypeResourceQuota ResourceType = "resourcequota"
)

// ResourceInfo holds information about a Kubernetes resource
type ResourceInfo struct {
	Name      string
	Namespace string
	Type      ResourceType
	Status    string
	Age       string
	Ready     string
	Spec      interface{}
}

// ResourceManager handles loading and managing different resource types
func (a *App) LoadResources(resourceType ResourceType) error {
	switch resourceType {
	case ResourceTypeDeployment:
		return a.LoadDeployments()
	case ResourceTypeService:
		return a.LoadServices()
	case ResourceTypeConfigMap:
		return a.LoadConfigMaps()
	case ResourceTypeSecret:
		return a.LoadSecrets()
	case ResourceTypeIngress:
		return a.LoadIngresses()
	case ResourceTypeNode:
		return a.LoadNodes()
	case ResourceTypeReplicaSet:
		return a.LoadReplicaSets()
	case ResourceTypeStatefulSet:
		return a.LoadStatefulSets()
	case ResourceTypeDaemonSet:
		return a.LoadDaemonSets()
	case ResourceTypeJob:
		return a.LoadJobs()
	case ResourceTypeCronJob:
		return a.LoadCronJobs()
	case ResourceTypePVC:
		return a.LoadPVCs()
	case ResourceTypePV:
		return a.LoadPVs()
	case ResourceTypeNetworkPolicy:
		return a.LoadNetworkPolicies()
	case ResourceTypeServiceAccount:
		return a.LoadServiceAccounts()
	case ResourceTypeRole:
		return a.LoadRoles()
	case ResourceTypeRoleBinding:
		return a.LoadRoleBindings()
	case ResourceTypeClusterRole:
		return a.LoadClusterRoles()
	case ResourceTypeClusterRoleBinding:
		return a.LoadClusterRoleBindings()
	case ResourceTypeEndpoint:
		return a.LoadEndpoints()
	case ResourceTypeHPA:
		return a.LoadHPAs()
	case ResourceTypeLimitRange:
		return a.LoadLimitRanges()
	case ResourceTypeResourceQuota:
		return a.LoadResourceQuotas()
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// LoadDeployments loads all deployments in the current namespace
func (a *App) LoadDeployments() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	deployments, err := a.KubeClient.AppsV1().Deployments(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing deployments: %v", err)
	}

	a.ResourceList.Clear()
	for _, deployment := range deployments.Items {
		deployment := deployment // capture for closure
		readyReplicas := fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas)
		
		a.ResourceList.AddItem(deployment.Name, readyReplicas, 0, func() {
			a.SelectedResource = deployment.Name
			a.SelectedResourceType = ResourceTypeDeployment
			a.showDeploymentInfo(&deployment)
		})
	}

	return nil
}

// LoadServices loads all services in the current namespace
func (a *App) LoadServices() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	services, err := a.KubeClient.CoreV1().Services(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing services: %v", err)
	}

	a.ResourceList.Clear()
	for _, service := range services.Items {
		service := service // capture for closure
		
		typeStr := string(service.Spec.Type)
		clusterIP := service.Spec.ClusterIP
		if clusterIP == "" {
			clusterIP = "None"
		}
		
		info := fmt.Sprintf("%s %s", typeStr, clusterIP)
		
		a.ResourceList.AddItem(service.Name, info, 0, func() {
			a.SelectedResource = service.Name
			a.SelectedResourceType = ResourceTypeService
			a.showServiceInfo(&service)
		})
	}

	return nil
}

// LoadConfigMaps loads all configmaps in the current namespace
func (a *App) LoadConfigMaps() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	configMaps, err := a.KubeClient.CoreV1().ConfigMaps(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing configmaps: %v", err)
	}

	a.ResourceList.Clear()
	for _, cm := range configMaps.Items {
		cm := cm // capture for closure
		
		dataCount := fmt.Sprintf("%d keys", len(cm.Data))
		
		a.ResourceList.AddItem(cm.Name, dataCount, 0, func() {
			a.SelectedResource = cm.Name
			a.SelectedResourceType = ResourceTypeConfigMap
			a.showConfigMapInfo(&cm)
		})
	}

	return nil
}

// LoadSecrets loads all secrets in the current namespace
func (a *App) LoadSecrets() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	secrets, err := a.KubeClient.CoreV1().Secrets(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing secrets: %v", err)
	}

	a.ResourceList.Clear()
	for _, secret := range secrets.Items {
		secret := secret // capture for closure
		
		typeStr := string(secret.Type)
		dataCount := fmt.Sprintf("%d keys", len(secret.Data))
		
		a.ResourceList.AddItem(secret.Name, fmt.Sprintf("%s %s", typeStr, dataCount), 0, func() {
			a.SelectedResource = secret.Name
			a.SelectedResourceType = ResourceTypeSecret
			a.showSecretInfo(&secret)
		})
	}

	return nil
}

// LoadIngresses loads all ingresses in the current namespace
func (a *App) LoadIngresses() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	ingresses, err := a.KubeClient.NetworkingV1().Ingresses(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing ingresses: %v", err)
	}

	a.ResourceList.Clear()
	for _, ingress := range ingresses.Items {
		ingress := ingress // capture for closure
		
		var hosts []string
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" {
				hosts = append(hosts, rule.Host)
			}
		}
		
		hostInfo := "No hosts"
		if len(hosts) > 0 {
			hostInfo = strings.Join(hosts, ", ")
		}
		
		a.ResourceList.AddItem(ingress.Name, hostInfo, 0, func() {
			a.SelectedResource = ingress.Name
			a.SelectedResourceType = ResourceTypeIngress
			a.showIngressInfo(&ingress)
		})
	}

	return nil
}

// LoadNodes loads all nodes in the cluster
func (a *App) LoadNodes() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	nodes, err := a.KubeClient.CoreV1().Nodes().List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing nodes: %v", err)
	}

	a.ResourceList.Clear()
	for _, node := range nodes.Items {
		node := node // capture for closure
		
		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				status = "NotReady"
				break
			}
		}
		
		info := fmt.Sprintf("%s %s", node.Status.NodeInfo.KubeletVersion, status)
		
		a.ResourceList.AddItem(node.Name, info, 0, func() {
			a.SelectedResource = node.Name
			a.SelectedResourceType = ResourceTypeNode
			a.showNodeInfo(&node)
		})
	}

	return nil
}

// GetResourceDisplayName returns a human-readable name for resource types
func GetResourceDisplayName(resourceType ResourceType) string {
	switch resourceType {
	case ResourceTypeDeployment:
		return "Deployments"
	case ResourceTypeReplicaSet:
		return "ReplicaSets"
	case ResourceTypeStatefulSet:
		return "StatefulSets"
	case ResourceTypeDaemonSet:
		return "DaemonSets"
	case ResourceTypeJob:
		return "Jobs"
	case ResourceTypeCronJob:
		return "CronJobs"
	case ResourceTypeService:
		return "Services"
	case ResourceTypeConfigMap:
		return "ConfigMaps"
	case ResourceTypeSecret:
		return "Secrets"
	case ResourceTypeIngress:
		return "Ingresses"
	case ResourceTypeNetworkPolicy:
		return "NetworkPolicies"
	case ResourceTypePVC:
		return "PersistentVolumeClaims"
	case ResourceTypePV:
		return "PersistentVolumes"
	case ResourceTypeServiceAccount:
		return "ServiceAccounts"
	case ResourceTypeRole:
		return "Roles"
	case ResourceTypeRoleBinding:
		return "RoleBindings"
	case ResourceTypeClusterRole:
		return "ClusterRoles"
	case ResourceTypeClusterRoleBinding:
		return "ClusterRoleBindings"
	case ResourceTypeEndpoint:
		return "Endpoints"
	case ResourceTypeHPA:
		return "HorizontalPodAutoscalers"
	case ResourceTypeLimitRange:
		return "LimitRanges"
	case ResourceTypeResourceQuota:
		return "ResourceQuotas"
	case ResourceTypeNode:
		return "Nodes"
	default:
		return string(resourceType)
	}
}

// GetAllResourceTypes returns all supported resource types
func GetAllResourceTypes() []ResourceType {
	return []ResourceType{
		ResourceTypeDeployment,
		ResourceTypeReplicaSet,
		ResourceTypeStatefulSet,
		ResourceTypeDaemonSet,
		ResourceTypeJob,
		ResourceTypeCronJob,
		ResourceTypeService,
		ResourceTypeConfigMap,
		ResourceTypeSecret,
		ResourceTypeIngress,
		ResourceTypeNetworkPolicy,
		ResourceTypePVC,
		ResourceTypePV,
		ResourceTypeServiceAccount,
		ResourceTypeRole,
		ResourceTypeRoleBinding,
		ResourceTypeClusterRole,
		ResourceTypeClusterRoleBinding,
		ResourceTypeEndpoint,
		ResourceTypeHPA,
		ResourceTypeLimitRange,
		ResourceTypeResourceQuota,
		ResourceTypeNode,
	}
}
