package cmd

import (
	"context"
	"github.com/nokamoto/print-github-contrib/cmd/github"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

func newPrint() *cobra.Command {
	var (
		debug         bool
		enterpriseURL string
		sleep         time.Duration
	)

	cmd := &cobra.Command{
		Use:           "print ORGANIZATION",
		Short:         "Print contribution to organization's repositories",
		Long:          `Print contribution to organization's repositories.`,
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := zap.NewProductionConfig()
			if debug {
				cfg.Level.SetLevel(zap.DebugLevel)
			}

			logger, err := cfg.Build()
			if err != nil {
				panic(err)
			}
			defer logger.Sync()

			org := args[0]

			client := github.NewClient(logger, sleep)
			if len(enterpriseURL) != 0 {
				client, err = github.NewEnterpriseClient(logger, enterpriseURL, sleep)
				if err != nil {
					return err
				}
			}

			ctx := context.Background()

			_, err = client.ListRepositoryByOrg(ctx, org)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&debug, "debug", false, "turn on debug logging")

	cmd.Flags().StringVar(&enterpriseURL, "enterprise-url", "", "http(s)://hostname/api/v3/")

	cmd.Flags().DurationVar(&sleep, "sleep", 1*time.Second, "api request interval")

	return cmd
}

func init() {
	rootCmd.AddCommand(newPrint())
}
