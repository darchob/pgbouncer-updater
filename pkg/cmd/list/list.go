package list

import (
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/configuration"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/databases"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/userlist"
)

const (
	getApplicationExample = `
	# Get all role from db write they to default file
	%[1]s list

	# Get all role from db write they to a file
	%[1]s list -file /tmp/userlist.txt 

	# Get all role from db with config
	%[1]s list -config /etc/pgbouncer-updater/config.yaml
	`

	getUsage = `
	Exec query inside postgres database to find role name and role password, write they to a file.
	`
)

func NewCmdUpdateUserList(o *options.Options) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "pgbouncer-updater list",
		Short:        "Get user list from DB",
		Long:         getUsage,
		Aliases:      []string{"ll", "list"},
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

			if err := ListCmd(c, o, conf); err != nil {
				return err
			}

			return nil
		},
	}

	o.WithDefaultFlags(cmd)
	cmd.Flags().StringVar(&o.Query, "query", o.WithDefaultOptions().Query, "Query to get Roles from DB")
	cmd.Flags().StringVar(&o.ConfigFilePath, "config", o.WithDefaultOptions().ConfigFilePath, "Config file path")
	cmd.Flags().StringVar(&o.File, "file", o.WithDefaultOptions().File, "USer list file")
	return cmd
}

func ListCmd(c *cobra.Command, o *options.Options, conf configuration.Configurations) error {

	// Parse config to DSN Format
	dsn, err := conf.GetPostgresDSN()
	if err != nil {
		return err
	}

	// Configure Postgres connection
	db, err := databases.NewQuery(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Exec query to map
	data, err := db.ToMap(o.Query)
	if err != nil {
		return err
	}

	// Configure file
	list, err := userlist.NewUserListToFile(o.File)
	if err != nil {
		return err
	}

	// Write data to a file
	err = list.WriteMany(data)
	if err != nil {
		return err
	}
	return nil
}
