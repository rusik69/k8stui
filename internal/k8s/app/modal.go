package app

import (
	"fmt"

	"github.com/rivo/tview"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// showConfirmationModal shows a confirmation dialog
func (a *App) showConfirmationModal(title, message string, callback func()) {
	// Create the modal
	modal := tview.NewModal()
	modal.SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.App.SetRoot(a.pages, true).SetFocus(a.getCurrentFocus())
			a.pages.RemovePage("confirmation")
			if buttonLabel == "Yes" && callback != nil {
				callback()
			}
		})

	// Add the modal to the pages
	a.pages.AddPage("confirmation", modal, true, true)
	a.App.SetFocus(modal)
}

// getCurrentFocus returns the currently focused UI element
func (a *App) getCurrentFocus() tview.Primitive {
	switch a.CurrentFocus {
	case 0:
		return a.NsList
	case 1:
		return a.ResourceTypeList
	case 2:
		return a.ResourceList
	default:
		return a.NsList
	}
}

// deleteCurrentResource deletes the currently selected resource
func (a *App) deleteCurrentResource() {
	switch a.CurrentFocus {
	case 0: // Namespace
		if a.SelectedNs == "" {
			return
		}
		a.showConfirmationModal(
			"Delete Namespace",
			fmt.Sprintf("Are you sure you want to delete namespace %s?\nThis action cannot be undone.", a.SelectedNs),
			a.deleteSelectedNamespace,
		)
	case 1: // Pod
		if a.SelectedPod == "" || a.CurrentNs == "" {
			return
		}
		a.showConfirmationModal(
			"Delete Pod",
			fmt.Sprintf("Are you sure you want to delete pod %s in namespace %s?\nThis action cannot be undone.", a.SelectedPod, a.CurrentNs),
			a.deleteSelectedPod,
		)
	}
}

// deleteSelectedNamespace deletes the currently selected namespace
func (a *App) deleteSelectedNamespace() {
	if a.KubeClient == nil || a.SelectedNs == "" {
		return
	}

	err := a.KubeClient.CoreV1().Namespaces().Delete(a.getContext(), a.SelectedNs, metav1.DeleteOptions{})
	if err != nil {
		a.showError(fmt.Sprintf("Error deleting namespace: %v", err))
		return
	}

	// Clear selection and reload namespaces
	a.SelectedNs = ""
	a.LoadNamespaces()
}

// deleteSelectedPod deletes the currently selected pod
func (a *App) deleteSelectedPod() {
	if a.KubeClient == nil || a.CurrentNs == "" || a.SelectedPod == "" {
		return
	}

	err := a.KubeClient.CoreV1().Pods(a.CurrentNs).Delete(a.getContext(), a.SelectedPod, metav1.DeleteOptions{})
	if err != nil {
		a.showError(fmt.Sprintf("Error deleting pod: %v", err))
		return
	}

	// Clear selection and reload pods
	a.SelectedPod = ""
	a.LoadPods()
}

// showError displays an error message in a modal
func (a *App) showError(message string) {
	modal := tview.NewModal()
	modal.SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.App.SetRoot(a.pages, true).SetFocus(a.getCurrentFocus())
			a.pages.RemovePage("error")
		})

	a.pages.AddPage("error", modal, true, true)
	a.App.SetFocus(modal)
}
