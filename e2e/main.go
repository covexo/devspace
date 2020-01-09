package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/devspace-cloud/devspace/e2e/tests/deploy"
	"github.com/devspace-cloud/devspace/e2e/tests/enter"
	"github.com/devspace-cloud/devspace/e2e/tests/examples"
	"github.com/devspace-cloud/devspace/e2e/tests/initcmd"
	"github.com/devspace-cloud/devspace/e2e/tests/logs"
	"github.com/devspace-cloud/devspace/e2e/tests/space"
	"github.com/devspace-cloud/devspace/e2e/tests/sync"
	"github.com/devspace-cloud/devspace/e2e/utils"
	"github.com/devspace-cloud/devspace/pkg/util/log"
)

var testNamespace = "testing-test-namespace"

// Create a new type for a list of Strings
type stringList []string

// Implement the flag.Value interface
func (s *stringList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

type Test interface {
	Run(subTests []string, ns string, pwd string, logger log.Logger, verbose bool, timeout int) error
	SubTests() []string
}

var availableTests = map[string]Test{
	"examples": examples.RunNew,
	"deploy":   deploy.RunNew,
	"init":     initcmd.RunNew,
	"enter":    enter.RunNew,
	"sync":     sync.RunNew,
	"logs":     logs.RunNew,
	"space":    space.RunNew,
}

var subTests = map[string]*stringList{}

func main() {
	logger := log.GetInstance()
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testCommand := flag.NewFlagSet("test", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	for t := range availableTests {
		subTests[t] = &stringList{}
		testCommand.Var(subTests[t], "test-"+t, "A comma seperated list of sub tests to be passed")
	}

	var test stringList
	testCommand.Var(&test, "test", "A comma seperated list of group tests to pass")

	var verbose bool
	testCommand.BoolVar(&verbose, "verbose", false, "Displays the tests outputs in real time (default: false)")

	var timeout int
	testCommand.IntVar(&timeout, "timeout", 200, "Sets a timeout limit in seconds for each test (default: 200)")

	var testlist stringList
	listCommand.Var(&testlist, "test", "A comma seperated list of group tests to list (leave empty to list all group tests)")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("test or list subcommand is required")
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "test":
		testCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	// If "list" and "test" are used together, only the former will be parsed and recognized, the latter will be ignored
	if listCommand.Parsed() {
		// Required Flags
		fmt.Println("listCommand parsed!")
	}
	if testCommand.Parsed() {
		// We gather all the group tests called with the --test flag. e.g: --test=examples,init
		var testsToRun = map[string]Test{}
		for _, testName := range test {
			if availableTests[testName] == nil {
				// arg is not valid
				fmt.Printf("'%v' is not a valid argument for --test. Valid arguments are the following: [ ", testName)
				for key := range availableTests {
					fmt.Printf("%v ", key)
				}
				fmt.Printf("]\n ")
				os.Exit(1)
			}
			testsToRun[testName] = availableTests[testName]
		}

		if len(testsToRun) == 0 {
			for testName, args := range subTests {
				if args != nil && len(*args) > 0 {
					testsToRun[testName] = availableTests[testName]
				}
			}
		}

		// If cmd test alone (if no --test flag), we want to run all available tests
		if len(testsToRun) == 0 {

			for testName := range availableTests {
				testsToRun[testName] = availableTests[testName]
			}
		}

		for testName, testRun := range testsToRun {
			parameterSubTests := []string{}
			if t, ok := subTests[testName]; ok && t != nil && len(*t) > 0 {
				for _, s := range *t {
					if !utils.StringInSlice(s, testRun.SubTests()) {
						// arg is not valid
						fmt.Printf("'%v' is not a valid argument for --test-%v. Valid arguments are the following: [ ", s, testName)
						for _, st := range testRun.SubTests() {
							fmt.Printf("%v ", st)
						}
						fmt.Printf("]\n ")
						os.Exit(1)
					}

					parameterSubTests = append(parameterSubTests, s)
				}
			}

			// We run the actual group tests by passing the sub tests
			err := testRun.Run(parameterSubTests, testNamespace, pwd, logger, verbose, timeout)
			if err != nil {
				logger.Error(err)
				os.Exit(1)
			}
		}
	}
}

// func runTestWithTimeoutA(testRun Test, parameterSubTests []string, testNamespace string, pwd string, logger log.Logger, verbose bool, timeout int) error {
// 	c1 := make(chan error, 1)

// 	go func() {
// 		err := testRun.Run(parameterSubTests, testNamespace, pwd, logger, verbose, timeout)
// 		c1 <- err
// 	}()

// 	select {
// 	case res := <-c1:
// 		return res
// 	case <-time.After(time.Duration(timeout) * time.Second):
// 		return errors.Errorf("Timeout error: the test did not return within the specified timeout of %v seconds", timeout)
// 	}
// }
