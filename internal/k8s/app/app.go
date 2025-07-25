package app

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

// NewApp creates a new instance of the application
func NewApp() *App {
	a := &App{
		App:              tview.NewApplication(),
		pages:            tview.NewPages(),
		NsList:           tview.NewList(),

		ResourceList:     tview.NewList(),
		ResourceTypeList: tview.NewList(),
		InfoView:         tview.NewTextView().SetDynamicColors(true),
		LogsView:         tview.NewTextView().SetDynamicColors(true),
		CurrentFocus:     0,
		stopChan:         make(chan struct{}),
		logStopChan:      make(chan struct{}),
	}

	// Initialize the UI
	a.initUI()

	// Initialize the Kubernetes client
	if err := a.initKubeClient(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Kubernetes client: %v\n", err)
	}

	return a
}

// Run starts the application
func (a *App) Run() error {
	// Load initial data
	if err := a.LoadNamespaces(); err != nil {
		return fmt.Errorf("error loading namespaces: %v", err)
	}

	// Start the application
	if err := a.App.Run(); err != nil {
		return fmt.Errorf("error running application: %v", err)
	}

	// Clean up
	if a.logStream != nil {
		a.logStream.Close()
	}
	if a.logStopChan != nil {
		close(a.logStopChan)
	}

	return nil
}
