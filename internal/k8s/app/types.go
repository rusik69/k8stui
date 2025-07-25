package app

import (
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	App                  *tview.Application `json:"-"` // JSON tag "-" to avoid circular reference
	pages                *tview.Pages       // For modals and other pages
	grid                 *tview.Grid        // Main grid layout for responsive updates
	NsList               *tview.List

	ResourceList         *tview.List       // New list for other resource types
	ResourceTypeList     *tview.List       // List to select resource type
	InfoView             *tview.TextView
	LogsView             *tview.TextView
	KubeClient           kubernetes.Interface
	RestConfig           *rest.Config
	CurrentNs            string
	SelectedNs           string
	SelectedPod          string
	SelectedCont         string
	SelectedResource     string
	SelectedResourceType ResourceType
	CurrentFocus         int
	stopChan             chan struct{}
	logStream            io.ReadCloser
	logStopChan          chan struct{}
}
