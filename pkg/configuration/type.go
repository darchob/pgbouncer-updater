package configuration

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

const (
	dsn                = "dbname=%s host=%s port=%d user=%s password=%s sslmode=%s"
	DefaultPGPort      = 5432
	DefaultSSHPort     = 22
	DefaultSSHUsername = "ansible"
	DefaultSSMode      = "disable"
	DefaultFileName    = "config.yaml"
	DefaultPrivKeyPath = "%s/.ssh/id_rsa_ansible"
)

type Configuration struct {
	configStream   io.Reader
	Postgrescred   *PostGresCred    `yaml:"credentials"`
	PGbouncerHosts []*PGBouncerHost `yaml:"hosts"`
}

type PostGresCred struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
	UserName string `yaml:"username"`
}

type PGBouncerHost struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	UserName string `yaml:"username"`
	PrivKey  string `yaml:"privkey"`
}

func (conf *Configuration) GetPGBouncerHost() ([]*PGBouncerHost, error) {
	if err := conf.parseConfigFile(); err != nil {
		return nil, err
	}
	return conf.PGbouncerHosts, nil
}

func (conf *Configuration) GetPostgresDSN() (string, error) {
	if err := conf.parseConfigFile(); err != nil {
		return "", err
	}
	cred := conf.Postgrescred

	return fmt.Sprintf(dsn, cred.DBName, cred.Host, cred.Port, cred.UserName, cred.Password, cred.SSLmode), nil
}

func (conf *Configuration) GetPostgresCustomDSN(dbname, host, sslmode string, port int64) (string, error) {
	if err := conf.parseConfigFile(); err != nil {
		return "", err
	}

	cred := conf.Postgrescred
	return fmt.Sprintf(dsn, dbname, host, port, cred.UserName, cred.Password, sslmode), nil
}

func (conf *Configuration) GenerateDefaultConfig() *Configuration {
	defaultConf := &Configuration{
		Postgrescred: &PostGresCred{
			Host:     conf.Postgrescred.Host,
			UserName: conf.Postgrescred.UserName,
			DBName:   conf.Postgrescred.DBName,
			Password: conf.Postgrescred.Password,
			Port:     DefaultPGPort,
			SSLmode:  DefaultSSMode,
		},
		PGbouncerHosts: func(host []*PGBouncerHost) []*PGBouncerHost {
			hosts := []*PGBouncerHost{}
			for _, h := range host {
				pghost := PGBouncerHost{
					Host:     h.Host,
					Port:     DefaultSSHPort,
					UserName: DefaultSSHUsername,
				}
				hosts = append(hosts, &pghost)
			}

			return hosts
		}(conf.PGbouncerHosts),
	}

	return defaultConf
}

func (conf *Configuration) WriteToFile() (int, error) {
	f, err := os.Create(DefaultFileName)
	if err != nil {
		return 0, err
	}

	config, err := yaml.Marshal(conf)
	if err != nil {
		return 0, err
	}

	return fmt.Fprint(f, string(config))

}
func (conf *Configuration) WithPrivKeyFromFile(filePath string) (*Configuration, error) {
	if filePath == DefaultPrivKeyPath {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		filePath = fmt.Sprintf(DefaultPrivKeyPath, home)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	for _, host := range conf.PGbouncerHosts {
		host.PrivKey = buf.String()
	}

	return conf.DeepCopy(), nil
}

func (conf *Configuration) parseConfigFile() error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(conf.configStream)

	err := yaml.Unmarshal(buf.Bytes(), conf)
	if err != nil {
		return err
	}

	return nil
}

func (in *Configuration) deepCopyInto(out *Configuration) {
	*out = *in
}

func (in *Configuration) DeepCopy() *Configuration {
	if in == nil {
		return nil
	}

	out := new(Configuration)
	in.deepCopyInto(out)
	return out
}

func (in *PGBouncerHost) deepCopyInto(out *PGBouncerHost) {
	*out = *in
}

func (in *PGBouncerHost) DeepCopy() *PGBouncerHost {
	if in == nil {
		return nil
	}

	out := new(PGBouncerHost)
	in.deepCopyInto(out)
	return out
}

func (in *PostGresCred) deepCopyInto(out *PostGresCred) {
	*out = *in
}

func (in *PostGresCred) DeepCopy() *PostGresCred {
	if in == nil {
		return nil
	}

	out := new(PostGresCred)
	in.deepCopyInto(out)
	return out
}
