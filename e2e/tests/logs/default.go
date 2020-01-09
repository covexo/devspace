package logs

import (
	"strings"
	"time"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/devspace-cloud/devspace/cmd/flags"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
)

func runDefault(f *customFactory, logger log.Logger) error {
	logger.Info("Run sub test 'default' of test 'logs'")
	logger.StartWait("Run test...")
	defer logger.StopWait()

	lc := &cmd.LogsCmd{
		GlobalFlags: &flags.GlobalFlags{
			Namespace: f.namespace,
		},
		LastAmountOfLines: 20,
	}

	output := "test123"

	done := utils.Capture()

	err := lc.RunLogs(f, nil, nil)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	capturedOutput, err := done()
	if err != nil {
		return err
	}

	if strings.Index(capturedOutput, output) == -1 {
		return errors.Errorf("capturedOutput '%v' is different than output '%s' for the enter cmd", capturedOutput, output)
	}

	return nil
}
