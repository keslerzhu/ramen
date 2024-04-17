package dractions

import (
	"fmt"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DRActions struct {
	Ctx *util.Context
}

const OcmSchedulingDisable = "cluster.open-cluster-management.io/experimental-scheduling-disable"

func (r DRActions) EnableProtection(w workloads.Workload, d deployers.Deployer) error {
	// If AppSet/Subscription, find Placement
	// Determine DRPolicy
	// Determine preferredCluster
	// Determine PVC label selector
	// Determine KubeObjectProtection requirements if Imperative (?)
	// Create DRPC, in desired namespace
	util.LogEnter(r.Ctx.Log)
	defer util.LogExit(r.Ctx.Log)

	_, isSub := d.(*deployers.Subscription)
	_, isAppSet := d.(*deployers.ApplicationSet)
	if isSub || isAppSet {
		name := d.GetName()
		namespace := d.GetNameSpace()
		drPolicyName := r.Ctx.Config.DRPolicy
		appname := w.GetAppName()

		//TODO: improve placement name
		placementName := util.DefaultPlacement
		if isAppSet {
			placementName = d.GetName() + "-placement"
		}

		drpcName := name + "-drpc"
		client := r.Ctx.Hub.CtrlClient

		placement, placementDecisionName, err := r.waitPlacementDecision(client, namespace, placementName)
		if err != nil {
			return err
		}

		r.Ctx.Log.Info("get placementdecision " + placementDecisionName)

		placementDecision, err := getPlacementDecision(client, namespace, placementDecisionName)
		if err != nil {
			return err
		}

		clusterName := placementDecision.Status.Decisions[0].ClusterName
		r.Ctx.Log.Info("placementdecision clusterName: " + clusterName)

		// move update placement annotation after placement has been handled
		// otherwise if we first add ocm disable annotation then it might not
		// yet be handled by ocm and thus PlacementSatisfied=false

		if placement.Annotations == nil {
			placement.Annotations = make(map[string]string)
		}

		placement.Annotations[OcmSchedulingDisable] = "true"

		r.Ctx.Log.Info("update placement " + placementName + " annotation")

		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}

		r.Ctx.Log.Info("create drpc " + drpcName)
		drpc := &ramen.DRPlacementControl{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DRPlacementControl",
				APIVersion: "ramendr.openshift.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      drpcName,
				Namespace: namespace,
				Labels:    map[string]string{"app": name},
			},
			Spec: ramen.DRPlacementControlSpec{
				PreferredCluster: clusterName,
				DRPolicyRef: v1.ObjectReference{
					Name: drPolicyName,
				},
				PlacementRef: v1.ObjectReference{
					Kind: "placement",
					Name: placementName,
				},
				PVCSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"appname": appname},
				},
			},
		}

		err = r.createDRPC(client, drpc)
		if err != nil {
			return err
		}

		err = r.waitDRPCReady(client, namespace, drpcName)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("deployer not known")
	}

	return nil
}

func (r DRActions) DisableProtection(w workloads.Workload, d deployers.Deployer) error {
	// remove DRPC
	// update placement annotation
	util.LogEnter(r.Ctx.Log)
	defer util.LogExit(r.Ctx.Log)

	_, ok := d.(*deployers.Subscription)
	if ok {
		name := d.GetName()
		namespace := d.GetNameSpace()
		placementName := util.DefaultPlacement
		drpcName := name + "-drpc"
		client := r.Ctx.Hub.CtrlClient

		r.Ctx.Log.Info("delete drpc " + drpcName)

		err := deleteDRPC(client, namespace, drpcName)
		if err != nil {
			return err
		}

		r.Ctx.Log.Info("get placement " + placementName)

		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {
			return err
		}

		delete(placement.Annotations, OcmSchedulingDisable)

		r.Ctx.Log.Info("update placement " + placementName + " annotation")

		err = updatePlacement(client, placement)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("deployer not known")
	}

	return nil
}

func (r DRActions) Failover(w workloads.Workload, d deployers.Deployer) error {
	// Determine DRPC
	// Check Placement
	// Failover to alternate in DRPolicy as the failoverCluster
	// Update DRPC
	util.LogEnter(r.Ctx.Log)
	defer util.LogExit(r.Ctx.Log)

	name := d.GetName()
	namespace := d.GetNameSpace()
	// placementName := w.GetPlacementName()
	drPolicyName := r.Ctx.Config.DRPolicy
	drpcName := name + "-drpc"
	client := r.Ctx.Hub.CtrlClient

	// here we expect drpc should be ready before failover
	err := r.waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	// enable phase check when necessary
	// here we expect phase should be Deployed before failover
	// TODO: will update for other valid phases
	err = r.waitDRPCPhase(client, namespace, drpcName, "Deployed")
	if err != nil {
		return err
	}

	r.Ctx.Log.Info("get drpc " + drpcName)

	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	r.Ctx.Log.Info("get drpolicy " + drPolicyName)

	drpolicy, err := getDRPolicy(client, drPolicyName)
	if err != nil {
		return err
	}

	preferredCluster := drpc.Spec.PreferredCluster
	failoverCluster := ""

	if preferredCluster == drpolicy.Spec.DRClusters[0] {
		failoverCluster = drpolicy.Spec.DRClusters[1]
	} else {
		failoverCluster = drpolicy.Spec.DRClusters[0]
	}

	r.Ctx.Log.Info("preferredCluster: " + preferredCluster + " -> failoverCluster: " + failoverCluster)

	drpc.Spec.Action = "Failover"
	drpc.Spec.FailoverCluster = failoverCluster

	r.Ctx.Log.Info("update drpc " + drpcName)

	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	// check Phase
	err = r.waitDRPCPhase(client, namespace, drpcName, "FailedOver")
	if err != nil {
		return err
	}
	// then check Conditions
	err = r.waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	return nil
}

func (r DRActions) Relocate(w workloads.Workload, d deployers.Deployer) error {
	// Determine DRPC
	// Check Placement
	// Relocate to Primary in DRPolicy as the PrimaryCluster
	// Update DRPC
	util.LogEnter(r.Ctx.Log)
	defer util.LogExit(r.Ctx.Log)

	name := d.GetName()
	namespace := d.GetNameSpace()
	// placementName := w.GetPlacementName()
	drpcName := name + "-drpc"
	client := r.Ctx.Hub.CtrlClient

	// here we expect drpc should be ready before relocate
	err := r.waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	// enable phase check when necessary
	// here we expect phase should be FailedOver before relocate
	// TODO: will update for other valid phases
	err = r.waitDRPCPhase(client, namespace, drpcName, "FailedOver")
	if err != nil {
		return err
	}

	r.Ctx.Log.Info("get drpc " + drpcName)

	drpc, err := getDRPC(client, namespace, drpcName)
	if err != nil {
		return err
	}

	drpc.Spec.Action = "Relocate"

	r.Ctx.Log.Info("update drpc " + drpcName)

	err = updateDRPC(client, drpc)
	if err != nil {
		return err
	}

	// check Phase
	err = r.waitDRPCPhase(client, namespace, drpcName, "Relocated")
	if err != nil {
		return err
	}
	// then check Conditions
	err = r.waitDRPCReady(client, namespace, drpcName)
	if err != nil {
		return err
	}

	return nil
}
