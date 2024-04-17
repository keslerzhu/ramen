package suites

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/ramendr/ramen/e2e/util"
// 	rookv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// )

// const ramenSystemNamespace = "ramen-system"

// type PrecheckSuite struct {
// 	Ctx *util.Context
// }

// func (s *PrecheckSuite) SetContext(ctx *util.Context) {
// 	util.LogEnter(&ctx.Log)
// 	defer util.LogExit(&ctx.Log)

// 	s.Ctx = ctx
// }

// func (s *PrecheckSuite) SetupSuite() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	return nil
// }

// func (s *PrecheckSuite) TeardownSuite() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	return nil
// }

// func (s *PrecheckSuite) Tests() []Test {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	return []Test{
// 		s.TestRamenHubOperatorStatus,
// 		s.TestRamenSpokeOperatorStatus,
// 		s.TestCephClusterStatus,
// 		// check other operator like cert manager, csi, ocm, submariner, rook, velero, minio, volsync
// 	}
// }

// func (s *PrecheckSuite) TestRamenHubOperatorStatus() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	isRunning, podName, err := CheckRamenHubPodRunningStatus(s.Ctx.HubK8sClientSet())
// 	if err != nil {
// 		return err
// 	}

// 	if isRunning {
// 		s.Ctx.Log.Info("Ramen Hub Operator is running", "pod", podName)
// 	} else {
// 		return fmt.Errorf("no running Ramen Hub Operator pod")
// 	}

// 	s.Ctx.Log.Info("TestRamenHubOperatorStatus: Pass")

// 	return nil
// }

// func (s *PrecheckSuite) TestRamenSpokeOperatorStatus() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	isRunning, podName, err := CheckRamenSpokePodRunningStatus(s.Ctx.C1K8sClientSet())
// 	if err != nil {
// 		return err
// 	}

// 	if isRunning {
// 		s.Ctx.Log.Info("Ramen Spoke Operator is running on cluster 1", "pod", podName)
// 	} else {
// 		return fmt.Errorf("no running Ramen Spoke Operator pod on cluster 1")
// 	}

// 	isRunning, podName, err = CheckRamenSpokePodRunningStatus(s.Ctx.C2K8sClientSet())
// 	if err != nil {
// 		return err
// 	}

// 	if isRunning {
// 		s.Ctx.Log.Info("Ramen Spoke Operator is running on cluster 2", "pod", podName)
// 	} else {
// 		return fmt.Errorf("no running Ramen Spoke Operator pod on cluster 2")
// 	}

// 	s.Ctx.Log.Info("TestRamenSpokeOperatorStatus: Pass")

// 	return nil
// }

// // CheckPodRunningStatus checks if there is at least one pod matching the labelSelector
// // in the given namespace that is in the "Running" phase and contains the podIdentifier in its name.
// func CheckPodRunningStatus(client *kubernetes.Clientset, namespace, labelSelector, podIdentifier string) (
// 	bool, string, error,
// ) {
// 	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
// 		LabelSelector: labelSelector,
// 	})
// 	if err != nil {
// 		return false, "", fmt.Errorf("failed to list pods in namespace %s: %v", namespace, err)
// 	}

// 	for _, pod := range pods.Items {
// 		if strings.Contains(pod.Name, podIdentifier) && pod.Status.Phase == "Running" {
// 			return true, pod.Name, nil
// 		}
// 	}

// 	return false, "", nil
// }

// func GetRamenNameSpace(k8sClient *kubernetes.Clientset) (string, error) {
// 	isOpenShift, err := IsOpenShiftCluster(k8sClient)
// 	if err != nil {
// 		return "", err
// 	}

// 	if isOpenShift {
// 		return "openshift-operators", nil
// 	}

// 	return ramenSystemNamespace, nil
// }

// // IsOpenShiftCluster checks if the given Kubernetes cluster is an OpenShift cluster.
// // It returns true if the cluster is OpenShift, false otherwise, along with any error encountered.
// func IsOpenShiftCluster(k8sClient *kubernetes.Clientset) (bool, error) {
// 	discoveryClient := k8sClient.Discovery()

// 	apiGroups, err := discoveryClient.ServerGroups()
// 	if err != nil {
// 		return false, err
// 	}

// 	for _, group := range apiGroups.Groups {
// 		if group.Name == "route.openshift.io" {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

// func CheckRamenHubPodRunningStatus(k8sClient *kubernetes.Clientset) (bool, string, error) {
// 	labelSelector := "app=ramen-hub"
// 	podIdentifier := "ramen-hub-operator"

// 	ramenNameSpace, err := GetRamenNameSpace(k8sClient)
// 	if err != nil {
// 		return false, "", err
// 	}

// 	return CheckPodRunningStatus(k8sClient, ramenNameSpace, labelSelector, podIdentifier)
// }

// func CheckRamenSpokePodRunningStatus(k8sClient *kubernetes.Clientset) (bool, string, error) {
// 	labelSelector := "app=ramen-dr-cluster"
// 	podIdentifier := "ramen-dr-cluster-operator"

// 	ramenNameSpace, err := GetRamenNameSpace(k8sClient)
// 	if err != nil {
// 		return false, "", err
// 	}

// 	return CheckPodRunningStatus(k8sClient, ramenNameSpace, labelSelector, podIdentifier)
// }

// func (s *PrecheckSuite) TestCephClusterStatus() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	c1CtrlClient := s.Ctx.C1CtrlClient()

// 	err := CheckCephClusterRunningStatus(c1CtrlClient)
// 	if err != nil {
// 		return err
// 	}

// 	c2CtrlClient := s.Ctx.C2CtrlClient()

// 	err = CheckCephClusterRunningStatus(c2CtrlClient)
// 	if err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestCephClusterStatus: Pass")

// 	return nil
// }

// func CheckCephClusterRunningStatus(ctrlClient client.Client) error {
// 	rookNamespace := "rook-ceph"

// 	cephclusterList := &rookv1.CephClusterList{}

// 	err := ctrlClient.List(context.Background(), cephclusterList, &client.ListOptions{Namespace: rookNamespace})
// 	if err != nil {
// 		return fmt.Errorf("could not list cephcluster")
// 	}

// 	if len(cephclusterList.Items) == 0 {
// 		return fmt.Errorf("found 0 cephcluster")
// 	}

// 	if len(cephclusterList.Items) > 1 {
// 		return fmt.Errorf("found more than 1 cephcluster")
// 	}

// 	phase := fmt.Sprint(cephclusterList.Items[0].Status.Phase)
// 	if phase != "Ready" {
// 		return fmt.Errorf("cephcluster is not Ready")
// 	}

// 	return nil
// }
