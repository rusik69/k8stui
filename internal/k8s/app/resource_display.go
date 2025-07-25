package app

import (
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
)

// showDeploymentInfo displays detailed information about a deployment
func (a *App) showDeploymentInfo(deployment *appsv1.Deployment) error {
	if deployment == nil {
		return fmt.Errorf("deployment is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Deployment: [white]%s\n", deployment.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", deployment.Namespace))
	status.WriteString(fmt.Sprintf("[green]Status: [white]%d/%d ready\n", 
		deployment.Status.ReadyReplicas, *deployment.Spec.Replicas))
	status.WriteString(fmt.Sprintf("[green]Strategy: [white]%s\n", deployment.Spec.Strategy.Type))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(deployment.CreationTimestamp.Time)))
	
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		status.WriteString("\n[green]Containers:\n")
		for _, container := range deployment.Spec.Template.Spec.Containers {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white] (%s)\n", container.Name, container.Image))
		}
	}
	
	if len(deployment.Spec.Selector.MatchLabels) > 0 {
		status.WriteString("\n[green]Labels:\n")
		for k, v := range deployment.Spec.Selector.MatchLabels {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// showServiceInfo displays detailed information about a service
func (a *App) showServiceInfo(service *corev1.Service) error {
	if service == nil {
		return fmt.Errorf("service is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Service: [white]%s\n", service.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", service.Namespace))
	status.WriteString(fmt.Sprintf("[green]Type: [white]%s\n", service.Spec.Type))
	status.WriteString(fmt.Sprintf("[green]ClusterIP: [white]%s\n", service.Spec.ClusterIP))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(service.CreationTimestamp.Time)))
	
	if len(service.Spec.Ports) > 0 {
		status.WriteString("\n[green]Ports:\n")
		for _, port := range service.Spec.Ports {
			portStr := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
			if port.NodePort != 0 {
				portStr = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
			}
			if port.TargetPort.String() != "" {
				portStr = fmt.Sprintf("%s → %s (%s)", portStr, port.TargetPort.String(), port.Protocol)
			}
			status.WriteString(fmt.Sprintf("  [yellow]%s\n", portStr))
		}
	}
	
	if len(service.Spec.Selector) > 0 {
		status.WriteString("\n[green]Selectors:\n")
		for k, v := range service.Spec.Selector {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}
	
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		status.WriteString("\n[green]LoadBalancer IPs:\n")
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				status.WriteString(fmt.Sprintf("  [yellow]%s\n", ingress.IP))
			}
			if ingress.Hostname != "" {
				status.WriteString(fmt.Sprintf("  [yellow]%s\n", ingress.Hostname))
			}
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// showConfigMapInfo displays detailed information about a configmap
func (a *App) showConfigMapInfo(configMap *corev1.ConfigMap) error {
	if configMap == nil {
		return fmt.Errorf("configmap is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]ConfigMap: [white]%s\n", configMap.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", configMap.Namespace))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(configMap.CreationTimestamp.Time)))
	
	if len(configMap.Data) > 0 {
		status.WriteString("\n[green]Data:\n")
		for key, value := range configMap.Data {
			// Truncate long values for display
			valueStr := value
			if len(valueStr) > 100 {
				valueStr = valueStr[:97] + "..."
			}
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", key, valueStr))
		}
	}
	
	if len(configMap.BinaryData) > 0 {
		status.WriteString("\n[green]Binary Data:\n")
		for key := range configMap.BinaryData {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %d bytes\n", key, len(configMap.BinaryData[key])))
		}
	}
	
	if len(configMap.Labels) > 0 {
		status.WriteString("\n[green]Labels:\n")
		for k, v := range configMap.Labels {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// showSecretInfo displays detailed information about a secret
func (a *App) showSecretInfo(secret *corev1.Secret) error {
	if secret == nil {
		return fmt.Errorf("secret is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Secret: [white]%s\n", secret.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", secret.Namespace))
	status.WriteString(fmt.Sprintf("[green]Type: [white]%s\n", secret.Type))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(secret.CreationTimestamp.Time)))
	
	if len(secret.Data) > 0 {
		status.WriteString("\n[green]Data:\n")
		for key := range secret.Data {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %d bytes (redacted)\n", key, len(secret.Data[key])))
		}
	}
	
	if len(secret.Labels) > 0 {
		status.WriteString("\n[green]Labels:\n")
		for k, v := range secret.Labels {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}
	
	if len(secret.Annotations) > 0 {
		status.WriteString("\n[green]Annotations:\n")
		for k, v := range secret.Annotations {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// showIngressInfo displays detailed information about an ingress
func (a *App) showIngressInfo(ingress *netv1.Ingress) error {
	if ingress == nil {
		return fmt.Errorf("ingress is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Ingress: [white]%s\n", ingress.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", ingress.Namespace))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(ingress.CreationTimestamp.Time)))
	
	if len(ingress.Spec.Rules) > 0 {
		status.WriteString("\n[green]Rules:\n")
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != "" {
				status.WriteString(fmt.Sprintf("  [yellow]Host[white]: %s\n", rule.Host))
			}
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					backend := ""
					if path.Backend.Service != nil {
						backend = fmt.Sprintf("%s:%d", path.Backend.Service.Name, path.Backend.Service.Port.Number)
					}
					status.WriteString(fmt.Sprintf("    [yellow]Path[white]: %s → %s\n", path.Path, backend))
				}
			}
		}
	}
	
	if ingress.Spec.DefaultBackend != nil {
		status.WriteString("\n[green]Default Backend:\n")
		if ingress.Spec.DefaultBackend.Service != nil {
			status.WriteString(fmt.Sprintf("  [yellow]Service[white]: %s\n", ingress.Spec.DefaultBackend.Service.Name))
			if ingress.Spec.DefaultBackend.Service.Port.Number != 0 {
				status.WriteString(fmt.Sprintf("  [yellow]Port[white]: %d\n", ingress.Spec.DefaultBackend.Service.Port.Number))
			}
		}
	}
	
	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		status.WriteString("\n[green]LoadBalancer IPs:\n")
		for _, ingress := range ingress.Status.LoadBalancer.Ingress {
			if ingress.IP != "" {
				status.WriteString(fmt.Sprintf("  [yellow]%s\n", ingress.IP))
			}
			if ingress.Hostname != "" {
				status.WriteString(fmt.Sprintf("  [yellow]%s\n", ingress.Hostname))
			}
		}
	}
	
	if len(ingress.Labels) > 0 {
		status.WriteString("\n[green]Labels:\n")
		for k, v := range ingress.Labels {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// showNodeInfo displays detailed information about a node
func (a *App) showNodeInfo(node *corev1.Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Node: [white]%s\n", node.Name))
	status.WriteString(fmt.Sprintf("[green]Age: [white]%s\n", getAge(node.CreationTimestamp.Time)))
	
	// Node status
	nodeStatus := "Unknown"
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				nodeStatus = "Ready"
			} else {
				nodeStatus = "NotReady"
			}
			break
		}
	}
	status.WriteString(fmt.Sprintf("[green]Status: [white]%s\n", nodeStatus))
	
	// System info
	status.WriteString(fmt.Sprintf("[green]OS: [white]%s\n", node.Status.NodeInfo.OSImage))
	status.WriteString(fmt.Sprintf("[green]Kernel: [white]%s\n", node.Status.NodeInfo.KernelVersion))
	status.WriteString(fmt.Sprintf("[green]Kubelet: [white]%s\n", node.Status.NodeInfo.KubeletVersion))
	status.WriteString(fmt.Sprintf("[green]Container Runtime: [white]%s\n", node.Status.NodeInfo.ContainerRuntimeVersion))
	
	// Capacity and allocatable
	status.WriteString("\n[green]Capacity:\n")
	for resource, quantity := range node.Status.Capacity {
		status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", resource, quantity.String()))
	}
	
	status.WriteString("\n[green]Allocatable:\n")
	for resource, quantity := range node.Status.Allocatable {
		status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", resource, quantity.String()))
	}
	
	// Addresses
	if len(node.Status.Addresses) > 0 {
		status.WriteString("\n[green]Addresses:\n")
		for _, addr := range node.Status.Addresses {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", addr.Type, addr.Address))
		}
	}
	
	// Labels
	if len(node.Labels) > 0 {
		status.WriteString("\n[green]Labels:\n")
		for k, v := range node.Labels {
			status.WriteString(fmt.Sprintf("  [yellow]%s[white]: %s\n", k, v))
		}
	}

	a.InfoView.SetText(status.String())
	return nil
}

// getAge returns a human-readable age string
func getAge(t time.Time) string {
	duration := time.Since(t)
	
	days := int(duration.Hours() / 24)
	if days > 365 {
		years := days / 365
		return fmt.Sprintf("%dy", years)
	} else if days > 30 {
		months := days / 30
		return fmt.Sprintf("%dmo", months)
	} else if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	
	hours := int(duration.Hours())
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	
	minutes := int(duration.Minutes())
	return fmt.Sprintf("%dm", minutes)
}
