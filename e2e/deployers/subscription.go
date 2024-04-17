package deployers

import (
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type Subscription struct {
	// branch  string
	// path    string
	// channel string
	Ctx *util.Context

	ChannelName      string
	ChannelNamespace string
	SubscriptionName string

	Name      string // deployment-rbd
	Namespace string // deployment-rbd
}

func (s *Subscription) Init() {
	s.ChannelName = "ramen-gitops"
	s.ChannelNamespace = "ramen-samples"
	s.SubscriptionName = "subscription"
	s.Name = "deployment-rbd"
	s.Namespace = "deployment-rbd"
}

func (s Subscription) GetName() string {
	return s.Name
}

func (s Subscription) GetNameSpace() string {
	return s.Namespace
}

func (s Subscription) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Use the global Channel
	// Generate a Binding for the namespace (does this need clusters?)
	// Generate a Subscription for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	util.LogEnter(s.Ctx.Log)
	defer util.LogExit(s.Ctx.Log)

	// w.Kustomize()
	err := createNamespace(s.Ctx, s.Namespace)
	if err != nil {
		return err
	}

	err = createManagedClusterSetBinding(s.Ctx, "default", s.Namespace, s.Name)
	if err != nil {
		return err
	}

	err = createPlacement(s.Ctx, util.DefaultPlacement, s.Namespace, s.Name)
	if err != nil {
		return err
	}

	err = s.createSubscription(w)
	if err != nil {
		return err
	}

	return nil
}

func (s Subscription) Undeploy(w workloads.Workload) error {
	// Delete Subscription, Placement, Binding
	util.LogEnter(s.Ctx.Log)
	defer util.LogExit(s.Ctx.Log)

	err := s.deleteSubscription()
	if err != nil {
		return err
	}

	err = deletePlacement(s.Ctx, util.DefaultPlacement, s.Namespace)
	if err != nil {
		return err
	}

	err = deleteManagedClusterSetBinding(s.Ctx, "default", s.Namespace)
	if err != nil {
		return err
	}

	err = deleteNamespace(s.Ctx, s.Namespace)
	if err != nil {
		return err
	}

	return nil
}

func (s Subscription) Health(w workloads.Workload) error {
	util.LogEnter(s.Ctx.Log)
	defer util.LogExit(s.Ctx.Log)
	// w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
