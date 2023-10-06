package options

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"gitlab.infra.adelaidegroup.fr/infra/applications-sources/pgbouncer-updater/pkg/databases"
)

func TestWithDefaultOptions(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Query:           databases.DefaultQuery,
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/etc/pgbouncer/userlist.txt",
				ConfigFilePath:  "/etc/pgbouncer-updater/config.yaml",
				File:            "/tmp/userlist.txt",
				Sudo:            false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().WithDefaultOptions()
			if tt.want.Query != got.Query {
				t.Errorf("WithDefaultOptions() = %v, want %v", got, tt.want)
			}

			if tt.want.DestinationFile != got.DestinationFile {
				t.Errorf("WithDefaultOptions() = %v, want %v", got, tt.want)
			}

			if tt.want.Sudo != got.Sudo {
				t.Errorf("WithDefaultOptions() = %v, want %v", got, tt.want)
			}
			if tt.want.ConfigFilePath != got.ConfigFilePath {
				t.Errorf("WithDefaultOptions() = %v, want %v", got, tt.want)
			}
			if tt.want.File != got.File {
				t.Errorf("WithDefaultOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithQuery(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/etc/pgbouncer/userlist.txt",
				ConfigFilePath:  "/opt/config.yaml",
				Query:            "userlist.txt",
				Sudo:            false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().WithQuery("userlist.txt")
			if tt.want.File != got.File {
				t.Errorf("WithQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithConfigFilePath(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Query:           "select",
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/etc/pgbouncer/userlist.txt",
				ConfigFilePath:  "/opt/config.yaml",
				Sudo:            false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().WithConfigFilePath("/opt/config.yaml")
			if tt.want.ConfigFilePath != got.ConfigFilePath {
				t.Errorf("WithConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFile(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Query:           databases.DefaultQuery,
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/etc/pgbouncer/userlist.txt",
				ConfigFilePath:  "/opt/config.yaml",
				File:            "userlist.txt",
				Sudo:            false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().WithFile("userlist.txt")
			if tt.want.File != got.File {
				t.Errorf("WithConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDestinationFile(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Query:           databases.DefaultQuery,
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/tmp/userlist.txt",
				Sudo:            false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().WithDestinationFile("/tmp/userlist.txt")
			if tt.want.DestinationFile != got.DestinationFile {
				t.Errorf("WithDestinantionFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsSudoers(t *testing.T) {
	logCtx := log.New()
	logCtx.SetOutput(os.Stdout)

	tests := []struct {
		name string
		want *Options
	}{
		{
			want: &Options{
				Query:           databases.DefaultQuery,
				Log:             logCtx,
				LogLevel:        log.InfoLevel.String(),
				DestinationFile: "/tmp/userlist.txt",
				Sudo:            true,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := NewPGBouncerUpdaterOptions().AsSudoer()
			if tt.want.Sudo != got.Sudo {
				t.Errorf("AsSudoer() = %v, want %v", got, tt.want)
			}
		})
	}
}

