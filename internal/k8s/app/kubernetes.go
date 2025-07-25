package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// initKubeClient initializes the Kubernetes client
func (a *App) initKubeClient() error {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %v", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %v", err)
	}

	a.KubeClient = clientset
	a.RestConfig = config
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
	for _, ns := range namespaces.Items {
		a.NsList.AddItem(ns.Name, "", 0, func() {
			a.CurrentNs = ns.Name
		a.SelectedNs = ns.Name
		a.LoadPods()
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

	pods, err := a.KubeClient.CoreV1().Pods(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing pods: %v", err)
	}

	a.PodList.Clear()
	for _, pod := range pods.Items {
		a.PodList.AddItem(pod.Name, "", 0, func() {
			a.SelectedPod = pod.Name
		a.LoadContainers(pod.Name)
		a.showPodStatus(&pod)
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

	a.ContList.Clear()
	for _, container := range pod.Spec.Containers {
		containerName := container.Name // Capture the container name in closure
		a.ContList.AddItem(container.Name, "", 0, func() {
			// Automatically show logs when container is selected
			a.ShowContainerLogs(containerName)
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
