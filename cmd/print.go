package cmd

import (
	"context"
	"fmt"
	gogithub "github.com/google/go-github/v32/github"
	"github.com/nokamoto/print-github-contrib/cmd/contribution"
	"github.com/nokamoto/print-github-contrib/cmd/flags"
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
		start         flags.Time
		end           flags.Time
		output        flags.Output
		repositories  []string
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

			owner := contribution.NewOwner(org, start.Time, end.Time)

			repos, err := client.ListRepositoryByOrg(ctx, org)
			if err != nil {
				return err
			}

			if repositories != nil {
				var whitelist []*gogithub.Repository
				for _, repo := range repos {
					for _, n := range repositories {
						if repo.GetName() == n {
							whitelist = append(whitelist, repo)
						}
					}
				}
				repos = whitelist
			}

			for _, repo := range repos {
				prs, err := client.ListPullRequest(ctx, org, repo.GetName(), start.Time, end.Time)
				if err != nil {
					return err
				}

				contrib, err := client.ListContribution(ctx, org, repo.GetName(), prs)
				if err != nil {
					return err
				}

				owner.AddRepository(contribution.NewRepository(repo.GetName(), contrib))
			}

			switch output.Value {
			case flags.JSON:
				out, err := owner.JSON()
				if err != nil {
					return err
				}

				cmd.Println(out)

			case flags.CSV:
				out, err := owner.CSV()
				if err != nil {
					return err
				}

				cmd.Println(out)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&debug, "debug", false, "turn on debug logging")

	cmd.Flags().StringVar(&enterpriseURL, "enterprise-url", "", "http(s)://hostname/api/v3/")

	cmd.Flags().DurationVar(&sleep, "sleep", 1*time.Second, "api request interval")

	now := flags.Time{Time: time.Now()}

	_ = start.Set(now.String())
	cmd.Flags().Var(&start, "start", fmt.Sprintf("layout=%s", flags.Layout))

	_ = end.Set(now.String())
	cmd.Flags().Var(&end, "end", fmt.Sprintf("layout=%s", flags.Layout))

	output.Value = flags.CSV
	cmd.Flags().VarP(&output, "output", "o", fmt.Sprintf("%s or %s", flags.JSON, flags.CSV))

	cmd.Flags().StringSliceVar(&repositories, "repositories", nil, "whitelist if needed")

	return cmd
}

func init() {
	rootCmd.AddCommand(newPrint())
}
