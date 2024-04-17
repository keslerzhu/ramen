// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/go-logr/logr"
	"github.com/ramendr/ramen/e2e/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/spf13/viper"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Placement
	ocmclusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	// ManagedClusterSetBinding
	ocmclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
	// PlacementRule
	placementrule "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"
	// Channel
	// channel "open-cluster-management.io/multicloud-operators-channel/pkg/apis"

	ramen "github.com/ramendr/ramen/api/v1alpha1"
	rookv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"

	// Subscription
	argocdv1alpha1hack "github.com/ramendr/ramen/e2e/argocd"
	subscription "open-cluster-management.io/multicloud-operators-subscription/pkg/apis"
)

var (
	kubeconfigHub string
	kubeconfigC1  string
	kubeconfigC2  string
)

func init() {
	flag.StringVar(&kubeconfigHub, "kubeconfig-hub", "kubeconfig/hub/kubeconfig", "Path to the kubeconfig file for the hub cluster")
	flag.StringVar(&kubeconfigC1, "kubeconfig-c1", "kubeconfig/c1/kubeconfig", "Path to the kubeconfig file for the C1 managed cluster")
	flag.StringVar(&kubeconfigC2, "kubeconfig-c2", "kubeconfig/c2/kubeconfig", "Path to the kubeconfig file for the C2 managed cluster")
}

type Cluster struct {
	K8sClientSet *kubernetes.Clientset
	CtrlClient   client.Client
}

// type Context struct {
// 	Log *logr.Logger
// 	Hub Cluster
// 	C1  Cluster
// 	C2  Cluster
// }

var e2eContext *util.Context

func setupClient(kubeconfigPath string) (*kubernetes.Clientset, client.Client, error) {
	var err error

	if kubeconfigPath == "" {
		return nil, nil, fmt.Errorf("kubeconfigPath is empty")
	}

	kubeconfigPath, err = filepath.Abs(kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to determine absolute path to file (%s): %w", kubeconfigPath, err)
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build config from kubeconfig (%s): %w", kubeconfigPath, err)
	}

	k8sClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build k8s client set from kubeconfig (%s): %w", kubeconfigPath, err)
	}

	err = addAllSchemes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add sheme %w", err)
	}

	ctrlClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build controller client from kubeconfig (%s): %w", kubeconfigPath, err)
	}

	return k8sClientSet, ctrlClient, nil
}

func addAllSchemes() error {
	err := ocmclusterv1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	err = ocmclusterv1beta2.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	err = placementrule.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	// err = channel.AddToScheme(scheme.Scheme)
	// if err != nil {
	// 	return err
	// }

	err = subscription.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	err = rookv1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	err = ramen.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	err = argocdv1alpha1hack.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	return err
}

func newContext(log *logr.Logger, hub, c1, c2 string) (*util.Context, error) {
	var err error

	ctx := new(util.Context)
	ctx.Log = log

	ctx.Hub.K8sClientSet, ctx.Hub.CtrlClient, err = setupClient(hub)
	if err != nil {
		return nil, fmt.Errorf("failed to create clients for hub cluster: %w", err)
	}

	ctx.C1.K8sClientSet, ctx.C1.CtrlClient, err = setupClient(c1)
	if err != nil {
		return nil, fmt.Errorf("failed to create clients for c1 cluster: %w", err)
	}

	ctx.C2.K8sClientSet, ctx.C2.CtrlClient, err = setupClient(c2)
	if err != nil {
		return nil, fmt.Errorf("failed to create clients for c2 cluster: %w", err)
	}

	config, err := readConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if config == nil {
		return nil, fmt.Errorf("config is nil.")
	}

	err = configContext(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to config context: %w", err)
	}

	return ctx, nil
}

func validateConfig(config *util.Config) error {
	if config.Clusters["hub"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find hub cluster in configuration")
	}

	if config.Clusters["c1"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find c1 cluster in configuration")
	}

	if config.Clusters["c2"].KubeconfigPath == "" {
		return fmt.Errorf("failed to find c2 cluster in configuration")
	}

	if config.DRPolicy == "" {
		return fmt.Errorf("failed to find drpolicy in configuration")
	}

	if config.Github.Repo == "" {
		return fmt.Errorf("failed to find channel repo in configuration")
	}

	if config.Github.Branch == "" {
		return fmt.Errorf("failed to find channel branch in configuration")
	}

	if config.Timeout < 0 {
		return fmt.Errorf("timeout value is negative")
	}

	if config.Interval < 0 {
		return fmt.Errorf("interval value is negative")
	}

	return nil
}

func readConfig() (*util.Config, error) {
	config := &util.Config{}

	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("failed to find configuration file: %v", err)
		}

		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	if err := viper.UnmarshalExact(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %v", err)
	}

	timeout, success := os.LookupEnv("e2e_timeout")
	if success {
		timeoutInt, err := strconv.Atoi(timeout)
		if err == nil {
			config.Timeout = timeoutInt
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	interval, success := os.LookupEnv("e2e_interval")
	if success {
		intervalInt, err := strconv.Atoi(interval)
		if err == nil {
			config.Interval = intervalInt
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	return config, nil
}

func configContext(ctx *util.Context, config *util.Config) error {
	ctx.Config = config
	// ctx.Clusters = make(util.Clusters)

	// for clusterName, cluster := range config.Clusters {
	// 	k8sClientSet, ctrlClient, err := util.GetClientSetFromKubeConfigPath(cluster.KubeconfigPath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	ctx.Clusters[clusterName] = &util.Cluster{
	// 		K8sClientSet: k8sClientSet,
	// 		CtrlClient:   ctrlClient,
	// 	}
	// }

	return nil
}

func TestMain(m *testing.M) {
	var err error

	flag.Parse()

	log := zap.New(zap.UseFlagOptions(&zap.Options{
		Development: true,
		ZapOpts: []uberzap.Option{
			uberzap.AddCaller(),
		},
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}))

	e2eContext, err = newContext(&log, kubeconfigHub, kubeconfigC1, kubeconfigC2)
	if err != nil {
		log.Error(err, "unable to create new testing context")

		panic(err)
	}

	os.Exit(m.Run())
}

type testDef struct {
	name string
	test func(t *testing.T)
}

var Suites = []testDef{
	{"Basic", Basic},
	// {"Validate", Validate},
	// {"Exhaustive", Exhaustive},
}

func TestSuites(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	for idx := range Suites {
		t.Run(Suites[idx].name, Suites[idx].test)
	}
}
