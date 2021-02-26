package services

import (
	"context"
	"io"
	"os"

	"github.com/loft-sh/devspace/pkg/devspace/services/targetselector"
	"github.com/mgutz/ansi"
)

// StartLogs print the logs and then attaches to the container
func (serviceClient *client) StartLogs(options targetselector.Options, follow bool, tail int64) error {
	return serviceClient.StartLogsWithWriter(options, follow, tail, os.Stdout)
}

// StartLogsWithWriter prints the logs and then attaches to the container with the given stdout and stderr
func (serviceClient *client) StartLogsWithWriter(options targetselector.Options, follow bool, tail int64, writer io.Writer) error {
	targetSelector := targetselector.NewTargetSelector(serviceClient.client)
	options.FilterContainer = nil

	container, err := targetSelector.SelectSingleContainer(context.TODO(), options, serviceClient.log)
	if err != nil {
		return err
	}

	serviceClient.log.Infof("Printing logs of pod:container %s:%s", ansi.Color(container.Pod.Name, "white+b"), ansi.Color(container.Container.Name, "white+b"))
	reader, err := serviceClient.client.Logs(context.Background(), container.Pod.Namespace, container.Pod.Name, container.Container.Name, false, &tail, follow)
	if err != nil {
		return nil
	}

	_, err = io.Copy(writer, reader)
	return err
}
