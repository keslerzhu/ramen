package main

import (
	"fmt"
	"sync"

	"github.com/ramendr/ramen/e2e/suites"
	"github.com/ramendr/ramen/e2e/util"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func configureLogOptions() *zap.Options {
	opts := zap.Options{
		Development: true,
		ZapOpts: []uberzap.Option{
			uberzap.AddCaller(),
		},
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}

	return &opts
}

func main() {

	logOpts := configureLogOptions()
	log := zap.New(zap.UseFlagOptions(logOpts))

	ctx := new(util.TestContext)
	ctx.Log = log

	ctx.Log.Info("enter main()")

	config, err := readConfig()
	if err != nil {
		ctx.Log.Error(err, "failed to read configuration")
		panic(err)
	}

	if config == nil {
		ctx.Log.Error(fmt.Errorf("config is nill"), "config is nill")
		panic(config)
	}
	err = configContext(ctx, config)
	if err != nil {
		ctx.Log.Error(err, "failed to config TestContext")
		panic(err)
	}

	// All TestSuites can run simultaneously.

	var wg sync.WaitGroup
	var numTestSuites = 4
	wg.Add(numTestSuites)

	go func() {
		defer wg.Done()
		err = RunSuite(&suites.PrecheckSuite{}, ctx)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err = RunSuite(&suites.BasicSuite{}, ctx)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err = RunSuite(&suites.AppsetSuite{}, ctx)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err = RunSuite(&suites.ImperativeSuite{}, ctx)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	ctx.Log.Info("exit main()")
}

func RunSuite(suite suites.TestSuite, ctx *util.TestContext) error {
	suite.SetContext(ctx)

	if err := suite.SetupSuite(); err != nil {
		return fmt.Errorf("setup suite failed: %w", err)
	}

	defer func() {
		if err := suite.TeardownSuite(); err != nil {
			panic(fmt.Errorf("teardown suite failed: %w", err))
		}
	}()

	for _, test := range suite.Tests() {
		if err := test(); err != nil {
			ctx.Log.Error(err, "test failed")
			return fmt.Errorf("test failed: %w", err)
		}
	}

	return nil
}
