package app

import (
	"io"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	App            *tview.Application `json:"-"` // JSON tag "-" to avoid circular reference
	pages          *tview.Pages       // For modals and other pages
	grid           *tview.Grid        // Main grid layout for responsive updates
	NsList         *tview.List
	PodList        *tview.List
	ContList       *tview.List
	InfoView       *tview.TextView
	LogsView       *tview.TextView
	KubeClient     kubernetes.Interface
	RestConfig     *rest.Config
	CurrentNs      string
	SelectedNs     string
	SelectedPod    string
	SelectedCont   string
	CurrentFocus   int
	stopChan       chan struct{}
	logStream      io.ReadCloser
	logStopChan    chan struct{}
}
