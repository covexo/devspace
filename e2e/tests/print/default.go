package print

import (
	"io/ioutil"
	"strings"

	"github.com/devspace-cloud/devspace/cmd"
	"github.com/devspace-cloud/devspace/cmd/flags"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/pkg/errors"
)

func runDefault(f *customFactory, logger log.Logger) error {
	logger.Info("Run sub test 'default' of test 'print'")
	logger.StartWait("Run test...")
	defer logger.StopWait()

	err := beforeTest(f, "default")
	defer afterTest(f)
	if err != nil {
		return errors.Errorf("test 'print' failed: %s %v", f.GetLogContents(), err)
	}

	pc := &cmd.PrintCmd{
		GlobalFlags: &flags.GlobalFlags{
			Profile: "production",
			Vars:    []string{"MY_IMAGE=testimage"},
		},
		SkipInfo: true,
	}

	done := utils.Capture()

	err = pc.Run(f, nil, nil)
	if err != nil {
		return err
	}

	capturedOutput, err := done()
	if err != nil {
		return err
	}

	_ = utils.ChangeWorkingDir(f.Pwd+"/tests/print", f.GetLog())
	expectedOutput, err := ioutil.ReadFile("./expectedoutput/default")
	if err != nil {
		return err
	}

	if strings.Index(string(expectedOutput), capturedOutput) == -1 {
		return errors.Errorf("output does not match expected output")
	}

	return nil
}
