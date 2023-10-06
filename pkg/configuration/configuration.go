package configuration

import (
	"io"
	"os"
)

type Configurations interface {
	WriteToFile() (int, error)
	WithPrivKeyFromFile(filePath string) (*Configuration, error)
	GetPostgresDSN() (string, error)
	GetPostgresCustomDSN(dbname, host, sslmode string, port int64) (string, error)
	GetPGBouncerHost() ([]*PGBouncerHost, error)
}

func NewConfiguration(file io.Reader) Configurations {
	return &Configuration{
		configStream: file,
	}
}

func NewConfigurationFromFile(confPath string) (Configurations, error) {
	f, err := os.Open(confPath)
	if (err != nil) && err != os.ErrNotExist {
		return nil, err
	}

	if err == os.ErrNotExist {
		return nil, FileNotFoundFunc()
	}

	return NewConfiguration(f), nil
}

func NewDefaultConfiguration(username, dbname, pghost, password string, pgbouncerHost ...string) Configurations {
	conf := &Configuration{
		Postgrescred: &PostGresCred{
			Host:     pghost,
			UserName: username,
			DBName:   dbname,
			Port:     DefaultPGPort,
			SSLmode:  DefaultSSMode,
			Password: password,
		},
		PGbouncerHosts: func(hostnames ...string) []*PGBouncerHost {
			hosts := []*PGBouncerHost{}
			for _, h := range hostnames {
				pghost := new(PGBouncerHost)
				pghost.Host = h
				hosts = append(hosts,pghost)
			}
			return hosts
		}(pgbouncerHost...),
	}

	return conf.GenerateDefaultConfig()
}
