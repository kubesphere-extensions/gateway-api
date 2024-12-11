package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/gops/agent"
	"github.com/kubesphere-extensions/gateway-api/cmd/app/options"
	apiserverconfig "github.com/kubesphere-extensions/gateway-api/pkg/config"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	_ "sigs.k8s.io/gateway-api/apis"
)

func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()

	conf, err := apiserverconfig.TryLoadFromDisk()
	if err == nil {
		s.Config = conf
	} else {
		klog.Fatalf("Failed to load configuration from disk: %v", err)
	}

	cmd := &cobra.Command{
		Use: "Gateway API Server",
		Long: `The KubeSphere Gateway API server validates and configures data for the API objects. 
The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if s.GOPSEnabled {
				// Add agent to report additional information such as the current stack trace, Go version, memory stats, etc.
				// Bind to a random port on address 127.0.0.1.
				if err := agent.Listen(agent.Options{}); err != nil {
					klog.Fatal(err)
				}
			}

			return Run(s, signals.SetupSignalHandler())
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	// cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cols := 10
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of GatewayAPI apiserver",
		// TODO
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func Run(s *options.ServerRunOptions, ctx context.Context) error {
	apiserver, err := s.NewAPIServer()
	if err != nil {
		return err
	}

	err = apiserver.PrepareRun()
	if err != nil {
		return err
	}

	err = apiserver.Run(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
