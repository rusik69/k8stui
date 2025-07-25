package app

import (
	"github.com/rivo/tview"
)

// showResourceTypeModal displays a modal for selecting resource types
func (a *App) showResourceTypeModal() {
	// Create a modal for resource type selection
	modal := tview.NewModal().
		SetText("Select Resource Type").
		AddButtons([]string{
			"Deployments",
			"ReplicaSets",
			"StatefulSets",
			"DaemonSets",
			"Jobs",
			"CronJobs",
			"Services",
			"ConfigMaps",
			"Secrets",
			"Ingresses",
			"NetworkPolicies",
			"PVCs",
			"PVs",
			"ServiceAccounts",
			"Roles",
			"RoleBindings",
			"ClusterRoles",
			"ClusterRoleBindings",
			"Endpoints",
			"HPAs",
			"LimitRanges",
			"ResourceQuotas",
			"Nodes",
			"Pods",
			"Namespaces",
			"Cancel",
		}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.HidePage("resource_modal")
			
			switch buttonLabel {
			case "Deployments":
				a.SelectedResourceType = ResourceTypeDeployment
				a.LoadDeployments()
			case "Services":
				a.SelectedResourceType = ResourceTypeService
				a.LoadServices()
			case "ConfigMaps":
				a.SelectedResourceType = ResourceTypeConfigMap
				a.LoadConfigMaps()
			case "Secrets":
				a.SelectedResourceType = ResourceTypeSecret
				a.LoadSecrets()
			case "Ingresses":
				a.SelectedResourceType = ResourceTypeIngress
				a.LoadIngresses()
			case "Nodes":
				a.SelectedResourceType = ResourceTypeNode
				a.LoadNodes()
			case "Pods":
				a.SelectedResourceType = ResourceTypePod
				a.LoadPods()
			case "Namespaces":
				a.SelectedResourceType = ResourceTypeNamespace
				a.LoadNamespaces()
			}
		})

	// Add the modal to pages
	a.pages.AddPage("resource_modal", modal, true, true)
	
	// Set focus to the modal
	a.App.SetFocus(modal)
}
