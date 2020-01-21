package use

import (
	"sort"

	"github.com/devspace-cloud/devspace/pkg/util/factory"
	"github.com/devspace-cloud/devspace/pkg/util/survey"

	"github.com/mgutz/ansi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ContextCmd struct{}

func newContextCmd(f factory.Factory) *cobra.Command {
	cmd := &ContextCmd{}

	useContext := &cobra.Command{
		Use:   "context",
		Short: "Tells DevSpace which kube context to use",
		Long: `
#######################################################
############### devspace use context ##################
#######################################################
Switch the current kube context

Example:
devspace use context my-context
#######################################################
	`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.RunUseContext(f, cobraCmd, args)
		},
	}

	return useContext
}

// RunUseContext executes the functionality "devspace use namespace"
func (cmd *ContextCmd) RunUseContext(f factory.Factory, cobraCmd *cobra.Command, args []string) error {
	// Load kube-config
	log := f.GetLog()
	kubeLoader := f.NewKubeConfigLoader()
	kubeConfig, err := kubeLoader.LoadRawConfig()
	if err != nil {
		return errors.Wrap(err, "load kube config")
	}

	var context string
	if len(args) > 0 {
		// First arg is context name
		context = args[0]
	} else {
		contexts := []string{}
		for ctx := range kubeConfig.Contexts {
			contexts = append(contexts, ctx)
		}

		sort.Strings(contexts)

		context, err = log.Question(&survey.QuestionOptions{
			Question:     "Which context do you want to use?",
			DefaultValue: kubeConfig.CurrentContext,
			Options:      contexts,
		})
		if err != nil {
			return err
		}
	}

	// Save old context
	oldContext := kubeConfig.CurrentContext

	// Set current kube-context
	kubeConfig.CurrentContext = context

	if oldContext != context {
		// Save updated kube-config
		kubeLoader.SaveConfig(kubeConfig)

		log.Infof("Your kube-context has been updated to '%s'", ansi.Color(kubeConfig.CurrentContext, "white+b"))
		log.Infof("\r         To revert this operation, run: %s\n", ansi.Color("devspace use context "+oldContext, "white+b"))
	}

	log.Donef("Successfully set kube-context to '%s'", ansi.Color(context, "white+b"))
	return nil
}
