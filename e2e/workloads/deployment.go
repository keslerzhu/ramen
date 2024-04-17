package workloads

import (
	"github.com/ramendr/ramen/e2e/util"
)

type Deployment struct {
	//	RepoURL  string // Possibly all this is part of Workload than each implementation of the interfaces?
	Path string
	//	Revision string
	AppName string

	Ctx *util.Context
}

func (w *Deployment) Init() {
	//	w.RepoURL = "https://github.com/ramendr/ocm-ramen-samples.git"
	w.Path = "workloads/deployment/k8s-regional-rbd"
	// w.Revision = "main"
	w.AppName = "busybox"
}

func (w Deployment) GetAppName() string {
	return w.AppName
}

// func (w Deployment) GetRepoURL() string {
// 	return w.RepoURL
// }

func (w Deployment) GetPath() string {
	return w.Path
}

// func (w Deployment) GetRevision() string {
// 	return w.Revision
// }

// func (w Deployment) GetResourceURL() string {
// 	//by default the timeout is 27s, could fail sometimes
// 	return w.RepoURL + "/" + w.Path + "?ref=" + w.Revision + "&timeout=90s"
// }

func (w Deployment) Kustomize() error {
	util.LogEnter(w.Ctx.Log)
	defer util.LogExit(w.Ctx.Log)

	return nil
}

func (w Deployment) GetResources() error {
	// this would be a common function given the vars? But we need the resources Kustomized
	util.LogEnter(w.Ctx.Log)
	defer util.LogExit(w.Ctx.Log)

	return nil
}

func (w Deployment) Health() error {
	// Check the workload health on a targetCluster
	util.LogEnter(w.Ctx.Log)
	defer util.LogExit(w.Ctx.Log)

	return nil
}
