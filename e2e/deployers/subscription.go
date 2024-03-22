package deployers

import (
	"os/exec"

	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type Subscription struct {
	branch  string
	path    string
	channel string
	Ctx     *util.TestContext
}

func (s Subscription) Deploy(w workloads.Workload) error {
	// Generate a Placement for the Workload
	// Use the global Channel
	// Generate a Binding for the namespace (does this need clusters?)
	// Generate a Subscription for the Workload
	// - Kustomize the Workload; call Workload.Kustomize(StorageType)
	// Address namespace/label/suffix as needed for various resources
	s.Ctx.Log.Info("enter Subscription Deploy")
	// w.Kustomize()

	cmd := exec.Command("kubectl", "apply", "-k", w.GetResourceURL(), "--kubeconfig="+s.Ctx.HubKubeconfig())
	err, _ := util.RunCommand(cmd)
	return err
}

func (s Subscription) Undeploy(w workloads.Workload) error {
	// Delete Subscription, Placement, Binding
	s.Ctx.Log.Info("enter Subscription Undeploy")

	cmd := exec.Command("kubectl", "delete", "-k", w.GetResourceURL(), "--kubeconfig="+s.Ctx.HubKubeconfig())
	err, _ := util.RunCommand(cmd)
	return err
}

func (s Subscription) Health(w workloads.Workload) error {
	s.Ctx.Log.Info("enter Subscription Health")
	w.GetResources()
	// Check health using reflection to known types of the workload on the targetCluster
	// Again if using reflection can be a common function outside of deployer as such
	return nil
}
