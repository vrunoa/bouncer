package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/vrunoa/bouncer/internal/bouncer"
	"os"
)

var (
	Version   = "0.0.0+unknown"
	GitCommit = "unknown-commit-sha"
)

func setupLogging(verbose bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"})
}

var configFile string

func checkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check image",
		Long:  "check your image size and ready times",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := bouncer.NewBouncer(configFile)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create bouncer")
			}
			results, err := b.Check()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to run check")
			}
			for _, res := range results {
				if res.Status != 0 {
					log.Fatal().Msg(res.Message)
				} else {
					log.Info().Str("policy", res.Desc).Msg(res.Message)
				}
			}
		},
	}
	cmd.Flags().StringVarP(&configFile, "config-file", "c", "", "config file")
	return cmd
}

func main() {
	setupLogging(true)
	mainCmd := &cobra.Command{
		Use:     "bouncer [command]",
		Short:   "Security guard on docker image sizes",
		Version: fmt.Sprintf("%s\n(build %s)", Version, GitCommit),
	}
	mainCmd.AddCommand(checkCommand())
	if err := mainCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("wops! seems like we meesed up")
		os.Exit(1)
	}
}
