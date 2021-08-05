/*
Copyright 2021 The RamenDR authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	"github.com/ramendr/ramen/controllers/util"
)

// DRPolicyReconciler reconciles a DRPolicy object
type DRPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ramendr.openshift.io,resources=drpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ramendr.openshift.io,resources=drpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ramendr.openshift.io,resources=drpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DRPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *DRPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithName("controllers").WithName("drpolicy").WithValues("name", req.NamespacedName.Name)
	log.Info("reconcile enter")

	defer log.Info("reconcile exit")

	drpolicy := &ramen.DRPolicy{}
	if err := r.Client.Get(ctx, req.NamespacedName, drpolicy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(fmt.Errorf("get: %w", err))
	}

	manifestWorkUtil := util.MWUtil{Client: r.Client, Ctx: ctx, Log: log, InstName: "", InstNamespace: ""}

	switch drpolicy.ObjectMeta.DeletionTimestamp.IsZero() {
	case true:
		log.Info("create/update")

		if err := finalizerAdd(ctx, drpolicy, r.Client, log); err != nil {
			return ctrl.Result{}, fmt.Errorf("finalizer add update: %w", err)
		}

		if err := manifestWorkUtil.ClusterRolesCreate(drpolicy); err != nil {
			return ctrl.Result{}, fmt.Errorf("cluster roles create: %w", err)
		}
	default:
		log.Info("delete")

		if err := manifestWorkUtil.ClusterRolesDelete(drpolicy); err != nil {
			return ctrl.Result{}, fmt.Errorf("cluster roles delete: %w", err)
		}

		if err := finalizerRemove(ctx, drpolicy, r.Client, log); err != nil {
			return ctrl.Result{}, fmt.Errorf("finalizer remove update: %w", err)
		}
	}

	return ctrl.Result{}, nil
}

const finalizerName = "drpolicies.ramendr.openshift.io/ramen"

func finalizerAdd(ctx context.Context, drpolicy *ramen.DRPolicy, client client.Client, log logr.Logger) error {
	finalizerCount := len(drpolicy.ObjectMeta.Finalizers)
	controllerutil.AddFinalizer(drpolicy, finalizerName)

	if len(drpolicy.ObjectMeta.Finalizers) != finalizerCount {
		log.Info("finalizer add")

		return client.Update(ctx, drpolicy)
	}

	return nil
}

func finalizerRemove(ctx context.Context, drpolicy *ramen.DRPolicy, client client.Client, log logr.Logger) error {
	finalizerCount := len(drpolicy.ObjectMeta.Finalizers)
	controllerutil.RemoveFinalizer(drpolicy, finalizerName)

	if len(drpolicy.ObjectMeta.Finalizers) != finalizerCount {
		log.Info("finalizer remove")

		return client.Update(ctx, drpolicy)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DRPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ramen.DRPolicy{}).
		Complete(r)
}
