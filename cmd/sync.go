package cmd

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	latest "github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/services"
	"github.com/devspace-cloud/devspace/pkg/devspace/services/targetselector"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/spf13/cobra"
)

// SyncCmd is a struct that defines a command call for "sync"
type SyncCmd struct {
	Selector      string
	Namespace     string
	LabelSelector string
	Container     string
	Pod           string
	Pick          bool

	Exclude       []string
	ContainerPath string
	LocalPath     string
	Verbose       bool
}

// NewSyncCmd creates a new init command
func NewSyncCmd() *cobra.Command {
	cmd := &SyncCmd{}

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Starts a bi-directional sync between the target container and the local path",
		Long: `
#######################################################
################### devspace sync #####################
#######################################################
Starts a bi-directionaly sync between the target container
and the current path:

devspace sync
devspace sync --local-path=subfolder --container-path=/app
devspace sync --exclude=node_modules --exclude=test
devspace sync --pod=my-pod --container=my-container
devspace sync --container-path=/my-path
#######################################################`,
		Run: cmd.Run,
	}

	syncCmd.Flags().StringVarP(&cmd.Selector, "selector", "s", "", "Selector name (in config) to select pod/container for terminal")
	syncCmd.Flags().StringVarP(&cmd.Container, "container", "c", "", "Container name within pod where to execute command")
	syncCmd.Flags().StringVar(&cmd.Pod, "pod", "", "Pod to open a shell to")
	syncCmd.Flags().StringVarP(&cmd.LabelSelector, "label-selector", "l", "", "Comma separated key=value selector list (e.g. release=test)")
	syncCmd.Flags().StringVarP(&cmd.Namespace, "namespace", "n", "", "Namespace where to select pods")
	syncCmd.Flags().BoolVarP(&cmd.Pick, "pick", "p", false, "Select a pod")

	syncCmd.Flags().StringSliceVarP(&cmd.Exclude, "exclude", "e", []string{}, "Exclude directory from sync")
	syncCmd.Flags().StringVar(&cmd.LocalPath, "local-path", ".", "Local path to use (Default is current directory")
	syncCmd.Flags().StringVar(&cmd.ContainerPath, "container-path", "", "Container path to use (Default is working directory)")
	syncCmd.Flags().BoolVar(&cmd.Verbose, "verbose", false, "Shows every file that is synced")

	return syncCmd
}

// Run executes the command logic
func (cmd *SyncCmd) Run(cobraCmd *cobra.Command, args []string) {
	var config *latest.Config
	if configutil.ConfigExists() {
		// Load generated config
		generatedConfig, err := generated.LoadConfig()
		if err != nil {
			log.Fatalf("Error loading generated.yaml: %v", err)
		}

		// Get config with adjusted cluster config
		config, err := configutil.GetContextAjustedConfig(cmd.Namespace, "")
		if err != nil {
			log.Fatal(err)
		}

		// Signal that we are working on the space if there is any
		err = cloud.ResumeSpace(config, generatedConfig, true, log.GetInstance())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Build params
	params := targetselector.CmdParameter{}
	if cmd.Selector != "" {
		params.Selector = &cmd.Selector
	}
	if cmd.Container != "" {
		params.ContainerName = &cmd.Container
	}
	if cmd.LabelSelector != "" {
		params.LabelSelector = &cmd.LabelSelector
	}
	if cmd.Namespace != "" {
		params.Namespace = &cmd.Namespace
	}
	if cmd.Pod != "" {
		params.PodName = &cmd.Pod
	}
	if cmd.Pick != false {
		params.Pick = &cmd.Pick
	}

	// Start terminal
	err := services.StartSyncFromCmd(config, params, cmd.LocalPath, cmd.ContainerPath, cmd.Exclude, cmd.Verbose, log.GetInstance())
	if err != nil {
		log.Fatal(err)
	}
}
