package mocks

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// NewFakeKubernetesClient creates a new fake Kubernetes client with test data
func NewFakeKubernetesClient() *fake.Clientset {
	// Create test namespaces
	namespaces := []runtime.Object{
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "kube-system",
			},
		},
	}

	// Create test pods
	pods := []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "test-container",
						Image: "nginx:latest",
					},
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:  "test-container",
						Ready: true,
					},
				},
			},
		},
	}

	// Create a fake clientset with our test data
	return fake.NewSimpleClientset(append(namespaces, pods...)...)
}

// GetCoreV1Client returns a mock CoreV1Interface
func GetCoreV1Client() corev1client.CoreV1Interface {
	return NewFakeKubernetesClient().CoreV1()
}

// GetNamespaces returns a list of mock namespaces
func GetNamespaces() *corev1.NamespaceList {
	return &corev1.NamespaceList{
		Items: []corev1.Namespace{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "kube-system",
				},
			},
		},
	}
}

// GetPods returns a list of mock pods
func GetPods() *corev1.PodList {
	return &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
		},
	}
}
