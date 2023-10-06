package reload

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/configuration"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/databases"
)

const (
	getApplicationExample = `
	# Reload PGBouncer hosts
	%[1]s reload

	#  Reload PGBouncer hosts with config
	%[1]s reload -config /etc/pgbouncer-updater/config.yaml

	#  Reload PGBouncer with no default query
	%[1]s reload -query "reload"
	`

	getUsage = `
	Exec reload query in PGBouncer DB
	`
)

const (
	defaultPGBouncerDB  = "postgres"
	defaultPostgresPort = 5432
)

func NewCmdReload(o *options.Options) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "pgbouncer-updater reload",
		Short:        "Reload PGBouncer user list ",
		Long:         o.Exemple(getApplicationExample),
		Aliases:      []string{"reload", "r"},
		Example:      getApplicationExample,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			// Get config from file
			log.Info("Load configurations from file ", o.ConfigFilePath)
			conf, err := configuration.NewConfigurationFromFile(o.ConfigFilePath)
			if (err != nil) && err != configuration.FileNotFound {
				return err
			}

			if err == configuration.FileNotFound {
				log.Info("Failed to load configurations from file ", o.ConfigFilePath, " use default config with args")
				conf = configuration.NewDefaultConfiguration(o.UserName, o.DBName, o.PGHost, o.Password, o.PGBouncerHosts...)
			}

			if err := ReloadCmd(c, o, conf); err != nil {
				return err
			}

			return nil
		},
	}

	o.WithDefaultFlags(cmd)
	cmd.Flags().StringVar(&o.ConfigFilePath, "config", "/etc/pgboncer-updater/config.yaml", "Config file path")
	cmd.Flags().StringVar(&o.Query, "query", "reload", "Reload pgbouncer database")
	return cmd
}

func ReloadCmd(c *cobra.Command, o *options.Options, conf configuration.Configurations) error {
	hostVars, err := conf.GetPGBouncerHost()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error, len(hostVars))
	log.Info("Start PGBouncer reload PGBouncer hosts")
	for _, hostvar := range hostVars {
		host := hostvar.DeepCopy()

		wg.Add(1)
		go func(pgHost *configuration.PGBouncerHost, conf configuration.Configurations) {
			defer wg.Done()

			dsn, err := conf.GetPostgresCustomDSN(
				defaultPGBouncerDB,
				host.Host,
				"disable",
				defaultPostgresPort)

			if err != nil {
				errCh <- err
				return
			}

			// Configure Postgres connection
			log.Info("Launch PGBouncer reload query on host ", host.Host)
			db, err := databases.NewQuery(dsn)
			if err != nil {
				log.Error("Failed to exec query with dsn ", dsn)
				errCh <- err
				return
			}
			defer db.Close()

			if err := db.ToVoid(o.Query); err != nil {
				errCh <- err
				return
			}
			log.Info("PGBouncer reload query on host ", host.Host, " done")
		}(host, conf)
	}
	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
