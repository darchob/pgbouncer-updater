package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/configuration"
)

const (
	getApplicationExample = `
	# Write config in current dir with default vars
	%[1]s config 

	# Write config in current dir with flags args
	%[1]s config -dbname postgres -username pgbouncer 
	`

	getUsage = `
	Write config in current dir with default vars or flag args
	`
)

func NewCmdConfig(o *options.Options) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "pgbouncer-updater config",
		Short:        "Write config in current dir",
		Long:         getUsage,
		Aliases:      []string{"config", "cf"},
		Example:      o.Exemple(getApplicationExample),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			conf, err := configuration.NewDefaultConfiguration(o.UserName, o.DBName, o.PGHost, o.Password, o.PGBouncerHosts...).WithPrivKeyFromFile(configuration.DefaultPrivKeyPath)
			if err != nil {
				return err
			}

			if n, err := conf.WriteToFile(); (err != nil) || n <= 0 {
				return fmt.Errorf("config file have %d length and %v error", n, err)
			}

			return nil
		},
	}

	o.WithDefaultFlags(cmd)

	return cmd
}
