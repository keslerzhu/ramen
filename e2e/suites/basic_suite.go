package suites

// import (
// 	"github.com/ramendr/ramen/e2e/deployers"
// 	"github.com/ramendr/ramen/e2e/dractions"
// 	"github.com/ramendr/ramen/e2e/util"
// 	"github.com/ramendr/ramen/e2e/workloads"
// )

// type BasicSuite struct {
// 	w   workloads.Workload
// 	d   deployers.Deployer
// 	r   dractions.DRActions
// 	Ctx *util.Context
// }

// func (s *BasicSuite) SetContext(ctx *util.Context) {
// 	util.LogEnter(&ctx.Log)
// 	defer util.LogExit(&ctx.Log)

// 	s.Ctx = ctx
// }

// func (s *BasicSuite) SetupSuite() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	deployment := &workloads.Deployment{Ctx: s.Ctx}
// 	deployment.Init()
// 	s.w = deployment

// 	sub := &deployers.Subscription{Ctx: s.Ctx}
// 	sub.Init()
// 	s.d = sub

// 	s.r = dractions.DRActions{Ctx: s.Ctx}

// 	return nil
// }

// func (s *BasicSuite) TeardownSuite() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	return nil
// }

// func (s *BasicSuite) Tests() []Test {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	return []Test{
// 		s.TestWorkloadDeployment,
// 		s.TestEnableProtection,
// 		s.TestWorkloadFailover,
// 		s.TestWorkloadRelocation,
// 		s.TestDisableProtection,
// 		s.TestWorkloadUndeployment,
// 	}
// }

// func (s *BasicSuite) TestWorkloadDeployment() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	if err := s.d.Deploy(s.w); err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestWorkloadDeployment: Pass")

// 	return nil
// }

// func (s *BasicSuite) TestEnableProtection() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	err := s.r.EnableProtection(s.w, s.d)
// 	if err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestEnableProtection: Pass")

// 	return nil
// }

// func (s *BasicSuite) TestWorkloadFailover() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	if err := s.r.Failover(s.w, s.d); err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestWorkloadFailover: Pass")

// 	return nil
// }

// func (s *BasicSuite) TestWorkloadRelocation() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	if err := s.r.Relocate(s.w, s.d); err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestWorkloadRelocation: Pass")

// 	return nil
// }

// func (s *BasicSuite) TestDisableProtection() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	err := s.r.DisableProtection(s.w, s.d)
// 	if err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestDisableProtection: Pass")

// 	return nil
// }

// func (s *BasicSuite) TestWorkloadUndeployment() error {
// 	util.LogEnter(&s.Ctx.Log)
// 	defer util.LogExit(&s.Ctx.Log)

// 	if err := s.d.Undeploy(s.w); err != nil {
// 		return err
// 	}

// 	s.Ctx.Log.Info("TestWorkloadUndeployment: Pass")

// 	return nil
// }
