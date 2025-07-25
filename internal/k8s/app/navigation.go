package app


// navigate handles keyboard navigation between UI elements
func (a *App) navigate(forward bool) {
	if forward {
		a.CurrentFocus = (a.CurrentFocus + 1) % 3
	} else {
		a.CurrentFocus = (a.CurrentFocus + 2) % 3 // +2 is equivalent to -1 mod 3
	}
	a.UpdateFocus()
}

// UpdateFocus updates the focus based on CurrentFocus
func (a *App) UpdateFocus() {
	switch a.CurrentFocus {
	case 0:
		a.App.SetFocus(a.NsList)
	case 1:
		a.App.SetFocus(a.PodList)
	case 2:
		a.App.SetFocus(a.ContList)
	}
}
