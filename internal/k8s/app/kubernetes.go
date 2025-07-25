package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// initKubeClient initializes the Kubernetes client
func (a *App) initKubeClient() error {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Try to use the real Kubernetes config first
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err == nil {
		// Successfully got real config
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("error creating kubernetes client: %v", err)
		}
		a.KubeClient = clientset
		a.RestConfig = config
		return nil
	}

	// Fall back to fake client for demonstration/testing
	fmt.Fprintf(os.Stderr, "Warning: Using fake Kubernetes client for demonstration\n")
	
	// Create fake client with some sample data
	fakeClient := fake.NewSimpleClientset(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "default"},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "kube-system"},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod-1", Namespace: "default"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:latest"},
					{Name: "redis", Image: "redis:alpine"},
				},
			},
			Status: corev1.PodStatus{Phase: "Running"},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod-2", Namespace: "default"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "app", Image: "myapp:latest"},
				},
			},
			Status: corev1.PodStatus{Phase: "Running"},
		},
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "nginx-deployment", Namespace: "default"},
			Spec: appsv1.DeploymentSpec{
				Replicas: func(i int32) *int32 { return &i }(3),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "nginx"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "nginx"}},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{Name: "nginx", Image: "nginx:latest"},
						},
					},
				},
			},
			Status: appsv1.DeploymentStatus{ReadyReplicas: 3},
		},
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "redis-deployment", Namespace: "default"},
			Spec: appsv1.DeploymentSpec{
				Replicas: func(i int32) *int32 { return &i }(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "redis"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "redis"}},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{Name: "redis", Image: "redis:alpine"},
						},
					},
				},
			},
			Status: appsv1.DeploymentStatus{ReadyReplicas: 2},
		},
	)

	a.KubeClient = fakeClient
	// Use a mock config for fake client
	a.RestConfig = &rest.Config{Host: "fake-cluster"}
	return nil
}

// LoadNamespaces loads the list of namespaces from the Kubernetes cluster
func (a *App) LoadNamespaces() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	namespaces, err := a.KubeClient.CoreV1().Namespaces().List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing namespaces: %v", err)
	}

	a.NsList.Clear()
	a.ResourceList.Clear()
	for _, ns := range namespaces.Items {
		a.NsList.AddItem(ns.Name, "", 0, func() {
			a.CurrentNs = ns.Name
			a.SelectedNs = ns.Name
			a.ResourceList.Clear()
			a.InfoView.Clear()
			a.LoadResources(ResourceTypePod)
		})
	}

	return nil
}

// LoadPods loads the list of pods in the current namespace
func (a *App) LoadPods() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	// Hide logs window when displaying pods
	a.showLogsWindow(false)

	pods, err := a.KubeClient.CoreV1().Pods(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing pods: %v", err)
	}

	a.ResourceList.Clear()
	for _, pod := range pods.Items {
		podName := pod.Name // Capture the pod name in closure
		status := "Running"
		if pod.Status.Phase != corev1.PodRunning {
			status = string(pod.Status.Phase)
		}
		
		a.ResourceList.AddItem(pod.Name, status, 0, func() {
			a.SelectedPod = podName
			a.LoadContainers(podName)
			a.showPodStatus(&pod)
			// Update responsive layout after selection
			if a.grid != nil {
				a.updateGridLayout(a.grid)
			}
		})
	}

	return nil
}

// LoadContainers loads the containers for a given pod
func (a *App) LoadContainers(podName string) error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	pod, err := a.KubeClient.CoreV1().Pods(a.CurrentNs).Get(a.getContext(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting pod: %v", err)
	}

	// Show logs window when displaying containers
	a.showLogsWindow(true)

	a.ResourceList.Clear()
	for _, container := range pod.Spec.Containers {
		containerName := container.Name // Capture the container name in closure
		a.ResourceList.AddItem(container.Name, "", 0, func() {
			// Automatically show logs when container is selected
			a.ShowContainerLogs(containerName)
			// Update responsive layout after selection
			if a.grid != nil {
				a.updateGridLayout(a.grid)
			}
		})
	}

	return a.showPodStatus(pod)
}

// showPodStatus displays detailed status of a pod
func (a *App) showPodStatus(pod *corev1.Pod) error {
	if pod == nil {
		return fmt.Errorf("pod is nil")
	}

	// Stop any existing log stream
	if a.logStream != nil {
		a.logStream.Close()
	}
	if a.logStopChan != nil {
		close(a.logStopChan)
	}
	a.logStopChan = make(chan struct{})

	// Format pod status
	var status strings.Builder
	status.WriteString(fmt.Sprintf("[green]Pod: [white]%s\n", pod.Name))
	status.WriteString(fmt.Sprintf("[green]Namespace: [white]%s\n", pod.Namespace))
	status.WriteString(fmt.Sprintf("[green]Status: [white]%s\n", pod.Status.Phase))
	status.WriteString(fmt.Sprintf("[green]Node: [white]%s\n", pod.Spec.NodeName))
	status.WriteString(fmt.Sprintf("[green]IP: [white]%s\n", pod.Status.PodIP))
	status.WriteString("\n[green]Containers:\n")

	for _, container := range pod.Spec.Containers {
		status.WriteString(fmt.Sprintf("  [yellow]%s[white] (%s)\n", container.Name, container.Image))
	}

	a.InfoView.SetText(status.String())
	return nil
}

// getContext returns a background context for Kubernetes operations
func (a *App) getContext() context.Context {
	return context.Background()
}
