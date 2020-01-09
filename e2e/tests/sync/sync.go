package sync

import (
	"bytes"
	"time"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/devspace-cloud/devspace/cmd/flags"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/util/factory"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type customFactory struct {
	*factory.DefaultFactoryImpl
	verbose     bool
	timeout     int
	namespace   string
	pwd         string
	cacheLogger log.Logger
	dirPath     string
	client      kubectl.Client
}

// GetLog implements interface
func (c *customFactory) GetLog() log.Logger {
	return c.cacheLogger
}

type Runner struct{}

var RunNew = &Runner{}

func (r *Runner) SubTests() []string {
	subTests := []string{}
	for k := range availableSubTests {
		subTests = append(subTests, k)
	}

	return subTests
}

var availableSubTests = map[string]func(factory *customFactory, logger log.Logger) error{
	"default": runDefault,
}

func (r *Runner) Run(subTests []string, ns string, pwd string, logger log.Logger, verbose bool, timeout int) error {
	buff := &bytes.Buffer{}
	var cacheLogger log.Logger
	cacheLogger = log.NewStreamLogger(buff, logrus.InfoLevel)

	var buffString string
	buffString = buff.String()

	if verbose {
		cacheLogger = logger
		buffString = ""
	}

	logger.Info("Run 'sync' test")

	// Populates the tests to run with all the available sub tests if no sub tests are specified
	if len(subTests) == 0 {
		for subTestName := range availableSubTests {
			subTests = append(subTests, subTestName)
		}
	}

	f := &customFactory{
		namespace:   ns,
		pwd:         pwd,
		cacheLogger: cacheLogger,
		verbose:     verbose,
		timeout:     timeout,
	}

	// Runs the tests
	for _, subTestName := range subTests {
		c1 := make(chan error)

		go func() {
			err := func() error {
				f.namespace = utils.GenerateNamespaceName("test-sync-" + subTestName)

				err := beforeTest(f)
				defer afterTest(f)
				if err != nil {
					return errors.Errorf("test 'sync' failed: %s %v", buffString, err)
				}

				err = availableSubTests[subTestName](f, logger)
				utils.PrintTestResult("sync", subTestName, err, logger)
				if err != nil {
					return errors.Errorf("test 'sync' failed: %s %v", buffString, err)
				}

				return nil
			}()
			c1 <- err
		}()

		select {
		case err := <-c1:
			if err != nil {
				return err
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			return errors.Errorf("Timeout error: the test did not return within the specified timeout of %v seconds", timeout)
		}
	}

	return nil
}

func beforeTest(f *customFactory) error {
	deployConfig := &cmd.DeployCmd{
		GlobalFlags: &flags.GlobalFlags{
			Namespace: f.namespace,
			NoWarn:    true,
		},
		ForceBuild:  true,
		ForceDeploy: true,
		SkipPush:    true,
	}

	dirPath, _, err := utils.CreateTempDir()
	if err != nil {
		return err
	}

	f.dirPath = dirPath

	err = utils.Copy(f.pwd+"/tests/sync/testdata", dirPath)
	if err != nil {
		return err
	}

	err = utils.ChangeWorkingDir(dirPath+"/quickstart", f.cacheLogger)
	if err != nil {
		return err
	}

	// Create kubectl client
	client, err := f.NewKubeClientFromContext(deployConfig.KubeContext, deployConfig.Namespace, deployConfig.SwitchContext)
	if err != nil {
		return errors.Errorf("Unable to create new kubectl client: %v", err)
	}

	f.client = client

	err = deployConfig.Run(f, nil, nil)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	// Checking if pods are running correctly
	err = utils.AnalyzePods(client, f.namespace, f.cacheLogger)
	if err != nil {
		return err
	}

	return nil
}

func afterTest(f *customFactory) {
	utils.DeleteTempAndResetWorkingDir(f.dirPath, f.pwd, f.cacheLogger)
	utils.DeleteNamespace(f.client, f.namespace)
}
