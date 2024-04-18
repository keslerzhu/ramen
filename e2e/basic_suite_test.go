// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"testing"

	"github.com/ramendr/ramen/e2e/deployers"
	"github.com/ramendr/ramen/e2e/dractions"
	"github.com/ramendr/ramen/e2e/util"
	"github.com/ramendr/ramen/e2e/workloads"
)

type BasicSuite struct {
	w workloads.Workload
	d deployers.Deployer
	r dractions.DRActions
}

func newBasicSuite(e2eContext *util.Context) (*BasicSuite, error) {
	deployment := &workloads.Deployment{Ctx: e2eContext}
	deployment.Init()

	subscription := &deployers.Subscription{Ctx: e2eContext}
	subscription.Init()

	dractions := dractions.DRActions{Ctx: e2eContext}

	bs := BasicSuite{w: deployment, d: subscription, r: dractions}

	return &bs, nil
}

var basicSuite *BasicSuite

func Basic(t *testing.T) {
	t.Helper()

	e2eContext.Log.Info(t.Name())

	var err error
	basicSuite, err = newBasicSuite(e2eContext)

	if basicSuite == nil {
		t.Error("basicSuite is nil")
	}

	if err != nil {
		t.Error(err)
	}

	// fmt.Println(basicSuite)
	// fmt.Println(basicSuite.w)
	// fmt.Println(basicSuite.d)
	// fmt.Println(basicSuite.r)

	t.Run("Deploy", Deploy)
	t.Run("Enable", Enable)
	t.Run("Failover", Failover)
	t.Run("Relocate", Relocate)
	t.Run("Disable", Disable)
	t.Run("Undeploy", Undeploy)
}

func Deploy(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.d.Deploy(basicSuite.w); err != nil {
		t.Error(err)
	}
}

func Enable(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.r.EnableProtection(basicSuite.w, basicSuite.d); err != nil {
		t.Error(err)
	}
}

func Failover(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.r.Failover(basicSuite.w, basicSuite.d); err != nil {
		t.Error(err)
	}
}

func Relocate(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.r.Relocate(basicSuite.w, basicSuite.d); err != nil {
		t.Error(err)
	}
}

func Disable(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.r.DisableProtection(basicSuite.w, basicSuite.d); err != nil {
		t.Error(err)
	}
}

func Undeploy(t *testing.T) {
	e2eContext.Log.Info(t.Name())

	if err := basicSuite.d.Undeploy(basicSuite.w); err != nil {
		t.Error(err)
	}
}
