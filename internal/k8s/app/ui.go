package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// initUI initializes the user interface
func (a *App) initUI() {
	// Configure the views
	a.NsList.SetBorder(true).SetTitle(" Namespaces ")
	a.PodList.SetBorder(true).SetTitle(" Pods ")
	a.ContList.SetBorder(true).SetTitle(" Containers ")
	// Configure InfoView with scrolling
	a.InfoView.SetBorder(true).SetTitle(" Pod Info ")
	a.InfoView.SetChangedFunc(func() {
		a.App.Draw()
	})

	// Configure LogsView with auto-scrolling
	a.LogsView.SetBorder(true).SetTitle(" Logs ")
	a.LogsView.SetChangedFunc(func() {
		a.LogsView.ScrollToEnd()
		a.App.Draw()
	})

	// Create the main grid layout
	a.grid = tview.NewGrid()
	hotkeyHelp := "[::b]TAB/Shift+TAB[::-] Navigate | [::b]ENTER[::-] Select | [::b]Ctrl+D[::-] Delete | [::b]Q[::-] Quit | [::b]↑/↓/←/→[::-] Scroll"
	a.grid.SetBorder(true).SetTitle(" K8s TUI - " + hotkeyHelp + " ")

	// Set up the grid layout to be responsive to terminal size
	a.grid.SetRows(0, 0). // Two rows of equal height
		SetBorders(true)

	// Add items with dynamic sizing based on terminal width
	a.grid.AddItem(a.NsList, 0, 0, 1, 1, 0, 0, true).
		AddItem(a.PodList, 0, 1, 1, 1, 0, 0, false).
		AddItem(a.ContList, 0, 2, 1, 1, 0, 0, false).
		AddItem(a.InfoView, 1, 0, 1, 2, 0, 0, false).
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
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				a.App.Stop()
				return nil
			}
		}

		return event
	})
}

// updateGridLayout updates the grid column layout to utilize full terminal width
func (a *App) updateGridLayout(grid *tview.Grid) {
	_, _, width, _ := grid.GetRect()
	
	// Calculate optimal column widths to utilize full terminal width
	// with reasonable minimums and proportional distribution
	switch {
	case width < 80: // Very narrow terminal
		// Use full width with reasonable proportions
		grid.SetColumns(width/3, width/3, width-(2*width/3))
	case width < 120: // Medium width terminal
		// Better proportional distribution
		grid.SetColumns(width/4, width/3, width-(width/4+width/3))
	default: // Wide terminal (120+)
		// Optimal distribution for wide screens
		col1 := max(25, width/5)          // Namespace list
		col2 := max(30, width/4)         // Pod list  
		col3 := width - col1 - col2 - 5  // Remaining space for info/logs
		grid.SetColumns(col1, col2, col3)
	}
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
