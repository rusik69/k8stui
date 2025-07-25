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
	grid := tview.NewGrid()
	hotkeyHelp := "[::b]TAB/Shift+TAB[::-] Navigate | [::b]ENTER[::-] Select | [::b]Ctrl+D[::-] Delete | [::b]Q[::-] Quit | [::b]↑/↓/←/→[::-] Scroll"
	grid.SetBorder(true).SetTitle(" K8s TUI - " + hotkeyHelp + " ")

	// Set up the grid layout to be responsive to terminal size
	grid.SetRows(0, 0). // Two rows of equal height
		SetBorders(true)

	// Add items with dynamic sizing based on terminal width
	grid.AddItem(a.NsList, 0, 0, 1, 1, 0, 0, true).
		AddItem(a.PodList, 0, 1, 1, 1, 0, 0, false).
		AddItem(a.ContList, 0, 2, 1, 1, 0, 0, false).
		AddItem(a.InfoView, 1, 0, 1, 2, 0, 0, false).
		AddItem(a.LogsView, 1, 2, 1, 1, 0, 0, false)

	// Set up dynamic column sizes based on terminal width
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		_, _, width, _ := grid.GetRect()
		
		// Calculate column widths based on terminal width
		switch {
		case width < 80: // Very narrow terminal
			grid.SetColumns(15, 15, 20)
		case width < 120: // Medium width terminal
			grid.SetColumns(20, 25, 35)
		default: // Wide terminal
			grid.SetColumns(25, 30, 40)
		}
		
		return event
	})

	// Create a pages container that will hold our main UI and modals
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(grid, 0, 1, true)

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




