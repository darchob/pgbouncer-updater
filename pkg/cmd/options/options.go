package options

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/databases"
)

type GetOptions interface {
	WithDefaultOptions() *Options
	WithQuery(query string) *Options
	AsSudoer() *Options
	WithDestinationFile(pathFile string) *Options
	WithFile(srcFilePath string) *Options
	WithConfigFilePath(configPath string) *Options
}

const (
	cliName = "pgbouncer-updater"
)

type Options struct {
	Query           string
	File            string
	ConfigFilePath  string
	Log             *log.Logger
	LogLevel        string
	DestinationFile string
	Sudo            bool
	UserName        string
	PGHost          string
	DBName          string
	Password        string
	PGBouncerHosts  []string
}

func NewPGBouncerUpdaterOptions() GetOptions {
	return new(Options)
}

func (o *Options) WithDefaultFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Password, "password", "", "DB password")
	cmd.Flags().StringVar(&o.UserName, "username", "", "DB Username")
	cmd.Flags().StringVar(&o.DBName, "dbname", "", "DB name")
	cmd.Flags().StringVar(&o.PGHost, "pghost", "", "DB Hostname")
	cmd.Flags().StringArrayVar(&o.PGBouncerHosts, "pgbouncerhosts", nil, "PGBOuncer hostnames")
}

func (o *Options) WithDefaultOptions() *Options {
	o = &Options{
		Sudo:            false,
		DestinationFile: "/etc/pgbouncer/userlist.txt",
		Query:           databases.DefaultQuery,
		ConfigFilePath:  "/etc/pgbouncer-updater/config.yaml",
		File:            "/tmp/userlist.txt",
	}

	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) WithQuery(query string) *Options {
	o.Query = query
	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) WithConfigFilePath(configPath string) *Options {
	o.ConfigFilePath = configPath
	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) WithDestinationFile(dst string) *Options {
	o.DestinationFile = dst
	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) WithFile(srcFilePath string) *Options {
	o.File = srcFilePath
	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) AsSudoer() *Options {
	o.Sudo = true
	return o.newPGBouncerUpdaterOptions()
}

func (o *Options) newPGBouncerUpdaterOptions() *Options {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	defaultOpts := &Options{
		Query:           o.Query,
		Log:             logCtx,
		LogLevel:        log.InfoLevel.String(),
		DestinationFile: o.DestinationFile,
		Sudo:            o.Sudo,
		ConfigFilePath:  o.ConfigFilePath,
		File:            o.File,
	}

	return defaultOpts

}

func (o *Options) UsageErr(c *cobra.Command) error {
	c.Usage()
	c.SilenceErrors = true
	return errors.New(c.UsageString())
}

func (o *Options) Exemple(example string) string {
	return strings.Trim(fmt.Sprintf(example, cliName), "\n")
}
