package workloads

type Workload interface {
	Kustomize() error    // Can differ based on the workload, hence part of the Workload interface
	GetResources() error // Get the actual workload resources

	GetName() string
	GetNameSpace() string
	GetPVCLabel() string
	GetPlacementName() string

	Init()
	Deploy() error
	Undeploy() error
}
