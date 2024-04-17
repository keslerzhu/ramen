package e2e

// import (
// 	"fmt"
// 	"os"
// 	"strconv"

// 	"github.com/ramendr/ramen/e2e/util"
// 	"github.com/spf13/viper"
// )

// func validateConfig(config *util.Config) error {
// 	if config.Clusters["hub"].KubeconfigPath == "" {
// 		return fmt.Errorf("failed to find hub cluster in configuration")
// 	}

// 	if config.Clusters["c1"].KubeconfigPath == "" {
// 		return fmt.Errorf("failed to find c1 cluster in configuration")
// 	}

// 	if config.Clusters["c2"].KubeconfigPath == "" {
// 		return fmt.Errorf("failed to find c2 cluster in configuration")
// 	}

// 	if config.DRPolicy == "" {
// 		return fmt.Errorf("failed to find drpolicy in configuration")
// 	}

// 	if config.Github.Repo == "" {
// 		return fmt.Errorf("failed to find channel repo in configuration")
// 	}

// 	if config.Github.Branch == "" {
// 		return fmt.Errorf("failed to find channel branch in configuration")
// 	}

// 	if config.Timeout < 0 {
// 		return fmt.Errorf("timeout value is negative")
// 	}

// 	if config.Interval < 0 {
// 		return fmt.Errorf("interval value is negative")
// 	}

// 	return nil
// }

// func readConfig() (*util.Config, error) {
// 	config := &util.Config{}

// 	viper.SetConfigFile("config.yaml")

// 	if err := viper.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			return nil, fmt.Errorf("failed to find configuration file: %v", err)
// 		}

// 		return nil, fmt.Errorf("failed to read configuration file: %v", err)
// 	}

// 	if err := viper.UnmarshalExact(config); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
// 	}

// 	if err := validateConfig(config); err != nil {
// 		return nil, fmt.Errorf("failed to validate configuration: %v", err)
// 	}

// 	timeout, success := os.LookupEnv("e2e_timeout")
// 	if success {
// 		timeoutInt, err := strconv.Atoi(timeout)
// 		if err == nil {
// 			config.Timeout = timeoutInt
// 		} else {
// 			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
// 			os.Exit(1)
// 		}
// 	}

// 	interval, success := os.LookupEnv("e2e_interval")
// 	if success {
// 		intervalInt, err := strconv.Atoi(interval)
// 		if err == nil {
// 			config.Interval = intervalInt
// 		} else {
// 			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
// 			os.Exit(1)
// 		}
// 	}

// 	return config, nil
// }

// func configContext(ctx *util.Context, config *util.Config) error {
// 	ctx.Config = config
// 	// ctx.Clusters = make(util.Clusters)

// 	// for clusterName, cluster := range config.Clusters {
// 	// 	k8sClientSet, ctrlClient, err := util.GetClientSetFromKubeConfigPath(cluster.KubeconfigPath)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}

// 	// 	ctx.Clusters[clusterName] = &util.Cluster{
// 	// 		K8sClientSet: k8sClientSet,
// 	// 		CtrlClient:   ctrlClient,
// 	// 	}
// 	// }

// 	return nil
// }
