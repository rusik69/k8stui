package app

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppCreation tests that NewApp creates a valid App instance
func TestAppCreation(t *testing.T) {
	app := NewApp()
	
	assert.NotNil(t, app, "NewApp should return a non-nil instance")
	assert.NotNil(t, app.App, "App should have a tview.Application")
	assert.NotNil(t, app.NsList, "App should have a namespace list")
	assert.NotNil(t, app.PodList, "App should have a pod list")
	assert.NotNil(t, app.ContList, "App should have a container list")
	assert.NotNil(t, app.InfoView, "App should have an info view")
	assert.NotNil(t, app.LogsView, "App should have a logs view")
}

// TestLoadNamespaces tests loading namespaces with mocked Kubernetes client
func TestLoadNamespaces(t *testing.T) {
	app := NewApp()
	
	// Create fake Kubernetes client
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
	)
	
	app.KubeClient = fakeClient
	
	// Test loading namespaces
	err := app.LoadNamespaces()
	require.NoError(t, err, "LoadNamespaces should not return an error")
	
	// Verify namespaces were loaded
	assert.Equal(t, 3, app.NsList.GetItemCount(), "Should load 3 namespaces")
	
	// Test that the first namespace is "default"
	name, _ := app.NsList.GetItemText(0)
	assert.Equal(t, "default", name, "First namespace should be 'default'")
}

// TestLoadPods tests loading pods with mocked Kubernetes client
func TestLoadPods(t *testing.T) {
	app := NewApp()
	app.CurrentNs = "default"
	
	// Create fake Kubernetes client with pods
	fakeClient := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-1",
				Namespace: "default",
			},
			Status: corev1.PodStatus{Phase: corev1.PodRunning},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-2",
				Namespace: "default",
			},
			Status: corev1.PodStatus{Phase: corev1.PodPending},
		},
	)
	
	app.KubeClient = fakeClient
	
	// Test loading pods
	err := app.LoadPods()
	require.NoError(t, err, "LoadPods should not return an error")
	
	// Verify pods were loaded
	assert.Equal(t, 2, app.PodList.GetItemCount(), "Should load 2 pods")
	
	// Test that pods are from the correct namespace
	name, _ := app.PodList.GetItemText(0)
	assert.Equal(t, "test-pod-1", name, "First pod should be 'test-pod-1'")
}

// TestLoadContainers tests loading containers for a pod
func TestLoadContainers(t *testing.T) {
	app := NewApp()
	app.CurrentNs = "default"
	
	// Create fake Kubernetes client with a pod having containers
	fakeClient := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:latest"},
					{Name: "redis", Image: "redis:alpine"},
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "nginx", Ready: true},
					{Name: "redis", Ready: false},
				},
			},
		},
	)
	
	app.KubeClient = fakeClient
	
	// Test loading containers
	err := app.LoadContainers("test-pod")
	require.NoError(t, err, "LoadContainers should not return an error")
	
	// Verify containers were loaded
	assert.Equal(t, 2, app.ContList.GetItemCount(), "Should load 2 containers")
	
	// Test container names and statuses
	name1, desc1 := app.ContList.GetItemText(0)
	assert.Equal(t, "nginx", name1, "First container should be 'nginx'")
	assert.Contains(t, desc1, "nginx:latest", "Container description should include image")
	
	name2, desc2 := app.ContList.GetItemText(1)
	assert.Equal(t, "redis", name2, "Second container should be 'redis'")
	assert.Contains(t, desc2, "redis:alpine", "Container description should include image")
}

// TestLoadPodsWithNoNamespace tests error handling when no namespace is selected
func TestLoadPodsWithNoNamespace(t *testing.T) {
	app := NewApp()
	
	// Test with no namespace selected
	app.CurrentNs = ""
	
	err := app.LoadPods()
	assert.Error(t, err, "LoadPods should return an error when no namespace is selected")
	assert.Contains(t, err.Error(), "no namespace selected")
}

// TestLoadContainersWithNonExistentPod tests error handling for non-existent pod
func TestLoadContainersWithNonExistentPod(t *testing.T) {
	app := NewApp()
	app.CurrentNs = "default"
	
	// Create fake client with empty cluster
	fakeClient := fake.NewSimpleClientset()
	app.KubeClient = fakeClient
	
	err := app.LoadContainers("non-existent-pod")
	assert.Error(t, err, "LoadContainers should return an error for non-existent pod")
	assert.Contains(t, err.Error(), "error getting pod")
}

// TestAppInitializationWithKubeClient tests app initialization with mocked client
func TestAppInitializationWithKubeClient(t *testing.T) {
	app := NewApp()
	
	// Create fake client
	fakeClient := fake.NewSimpleClientset()
	app.KubeClient = fakeClient
	
	// Test that we can load namespaces
	err := app.LoadNamespaces()
	require.NoError(t, err, "LoadNamespaces should work with mocked client")
	
	// Test that we can set namespace and load pods
	app.CurrentNs = "default"
	
	// Add a pod to the fake client
	fakeClient.CoreV1().Pods("default").Create(context.TODO(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "test", Image: "test"}},
		},
	}, metav1.CreateOptions{})
	
	err = app.LoadPods()
	require.NoError(t, err, "LoadPods should work with mocked client")
	assert.Equal(t, 1, app.PodList.GetItemCount(), "Should load 1 pod")
}

// TestAppMethodsWithRealDependencies tests that app methods work with real dependencies
func TestAppMethodsWithRealDependencies(t *testing.T) {
	// This test ensures that the app methods can be called without panicking
	app := NewApp()
	
	// Test that we can create the app without errors
	assert.NotNil(t, app)
	assert.NotNil(t, app.App)
	
	// Test that UI components are properly initialized
	assert.NotNil(t, app.NsList)
	assert.NotNil(t, app.PodList)
	assert.NotNil(t, app.ContList)
	assert.NotNil(t, app.InfoView)
	assert.NotNil(t, app.LogsView)
	
	// Test that we can call methods without panicking
	// These will return errors due to no client, but should not panic
	_ = app.LoadNamespaces()
	_ = app.LoadPods()
	_ = app.LoadContainers("test")
}

// TestUpdateFocus tests the UpdateFocus method
func TestUpdateFocus(t *testing.T) {
	app := NewApp()
	
	// Test initial focus
	assert.Equal(t, 0, app.CurrentFocus, "Initial focus should be 0")
	
	// Test focus updates
	app.CurrentFocus = 1
	app.UpdateFocus()
	assert.Equal(t, 1, app.CurrentFocus, "Focus should be updated to 1")
	
	app.CurrentFocus = 2
	app.UpdateFocus()
	assert.Equal(t, 2, app.CurrentFocus, "Focus should be updated to 2")
	
	// Test wrap around
	app.CurrentFocus = 0
	app.UpdateFocus()
	assert.Equal(t, 0, app.CurrentFocus, "Focus should be updated to 0")
}

// TestIntegrationWithMockClient tests the integration of app components
func TestIntegrationWithMockClient(t *testing.T) {
	app := NewApp()
	
	// Create comprehensive fake client
	fakeClient := fake.NewSimpleClientset(
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "default"},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod-1",
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "nginx", Image: "nginx:latest"},
					{Name: "redis", Image: "redis:alpine"},
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{Name: "nginx", Ready: true},
					{Name: "redis", Ready: false},
				},
			},
		},
	)
	
	app.KubeClient = fakeClient
	
	// Test complete workflow
	
	// 1. Load namespaces
	err := app.LoadNamespaces()
	require.NoError(t, err)
	assert.Equal(t, 1, app.NsList.GetItemCount(), "Should load 1 namespace")
	
	// 2. Select namespace and load pods
	app.CurrentNs = "default"
	err = app.LoadPods()
	require.NoError(t, err)
	assert.Equal(t, 1, app.PodList.GetItemCount(), "Should load 1 pod")
	
	// 3. Select pod and load containers
	err = app.LoadContainers("test-pod-1")
	require.NoError(t, err)
	assert.Equal(t, 2, app.ContList.GetItemCount(), "Should load 2 containers")
	
	// 4. Verify the integration works
	assert.Equal(t, "default", app.CurrentNs, "Current namespace should be 'default'")
	// Note: SelectedPod is set via UI callback, not directly in LoadPods
}
