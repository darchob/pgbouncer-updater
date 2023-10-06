package aio

import (
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/copy"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/list"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/reload"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/configuration"
)

const (
	getApplicationExample = `
	# Write config in current dir with default vars
	%[1]s aio --config config.yaml 
	`

	getUsage = `
	Write config in current dir with default vars or flag args
	`
)

func NewCmdAIO(o *options.Options) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "pgbouncer-updater aio",
		Short:        "All in one command",
		Long:         getUsage,
		Aliases:      []string{"aio", "all"},
		Example:      o.Exemple(getApplicationExample),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			// Get config from file
			conf, err := configuration.NewConfigurationFromFile(o.ConfigFilePath)
			if (err != nil) && err != configuration.FileNotFound {
				return err
			}

			if err == configuration.FileNotFound {
				conf = configuration.NewDefaultConfiguration(o.UserName, o.DBName, o.PGHost, o.Password, o.PGBouncerHosts...)
			}

			if err := list.ListCmd(c, o, conf); err != nil {
				return err
			}

			if err := copy.CopyCmd(c, o, conf); err != nil {
				return err
			}

			if err := reload.ReloadCmd(c, o, conf); err != nil {
				return err
			}

			return nil
		},
	}

	o.WithDefaultFlags(cmd)
	cmd.Flags().BoolVar(&o.Sudo, "sudo", o.Sudo, "Copy file as sudoer")
	cmd.Flags().StringVar(&o.ConfigFilePath, "config", o.WithDefaultOptions().ConfigFilePath, "Config file path")
	return cmd
}
