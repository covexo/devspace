package services

import (
	"os"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/kubectl"
	"github.com/devspace-cloud/devspace/pkg/devspace/services/targetselector"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"

	"github.com/mgutz/ansi"
	kubectlExec "k8s.io/client-go/util/exec"
)

// StartTerminal opens a new terminal
func StartTerminal(config *latest.Config, generatedConfig *generated.Config, client kubectl.Client, selectorParameter *targetselector.SelectorParameter, args []string, imageSelector []string, interrupt chan error, wait bool, log log.Logger) (int, error) {
	command := getCommand(config, args)

	targetSelector, err := targetselector.NewTargetSelector(config, client, selectorParameter, true, imageSelector)
	if err != nil {
		return 0, err
	}

	// Allow picking non running pods
	targetSelector.AllowNonRunning = true
	targetSelector.SkipWait = !wait
	targetSelector.PodQuestion = ptr.String("Which pod do you want to open the terminal for?")

	pod, container, err := targetSelector.GetContainer(true, log)
	if err != nil {
		return 0, err
	}

	wrapper, upgradeRoundTripper, err := client.GetUpgraderWrapper()
	if err != nil {
		return 0, err
	}

	log.Infof("Opening shell to pod:container %s:%s", ansi.Color(pod.Name, "white+b"), ansi.Color(container.Name, "white+b"))
	if selectorParameter.CmdParameter.Interactive == true && len(container.Command) > 0 && config != nil && generatedConfig != nil && config.Dev != nil && config.Dev.Interactive != nil && len(config.Dev.Interactive.Images) > 0 {
		for _, image := range config.Dev.Interactive.Images {
			imageSelector = []string{}
			imageConfigCache := generatedConfig.GetActive().GetImageCache(image.Name)
			if imageConfigCache != nil && imageConfigCache.ImageName != "" {
				imageName := imageConfigCache.ImageName + ":" + imageConfigCache.Tag
				if imageName == container.Image && (len(image.Entrypoint) > 0 || len(image.Cmd) > 0) {
					log.Warnf("The container you are entering was started with a Kubernetes `command` option (%s) instead of the original Dockerfile ENTRYPOINT. Interactive mode ENTRYPOINT override does not work for containers started using with Kubernetes command.", container.Command)
				}
			}
		}
	}

	go func() {
		interrupt <- client.ExecStreamWithTransport(wrapper, upgradeRoundTripper, pod, container.Name, command, true, os.Stdin, os.Stdout, os.Stderr, kubectl.SubResourceExec)
	}()

	err = <-interrupt
	upgradeRoundTripper.Close()
	if err != nil {
		if exitError, ok := err.(kubectlExec.CodeExitError); ok {
			return exitError.Code, nil
		}

		return 0, err
	}

	return 0, nil
}

func getCommand(config *latest.Config, args []string) []string {
	var command []string

	if config != nil && config.Dev != nil && config.Dev.Interactive != nil && config.Dev.Interactive.Terminal != nil && len(config.Dev.Interactive.Terminal.Command) > 0 {
		for _, cmd := range config.Dev.Interactive.Terminal.Command {
			command = append(command, cmd)
		}
	}

	if len(args) > 0 {
		command = args
	} else {
		if len(command) == 0 {
			command = []string{
				"sh",
				"-c",
				"command -v bash >/dev/null 2>&1 && exec bash || exec sh",
			}
		}
	}

	return command
}
