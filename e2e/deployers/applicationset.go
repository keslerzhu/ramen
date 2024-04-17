package deployers

import (
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type ApplicationSet struct {
	Ctx *util.Context

	Name      string
	Namespace string

	// ArgoCDNamespace                 string
	PlacementName                   string
	McsbName                        string
	ClusterDecisionConfigMapName    string
	ApplicationDestinationNamespace string
}

func (a *ApplicationSet) Init() {
	a.Name = "busybox-appset"
	a.Namespace = "argocd"
	// a.ArgoCDNamespace = "argocd"
	a.PlacementName = a.Name + "-placement"
	a.McsbName = "default"
	a.ClusterDecisionConfigMapName = a.Name + "-configmap"
	a.ApplicationDestinationNamespace = a.Name
}

func (a ApplicationSet) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Generate a Binding for the namespace?
	// Generate an ApplicationSet for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	// w.Kustomize()
	util.LogEnter(a.Ctx.Log)
	defer util.LogExit(a.Ctx.Log)

	err := a.addArgoCDClusters()
	if err != nil {
		return err
	}

	err = createManagedClusterSetBinding(a.Ctx, a.McsbName, a.Namespace, a.Name)
	if err != nil {
		return err
	}

	err = createPlacement(a.Ctx, a.PlacementName, a.Namespace, a.Name)
	if err != nil {
		return err
	}

	err = createPlacementDecisionConfigMap(a.Ctx, a.ClusterDecisionConfigMapName, a.Namespace)
	if err != nil {
		return err
	}

	err = a.createApplicationSet(w)
	if err != nil {
		return err
	}

	return err
}

func (a ApplicationSet) Undeploy(w workloads.Workload) error {
	// Delete Placement, Binding, ApplicationSet
	util.LogEnter(a.Ctx.Log)
	defer util.LogExit(a.Ctx.Log)

	err := a.deleteApplicationSet()
	if err != nil {
		return err
	}

	err = deleteConfigMap(a.Ctx, a.ClusterDecisionConfigMapName, a.Namespace)
	if err != nil {
		return err
	}

	err = deletePlacement(a.Ctx, a.PlacementName, a.Namespace)
	if err != nil {
		return err
	}

	err = deleteManagedClusterSetBinding(a.Ctx, a.McsbName, a.Namespace)
	if err != nil {
		return err
	}

	// don't use, this function is problematic
	//
	// 2024-04-15T13:07:21.319+0800    ERROR   util/cmd.go:22  ====== cmd start ======
	// cmd error: /usr/local/bin/argocd cluster rm rdr-hub  -y
	// time="2024-04-15T13:07:21+08:00" level=fatal msg=EOF

	// ====== cmd end ======   {"error": "exit status 20"}
	//
	// err = a.deleteArgoCDClusters()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a ApplicationSet) GetName() string {
	return a.Name
}

func (a ApplicationSet) GetNameSpace() string {
	return a.Namespace
}

func (a ApplicationSet) Health(w workloads.Workload) error {
	util.LogEnter(a.Ctx.Log)
	defer util.LogExit(a.Ctx.Log)

	// w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
