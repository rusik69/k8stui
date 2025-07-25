package app

// initResourceTypes initializes the resource type selection
func (a *App) initResourceTypes() {
	// Add resource types to the ResourceTypeList
	resourceTypes := GetAllResourceTypes()
	
	a.ResourceTypeList.Clear()
	
	// Add each resource type
	for _, resourceType := range resourceTypes {
		rt := resourceType // capture for closure
		a.ResourceTypeList.AddItem(GetResourceDisplayName(resourceType), "", 0, func() {
			a.SelectedResourceType = rt
			a.loadSelectedResourceType()
		})
	}
	
	// Add "Pods" as an option to switch back to pod view
	a.ResourceTypeList.AddItem("Pods", "", 0, func() {
		a.SelectedResourceType = ResourceTypePod
		a.LoadPods()
	})
	
	// Add "Namespaces" as an option
	a.ResourceTypeList.AddItem("Namespaces", "", 0, func() {
		a.SelectedResourceType = ResourceTypeNamespace
		a.LoadNamespaces()
	})
}

// loadSelectedResourceType loads the selected resource type
func (a *App) loadSelectedResourceType() {
	switch a.SelectedResourceType {
	case ResourceTypeDeployment:
		a.LoadDeployments()
	case ResourceTypeService:
		a.LoadServices()
	case ResourceTypeConfigMap:
		a.LoadConfigMaps()
	case ResourceTypeSecret:
		a.LoadSecrets()
	case ResourceTypeIngress:
		a.LoadIngresses()
	case ResourceTypeNode:
		a.LoadNodes()
	case ResourceTypePod:
		a.LoadPods()
	case ResourceTypeNamespace:
		a.LoadNamespaces()
	}
}

// updateGridLayoutForResources updates the grid layout based on current view
func (a *App) updateGridLayoutForResources() {
	// Remove all items
	a.grid.Clear()
	
	// Check if we're in resource type selection mode
	if a.SelectedResourceType != "" {
		// Show resource list - 3 column layout
		a.grid.AddItem(a.NsList, 0, 0, 1, 1, 0, 0, true).
			AddItem(a.ResourceTypeList, 0, 1, 1, 1, 0, 0, false).
			AddItem(a.ResourceList, 0, 2, 1, 1, 0, 0, false).
			AddItem(a.InfoView, 1, 0, 1, 2, 0, 0, false).
			AddItem(a.LogsView, 1, 2, 1, 1, 0, 0, false)
	} else {
		// Default layout - 3 column layout
		a.grid.AddItem(a.NsList, 0, 0, 1, 1, 0, 0, true).
			AddItem(a.ResourceTypeList, 0, 1, 1, 1, 0, 0, false).
			AddItem(a.ResourceList, 0, 2, 1, 1, 0, 0, false).
			AddItem(a.InfoView, 1, 0, 1, 2, 0, 0, false).
			AddItem(a.LogsView, 1, 2, 1, 1, 0, 0, false)
	}
	
	a.updateGridLayout(a.grid)
}
