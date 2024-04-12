package dractions

import (
	"fmt"
	"time"

	"open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r DRActions) waitPlacementDecision(client client.Client, namespace string, placementName string) (*v1beta1.Placement, string, error) {

	timeout := r.Ctx.Config.Timeout
	interval := r.Ctx.Config.Interval
	startTime := time.Now()
	placementDecisionName := ""

	for {
		placement, err := getPlacement(client, namespace, placementName)
		if err != nil {
			return nil, "", err
		}
		for _, cond := range placement.Status.Conditions {
			if cond.Type == "PlacementSatisfied" && cond.Status == "True" {
				placementDecisionName = placement.Status.DecisionGroups[0].Decisions[0]
				if placementDecisionName != "" {
					r.Ctx.Log.Info("got placementdecision name " + placementDecisionName)
					return placement, placementDecisionName, nil
				}
			}
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			return nil, "", fmt.Errorf("could not get placement decision before timeout")
		}
		r.Ctx.Log.Info(fmt.Sprintf("could not get placement decision, retry in %v seconds", interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func (r DRActions) waitDRPCReady(client client.Client, namespace string, drpcName string) error {

	timeout := r.Ctx.Config.Timeout
	interval := r.Ctx.Config.Interval
	startTime := time.Now()
	for {
		ready := true
		drpc, err := getDRPC(client, namespace, drpcName)
		if err != nil {
			return err
		}

		for _, cond := range drpc.Status.Conditions {
			if cond.Type == "Available" && cond.Status != "True" {
				r.Ctx.Log.Info("drpc status Available is not True")
				ready = false
				break
			}
			if cond.Type == "PeerReady" && cond.Status != "True" {
				r.Ctx.Log.Info("drpc status PeerReady is not True")
				ready = false
				break
			}
		}
		if ready {
			if drpc.Status.LastGroupSyncTime == nil {
				r.Ctx.Log.Info("drpc status LastGroupSyncTime is nil")
				ready = false
			}
		}
		if ready {
			r.Ctx.Log.Info("drpc status is ready")
			return nil
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			return fmt.Errorf(fmt.Sprintf("drpc status is not ready yet before timeout of %v", timeout))
		}
		r.Ctx.Log.Info(fmt.Sprintf("drpc status is not ready yet, retry in %v seconds", interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func (r DRActions) waitDRPCPhase(client client.Client, namespace string, drpcName string, phase string) error {

	timeout := r.Ctx.Config.Timeout
	interval := r.Ctx.Config.Interval
	startTime := time.Now()
	for {
		drpc, err := getDRPC(client, namespace, drpcName)
		if err != nil {
			return err
		}
		currentPhase := string(drpc.Status.Phase)
		if currentPhase == phase {
			r.Ctx.Log.Info("drpc phase is " + phase)
			return nil
		}
		if time.Since(startTime) > time.Second*time.Duration(timeout) {
			return fmt.Errorf(fmt.Sprintf("drpc status is not %s yet before timeout of %v", phase, timeout))
		}
		r.Ctx.Log.Info(fmt.Sprintf("current drpc phase is %s, expecting %s, retry in %v seconds", currentPhase, phase, interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}
