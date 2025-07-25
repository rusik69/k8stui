package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// initUI initializes the user interface
func (a *App) initUI() {
	// Configure the views
	a.NsList.SetBorder(true).SetTitle(" Namespaces ")
	a.ResourceTypeList.SetBorder(true).SetTitle(" Resource Types ")
	a.ResourceList.SetBorder(true).SetTitle(" Resources ")
	
	// Configure InfoView with scrolling
	a.InfoView.SetBorder(true).SetTitle(" Info ")
	a.InfoView.SetChangedFunc(func() {
		a.App.Draw()
	})

	// Configure LogsView with auto-scrolling
	a.LogsView.SetBorder(true).SetTitle(" Logs ")
	a.LogsView.SetChangedFunc(func() {
		a.LogsView.ScrollToEnd()
		a.App.Draw()
	})

	// Initialize resource type selection
	a.initResourceTypes()

	// Create the main grid layout
	a.grid = tview.NewGrid()
	hotkeyHelp := "[::b]TAB/Shift+TAB[::-] Navigate | [::b]ENTER[::-] Select | [::b]Ctrl+D[::-] Delete | [::b]Q[::-] Quit | [::b]↑/↓/←/→[::-] Scroll | [::b]Ctrl+R[::-] Resource Types"
	a.grid.SetBorder(true).SetTitle(" K8s TUI - " + hotkeyHelp + " ")

	// Set up the grid layout to be responsive to terminal size
	a.grid.SetRows(0, 0). // Two rows of equal height
		SetBorders(true)

	// Add items with dynamic sizing based on terminal width
	a.grid.AddItem(a.NsList, 0, 0, 1, 1, 0, 0, true).
		AddItem(a.ResourceTypeList, 0, 1, 1, 1, 0, 0, false).
		AddItem(a.ResourceList, 0, 2, 1, 1, 0, 0, false).
		AddItem(a.InfoView, 1, 0, 1, 3, 0, 0, false).
		AddItem(a.LogsView, 1, 2, 1, 1, 0, 0, false)

	// Set up dynamic column sizes based on terminal width
	a.updateGridLayout(a.grid)
	
	// Set up terminal resize handler to maintain responsive layout
	a.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		if screen != nil {
			a.updateGridLayout(a.grid)
		}
		return false
	})

	// Create a pages container that will hold our main UI and modals
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.grid, 0, 1, true)

	// Add the main UI to the pages
	a.pages.AddPage("main", mainFlex, true, true)

	// Set up the application with pages as root
	a.App.SetRoot(a.pages, true).
		SetFocus(a.NsList)

	// Set up key bindings
	a.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab:
			// Handle tab navigation
			a.navigate(event.Key() == tcell.KeyTab)
			return nil
		case tcell.KeyCtrlD:
			a.deleteCurrentResource()
			return nil
		case tcell.KeyCtrlR:
			// Show resource type selection
			a.showResourceTypeModal()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				a.App.Stop()
				return nil
			case 'r', 'R':
				// Show resource type selection
				a.showResourceTypeModal()
				return nil
			}
		}

		return event
	})
}

// updateGridLayout updates the grid column layout to utilize full terminal width
func (a *App) updateGridLayout(grid *tview.Grid) {
	_, _, width, _ := grid.GetRect()
	
	// Calculate optimal column widths for 3-column layout
	switch {
	case width < 80: // Narrow terminal
		// Equal distribution for 3 columns
		colWidth := width / 3
		grid.SetColumns(colWidth, colWidth, width-2*colWidth)
	case width < 120: // Medium width terminal
		// Proportional distribution
		col1 := width / 4   // Namespaces
		col2 := width / 3   // Resource Types
		col3 := width - col1 - col2  // Resources
		grid.SetColumns(col1, col2, col3)
	default: // Wide terminal (120+)
		// Optimal distribution for wide screens
		col1 := max(20, width/5)     // Namespaces
		col2 := max(25, width/4)     // Resource Types
		col3 := width - col1 - col2 - 5 // Resources (remaining space)
		grid.SetColumns(col1, col2, col3)
	}
}

// showLogsWindow shows or hides the logs window based on context
func (a *App) showLogsWindow(show bool) {
	if show {
		// Show logs window
		a.grid.AddItem(a.LogsView, 1, 2, 1, 1, 0, 0, false)
	} else {
		// Hide logs window - expand info view to full width
		a.grid.AddItem(a.InfoView, 1, 0, 1, 3, 0, 0, false)
	}
	a.updateGridLayout(a.grid)
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
