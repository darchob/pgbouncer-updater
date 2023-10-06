package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
)

func main() {

	opts := options.NewPGBouncerUpdaterOptions()
	rootCmd := cmd.NewCmdPGBouncerUpdate(opts.WithDefaultOptions())
	if err := rootCmd.Execute(); err != nil {
		log.Error("Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
