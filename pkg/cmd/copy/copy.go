package copy

import (
	"fmt"
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/configuration"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/sendfile"
)

const (
	getApplicationExample = `
	# Copy user list from default to server
	%[1]s copy

	# Copy user list from default to server with file
	%[1]s copy -file /tmp/userlist.txt 

	# Copy user list from default to server with config
	%[1]s copy -config /etc/pgbouncer-updater/config.yaml

	# Copy user list from default to server as suoder
	%[1]s copy -sudo
	`

	getUsage = `
	Copy User list from file to hosts
	`

	DefaultUserlistRemotePath = "/etc/pgbouncer/userlist.txt"
	DefaultUserlistOldPath    = "userlist.txt.old"
)

func NewCmdCopyUserList(o *options.Options) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "pgbouncer-updater copy",
		Short:        "Copy user list from file to hosts",
		Long:         getUsage,
		Aliases:      []string{"copy", "cp"},
		Example:      getApplicationExample,
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

			if err := CopyCmd(c, o, conf); err != nil {
				return err
			}

			return nil
		},
	}

	o.WithDefaultFlags(cmd)
	cmd.Flags().BoolVar(&o.Sudo, "sudo", o.Sudo, "Copy file as sudoer")
	cmd.Flags().StringVar(&o.ConfigFilePath, "config", "/etc/pgboncer-updater/config.yaml", "Config file path")
	cmd.Flags().StringVar(&o.File, "file", "userlist.txt", "User list file")
	cmd.Flags().StringVar(&o.DestinationFile, "remote", "/etc/pgbouncer/userlist.txt", "Remote users list file path")
	return cmd
}

func CopyCmd(c *cobra.Command, o *options.Options, conf configuration.Configurations) error {
	hostVars, err := conf.GetPGBouncerHost()
	if err != nil {
		return err
	}

	if len(hostVars) <= 0 {
		return fmt.Errorf("no hosts to copy userlists")
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error, len(hostVars))
	log.Info("Start copying userlist to hosts")
	for _, hostvar := range hostVars {
		host := hostvar.DeepCopy()
		log.Info("Try to connect to pgbouncer host ", host.Host)

		wg.Add(1)
		go func(pgHost *configuration.PGBouncerHost) {
			log.Info("Connect to ", pgHost.Host)
			host := fmt.Sprintf("%s:%d", pgHost.Host, pgHost.Port)
			scp := sendfile.NewScpClient(host, pgHost.UserName, pgHost.Port, []byte(pgHost.PrivKey), o.Sudo)
			defer wg.Done()

			log.Info("Save current userlist to ", DefaultUserlistOldPath)
			if err := scp.SaveOld(c.Context(), DefaultUserlistOldPath, DefaultUserlistRemotePath); err != nil {
				log.Error(err)
				errCh <- err
				return
			}
			readers, err := openFiles(DefaultUserlistOldPath, o.File)
			if err != nil {
				log.Error(err)
				errCh <- err
				return

			}
			log.Info("Compare userlist between ", DefaultUserlistOldPath, " and ", o.File)
			if err := scp.CompareFiles(readers["old"], readers["new"]); err != nil {
				if err != sendfile.ErrorDiff {
					errCh <- err
					return
				}

				log.Info("Copy new userlist to ", pgHost.Host)
				if err := scp.Copy(c.Context(), o.File, o.DestinationFile); err != nil {
					log.Error(err)
					errCh <- err
					return
				}

				return
			}

			log.Info("No changes found in userlist for host ", pgHost.Host)
		}(host)
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

func openFiles(old, new string) (map[string]io.Reader, error) {

	oldReader, err := os.Open(old)
	if err != nil {
		return nil, err
	}

	newReader, err := os.Open(new)
	if err != nil {
		return nil, err
	}

	return map[string]io.Reader{
		"old": oldReader,
		"new": newReader,
	}, nil
}
