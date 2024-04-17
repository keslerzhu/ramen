package deployers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	placementrulev1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	subscriptionv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"

	argocdv1alpha1hack "github.com/ramendr/ramen/e2e/argocd"
)

func (a *ApplicationSet) createApplicationSet(w workloads.Workload) error {
	util.LogEnter(a.Ctx.Log)
	defer util.LogExit(a.Ctx.Log)

	var requeueSeconds int64 = 180

	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
		Spec: argocdv1alpha1hack.ApplicationSetSpec{
			Generators: []argocdv1alpha1hack.ApplicationSetGenerator{
				{
					ClusterDecisionResource: &argocdv1alpha1hack.DuckTypeGenerator{
						ConfigMapRef: a.ClusterDecisionConfigMapName,
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"cluster.open-cluster-management.io/placement": a.PlacementName,
							},
						},
						RequeueAfterSeconds: &requeueSeconds,
					},
				},
			},
			Template: argocdv1alpha1hack.ApplicationSetTemplate{
				ApplicationSetTemplateMeta: argocdv1alpha1hack.ApplicationSetTemplateMeta{
					Name: "rbd-{{name}}",
				},
				Spec: argocdv1alpha1hack.ApplicationSpec{
					Source: &argocdv1alpha1hack.ApplicationSource{
						RepoURL:        a.Ctx.Config.Github.Repo, // w.GetRepoURL(),
						Path:           w.GetPath(),
						TargetRevision: a.Ctx.Config.Github.Branch, // w.GetRevision(),
					},
					Destination: argocdv1alpha1hack.ApplicationDestination{
						Server:    "{{server}}",
						Namespace: a.ApplicationDestinationNamespace,
					},
					Project: "default",
					SyncPolicy: &argocdv1alpha1hack.SyncPolicy{
						Automated: &argocdv1alpha1hack.SyncPolicyAutomated{
							Prune:    true,
							SelfHeal: true,
						},
						SyncOptions: []string{
							"CreateNamespace=true",
							"PruneLast=true",
						},
					},
				},
			},
		},
	}

	err := a.Ctx.Hub.CtrlClient.Create(context.Background(), appset)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			a.Ctx.Log.Error(err, "")

			return err
		}

		a.Ctx.Log.Info("applicationset " + appset.Name + " already Exists")
	}

	return nil
}

func (a *ApplicationSet) deleteApplicationSet() error {
	util.LogEnter(a.Ctx.Log)
	defer util.LogExit(a.Ctx.Log)

	appset := &argocdv1alpha1hack.ApplicationSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
	}

	err := a.Ctx.Hub.CtrlClient.Delete(context.Background(), appset)
	if err != nil {
		if !errors.IsNotFound(err) {
			a.Ctx.Log.Error(err, "")

			return err
		}

		a.Ctx.Log.Info("applicationset " + appset.Name + " not found")
	}

	return nil
}

func createPlacementDecisionConfigMap(ctx *util.Context, cmName string, cmNamespace string) error {
	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	data := map[string]string{
		"apiVersion":    "cluster.open-cluster-management.io/v1beta1",
		"kind":          "placementdecisions",
		"statusListKey": "decisions",
		"matchKey":      "clusterName",
	}

	configMap := &corev1.ConfigMap{ObjectMeta: object, Data: data}

	err := ctx.Hub.CtrlClient.Create(context.Background(), configMap)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			ctx.Log.Error(err, "")

			return fmt.Errorf("could not create configMap " + cmName)
		}

		ctx.Log.Info("configMap " + cmName + " already Exists")
	}

	return nil
}

func deleteConfigMap(ctx *util.Context, cmName string, cmNamespace string) error {
	object := metav1.ObjectMeta{Name: cmName, Namespace: cmNamespace}

	configMap := &corev1.ConfigMap{
		ObjectMeta: object,
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), configMap)
	if err != nil {
		if !errors.IsNotFound(err) {
			ctx.Log.Error(err, "")

			return fmt.Errorf("could not delete configMap " + cmName)
		}

		ctx.Log.Info("configMap " + cmName + " not found")
	}

	return nil
}

func (a *ApplicationSet) addArgoCDClusters() error {
	//TODO: clusternames better to be dynamically got from config
	// for _, c := range util.ClusterNames {
	// 	a.Ctx.Log.Info("add cluster " + c + " into ArgoCD")
	// 	cmd := exec.Command("argocd", "cluster", "add", c, " -y --namespace argocd --kubeconfig "+a.Ctx.Hub.Kubeconfig)

	// 	out, err := util.RunCommand(a.Ctx.Log, cmd)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	a.Ctx.Log.Info(out)
	// }

	return nil
}

// func (a *ApplicationSet) deleteArgoCDClusters() error {
// 	for _, c := range util.ClusterNames {
// 		a.Ctx.Log.Info("delete cluster " + c + " from ArgoCD")
// 		cmd := exec.Command("argocd", "cluster", "rm", c, " -y")

// 		out, err := util.RunCommand(&a.Ctx.Log, cmd)
// 		if err != nil {
// 			return err
// 		}
// 		a.Ctx.Log.Info(out)
// 	}

// 	return nil
// }

func createNamespace(ctx *util.Context, namespace string) error {
	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objNs)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("namespace " + namespace + " already Exists")
	}

	return nil
}

func (s *Subscription) createSubscription(w workloads.Workload) error {
	labels := make(map[string]string)
	labels["apps"] = s.Name

	annotations := make(map[string]string)
	annotations["apps.open-cluster-management.io/github-branch"] = s.Ctx.Config.Github.Branch // w.GetRevision()
	annotations["apps.open-cluster-management.io/github-path"] = w.GetPath()

	// Construct PlacementRef
	objReplacementRef := corev1.ObjectReference{
		Kind: "Placement",
		Name: util.DefaultPlacement,
	}

	objPlacementRulePlacement := &placementrulev1.Placement{}
	objPlacementRulePlacement.PlacementRef = &objReplacementRef

	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:        s.SubscriptionName,
			Namespace:   s.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: subscriptionv1.SubscriptionSpec{
			Channel:   s.ChannelNamespace + "/" + s.ChannelName,
			Placement: objPlacementRulePlacement,
		},
	}

	err := s.Ctx.Hub.CtrlClient.Create(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			s.Ctx.Log.Error(err, "")

			return err
		}

		s.Ctx.Log.Info("placement " + objSubscription.Name + " already Exists")
	}

	return nil
}

func createPlacement(ctx *util.Context, plName string, plNamespace string, appName string) error {
	labels := make(map[string]string)
	labels["apps"] = appName

	arrayClusterSets := []string{"default"}

	var numClusters int32 = 1

	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, // s.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta1.PlacementSpec{
			ClusterSets:      arrayClusterSets,
			NumberOfClusters: &numClusters,
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("placement " + objPlacement.Name + " already Exists")
	}

	return nil
}

func createManagedClusterSetBinding(ctx *util.Context, mcsbName string, mcsbNamespace string, appName string) error {

	labels := make(map[string]string)
	labels["apps"] = appName // s.AppName

	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // s.Namespace,
			Labels:    labels,
		},
		Spec: ocmclusterv1beta2.ManagedClusterSetBindingSpec{
			ClusterSet: "default",
		},
	}

	err := ctx.Hub.CtrlClient.Create(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("managedClusterSetBinding " + objMCSB.Name + " already Exists")
	}

	return nil
}

func deleteNamespace(ctx *util.Context, namespace string) error {
	objNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objNs)
	if err != nil {
		if !errors.IsNotFound(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("namespace " + namespace + " not found")
	}

	return nil
}

func (s *Subscription) deleteSubscription() error {
	objSubscription := &subscriptionv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.SubscriptionName,
			Namespace: s.Namespace,
		},
	}

	err := s.Ctx.Hub.CtrlClient.Delete(context.Background(), objSubscription)
	if err != nil {
		if !errors.IsNotFound(err) {
			s.Ctx.Log.Error(err, "")

			return err
		}

		s.Ctx.Log.Info("subscription " + s.SubscriptionName + " not found")
	}

	return nil
}

func deletePlacement(ctx *util.Context, plName string, plNamespace string) error {
	objPlacement := &ocmclusterv1beta1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      plName,      // util.DefaultPlacement,
			Namespace: plNamespace, // s.Namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objPlacement)
	if err != nil {
		if !errors.IsNotFound(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("placement " + plName + " not found")
	}

	return nil
}

func deleteManagedClusterSetBinding(ctx *util.Context, mcsbName string, mcsbNamespace string) error {
	objMCSB := &ocmclusterv1beta2.ManagedClusterSetBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcsbName,      // "default",
			Namespace: mcsbNamespace, // s.Namespace,
		},
	}

	err := ctx.Hub.CtrlClient.Delete(context.Background(), objMCSB)
	if err != nil {
		if !errors.IsNotFound(err) {
			ctx.Log.Error(err, "")

			return err
		}

		ctx.Log.Info("managedClusterSetBinding default not found")
	}

	return nil
}

/*
func (s *Subscription) createChannel() error {
	util.LogEnter(&s.Ctx.Log)
	defer util.LogExit(&s.Ctx.Log)

	objChannel := &channelv1.Channel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      w.ChannelName,
			Namespace: w.ChannelNamespace,
		},
		Spec: channelv1.ChannelSpec{
			Pathname: w.RepoURL,
			Type:     w.ChannelType,
		},
	}

	err := s.Ctx.HubCtrlClient().Create(context.Background(), objChannel)
	if err != nil {
		s.Ctx.Log.Error(err, "")
		return err
	}

	return nil
}

func (s *Subscription) deleteChannel() error {
	util.LogEnter(&s.Ctx.Log)
	defer util.LogExit(&s.Ctx.Log)

	objChannel := &channelv1.Channel{}

	key := types.NamespacedName{
		Name:      w.ChannelName,
		Namespace: w.ChannelNamespace,
	}

	err := s.Ctx.HubCtrlClient().Get(context.Background(), key, objChannel)
	if err != nil {
		if !errors.IsNotFound(err) {
			s.Ctx.Log.Error(err, "")
			return err
		}
		s.Ctx.Log.Info("channel " + w.ChannelName + " not found")
		return nil
	}

	err = s.Ctx.HubCtrlClient().Delete(context.Background(), objChannel)
	if err != nil {
		s.Ctx.Log.Error(err, "")
		return err
	}

	return nil
}
*/
