package app

import (
	"fmt"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadReplicaSets loads all ReplicaSets in the current namespace
func (a *App) LoadReplicaSets() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	replicasets, err := a.KubeClient.AppsV1().ReplicaSets(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing replicasets: %v", err)
	}

	a.ResourceList.Clear()
	for _, rs := range replicasets.Items {
		rs := rs // capture for closure
		desired := fmt.Sprintf("%d/%d", rs.Status.ReadyReplicas, *rs.Spec.Replicas)
		
		a.ResourceList.AddItem(rs.Name, desired, 0, func() {
			a.SelectedResource = rs.Name
			a.SelectedResourceType = ResourceTypeReplicaSet
		})
	}

	return nil
}

// LoadStatefulSets loads all StatefulSets in the current namespace
func (a *App) LoadStatefulSets() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	statefulsets, err := a.KubeClient.AppsV1().StatefulSets(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing statefulsets: %v", err)
	}

	a.ResourceList.Clear()
	for _, sts := range statefulsets.Items {
		sts := sts // capture for closure
		desired := fmt.Sprintf("%d/%d", sts.Status.ReadyReplicas, *sts.Spec.Replicas)
		
		a.ResourceList.AddItem(sts.Name, desired, 0, func() {
			a.SelectedResource = sts.Name
			a.SelectedResourceType = ResourceTypeStatefulSet
		})
	}

	return nil
}

// LoadDaemonSets loads all DaemonSets in the current namespace
func (a *App) LoadDaemonSets() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	daemonsets, err := a.KubeClient.AppsV1().DaemonSets(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing daemonsets: %v", err)
	}

	a.ResourceList.Clear()
	for _, ds := range daemonsets.Items {
		ds := ds // capture for closure
		desired := fmt.Sprintf("%d/%d", ds.Status.NumberReady, ds.Status.DesiredNumberScheduled)
		
		a.ResourceList.AddItem(ds.Name, desired, 0, func() {
			a.SelectedResource = ds.Name
			a.SelectedResourceType = ResourceTypeDaemonSet
		})
	}

	return nil
}

// LoadJobs loads all Jobs in the current namespace
func (a *App) LoadJobs() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	jobs, err := a.KubeClient.BatchV1().Jobs(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing jobs: %v", err)
	}

	a.ResourceList.Clear()
	for _, job := range jobs.Items {
		job := job // capture for closure
		status := fmt.Sprintf("%d/%d", job.Status.Succeeded, *job.Spec.Completions)
		
		a.ResourceList.AddItem(job.Name, status, 0, func() {
			a.SelectedResource = job.Name
			a.SelectedResourceType = ResourceTypeJob
		})
	}

	return nil
}

// LoadCronJobs loads all CronJobs in the current namespace
func (a *App) LoadCronJobs() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	cronjobs, err := a.KubeClient.BatchV1().CronJobs(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing cronjobs: %v", err)
	}

	a.ResourceList.Clear()
	for _, cj := range cronjobs.Items {
		cj := cj // capture for closure
		schedule := cj.Spec.Schedule
		
		a.ResourceList.AddItem(cj.Name, schedule, 0, func() {
			a.SelectedResource = cj.Name
			a.SelectedResourceType = ResourceTypeCronJob
		})
	}

	return nil
}

// LoadPVCs loads all PersistentVolumeClaims in the current namespace
func (a *App) LoadPVCs() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	pvcs, err := a.KubeClient.CoreV1().PersistentVolumeClaims(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing pvcs: %v", err)
	}

	a.ResourceList.Clear()
	for _, pvc := range pvcs.Items {
		pvc := pvc // capture for closure
		status := string(pvc.Status.Phase)
		capacity := pvc.Status.Capacity.Storage().String()
		
		info := fmt.Sprintf("%s %s", status, capacity)
		
		a.ResourceList.AddItem(pvc.Name, info, 0, func() {
			a.SelectedResource = pvc.Name
			a.SelectedResourceType = ResourceTypePVC
		})
	}

	return nil
}

// LoadPVs loads all PersistentVolumes in the cluster
func (a *App) LoadPVs() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	pvs, err := a.KubeClient.CoreV1().PersistentVolumes().List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing pvs: %v", err)
	}

	a.ResourceList.Clear()
	for _, pv := range pvs.Items {
		pv := pv // capture for closure
		status := string(pv.Status.Phase)
		capacity := pv.Spec.Capacity.Storage().String()
		
		info := fmt.Sprintf("%s %s", status, capacity)
		
		a.ResourceList.AddItem(pv.Name, info, 0, func() {
			a.SelectedResource = pv.Name
			a.SelectedResourceType = ResourceTypePV
		})
	}

	return nil
}

// LoadNetworkPolicies loads all NetworkPolicies in the current namespace
func (a *App) LoadNetworkPolicies() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	networkpolicies, err := a.KubeClient.NetworkingV1().NetworkPolicies(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing networkpolicies: %v", err)
	}

	a.ResourceList.Clear()
	for _, np := range networkpolicies.Items {
		np := np // capture for closure
		var policyTypes []string
		for _, pt := range np.Spec.PolicyTypes {
			policyTypes = append(policyTypes, string(pt))
		}
		policyTypeStr := strings.Join(policyTypes, ",")
		
		a.ResourceList.AddItem(np.Name, policyTypeStr, 0, func() {
			a.SelectedResource = np.Name
			a.SelectedResourceType = ResourceTypeNetworkPolicy
		})
	}

	return nil
}

// LoadServiceAccounts loads all ServiceAccounts in the current namespace
func (a *App) LoadServiceAccounts() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	serviceaccounts, err := a.KubeClient.CoreV1().ServiceAccounts(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing serviceaccounts: %v", err)
	}

	a.ResourceList.Clear()
	for _, sa := range serviceaccounts.Items {
		sa := sa // capture for closure
		
		a.ResourceList.AddItem(sa.Name, "", 0, func() {
			a.SelectedResource = sa.Name
			a.SelectedResourceType = ResourceTypeServiceAccount
		})
	}

	return nil
}

// LoadRoles loads all Roles in the current namespace
func (a *App) LoadRoles() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	roles, err := a.KubeClient.RbacV1().Roles(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing roles: %v", err)
	}

	a.ResourceList.Clear()
	for _, role := range roles.Items {
		role := role // capture for closure
		rules := fmt.Sprintf("%d rules", len(role.Rules))
		
		a.ResourceList.AddItem(role.Name, rules, 0, func() {
			a.SelectedResource = role.Name
			a.SelectedResourceType = ResourceTypeRole
		})
	}

	return nil
}

// LoadRoleBindings loads all RoleBindings in the current namespace
func (a *App) LoadRoleBindings() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	rolebindings, err := a.KubeClient.RbacV1().RoleBindings(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing rolebindings: %v", err)
	}

	a.ResourceList.Clear()
	for _, rb := range rolebindings.Items {
		rb := rb // capture for closure
		bindings := fmt.Sprintf("%d subjects", len(rb.Subjects))
		
		a.ResourceList.AddItem(rb.Name, bindings, 0, func() {
			a.SelectedResource = rb.Name
			a.SelectedResourceType = ResourceTypeRoleBinding
		})
	}

	return nil
}

// LoadClusterRoles loads all ClusterRoles in the cluster
func (a *App) LoadClusterRoles() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	clusterroles, err := a.KubeClient.RbacV1().ClusterRoles().List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing clusterroles: %v", err)
	}

	a.ResourceList.Clear()
	for _, cr := range clusterroles.Items {
		cr := cr // capture for closure
		rules := fmt.Sprintf("%d rules", len(cr.Rules))
		
		a.ResourceList.AddItem(cr.Name, rules, 0, func() {
			a.SelectedResource = cr.Name
			a.SelectedResourceType = ResourceTypeClusterRole
		})
	}

	return nil
}

// LoadClusterRoleBindings loads all ClusterRoleBindings in the cluster
func (a *App) LoadClusterRoleBindings() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	clusterrolebindings, err := a.KubeClient.RbacV1().ClusterRoleBindings().List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing clusterrolebindings: %v", err)
	}

	a.ResourceList.Clear()
	for _, crb := range clusterrolebindings.Items {
		crb := crb // capture for closure
		bindings := fmt.Sprintf("%d subjects", len(crb.Subjects))
		
		a.ResourceList.AddItem(crb.Name, bindings, 0, func() {
			a.SelectedResource = crb.Name
			a.SelectedResourceType = ResourceTypeClusterRoleBinding
		})
	}

	return nil
}

// LoadEndpoints loads all Endpoints in the current namespace
func (a *App) LoadEndpoints() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	endpoints, err := a.KubeClient.CoreV1().Endpoints(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing endpoints: %v", err)
	}

	a.ResourceList.Clear()
	for _, ep := range endpoints.Items {
		ep := ep // capture for closure
		addresses := 0
		for _, subset := range ep.Subsets {
			addresses += len(subset.Addresses)
		}
		
		info := fmt.Sprintf("%d addresses", addresses)
		
		a.ResourceList.AddItem(ep.Name, info, 0, func() {
			a.SelectedResource = ep.Name
			a.SelectedResourceType = ResourceTypeEndpoint
		})
	}

	return nil
}

// LoadHPAs loads all HorizontalPodAutoscalers in the current namespace
func (a *App) LoadHPAs() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	hpas, err := a.KubeClient.AutoscalingV1().HorizontalPodAutoscalers(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing hpas: %v", err)
	}

	a.ResourceList.Clear()
	for _, hpa := range hpas.Items {
		hpa := hpa // capture for closure
		current := strconv.Itoa(int(hpa.Status.CurrentReplicas))
		min := "0"
		if hpa.Spec.MinReplicas != nil {
			min = strconv.Itoa(int(*hpa.Spec.MinReplicas))
		}
		max := strconv.Itoa(int(hpa.Spec.MaxReplicas))
		
		status := fmt.Sprintf("%s/%s-%s", current, min, max)
		
		a.ResourceList.AddItem(hpa.Name, status, 0, func() {
			a.SelectedResource = hpa.Name
			a.SelectedResourceType = ResourceTypeHPA
		})
	}

	return nil
}

// LoadLimitRanges loads all LimitRanges in the current namespace
func (a *App) LoadLimitRanges() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	limitranges, err := a.KubeClient.CoreV1().LimitRanges(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing limitranges: %v", err)
	}

	a.ResourceList.Clear()
	for _, lr := range limitranges.Items {
		lr := lr // capture for closure
		
		a.ResourceList.AddItem(lr.Name, "", 0, func() {
			a.SelectedResource = lr.Name
			a.SelectedResourceType = ResourceTypeLimitRange
		})
	}

	return nil
}

// LoadResourceQuotas loads all ResourceQuotas in the current namespace
func (a *App) LoadResourceQuotas() error {
	if a.KubeClient == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	if a.CurrentNs == "" {
		return fmt.Errorf("no namespace selected")
	}

	resourcequotas, err := a.KubeClient.CoreV1().ResourceQuotas(a.CurrentNs).List(a.getContext(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing resourcequotas: %v", err)
	}

	a.ResourceList.Clear()
	for _, rq := range resourcequotas.Items {
		rq := rq // capture for closure
		
		a.ResourceList.AddItem(rq.Name, "", 0, func() {
			a.SelectedResource = rq.Name
			a.SelectedResourceType = ResourceTypeResourceQuota
		})
	}

	return nil
}
