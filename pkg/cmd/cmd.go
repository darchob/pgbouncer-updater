package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/aio"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/config"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/copy"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/list"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/options"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/cmd/reload"
)

const (
	example = `
	# Write default config in current dir
	%[1]s config

	# List users to a file with default file name userlist.txt
	%[1]s list

	# Copy userlist.txt to hosts
	%[1]s copy

	# Reload PGBouncer hosts
	%[1]s reload
 `
)

// NewCmdArgoRollouts returns new instance of rollouts command.
func NewCmdPGBouncerUpdate(o *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pgbouncer-update COMMAND",
		Short:        "Update PGBouncer Userlist from DB to server files",
		Long:         "List all user roles and md5 password from database write it to userlist.txt and next send they to PGBouncer server.",
		Example:      o.Exemple(example),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return o.UsageErr(c)
		},
	}

	cmd.AddCommand(config.NewCmdConfig(o))
	cmd.AddCommand(aio.NewCmdAIO(o))

	//WARN reload func must be called
	// before all other func use options query vars
	cmd.AddCommand(reload.NewCmdReload(o))
	cmd.AddCommand(list.NewCmdUpdateUserList(o))
	cmd.AddCommand(copy.NewCmdCopyUserList(o))

	return cmd
}
