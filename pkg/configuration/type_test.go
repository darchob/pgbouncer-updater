package configuration

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

func createStream() (io.Reader, error) {
	cred := PostGresCred{
		Host:     "10.29.0.0",
		UserName: "postgres",
		Password: "password",
		DBName:   "databaseName",
		SSLmode:  "allow",
		Port:     5432,
	}

	conf := Configuration{
		configStream: nil,
		Postgrescred: &cred,
	}

	out, err := yaml.Marshal(conf)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

func TestConfiguration_GetPostgresDSN(t *testing.T) {
	stream, err := createStream()
	if err != nil {
		t.Errorf("Error while generate a test stream %s", err)
		return
	}

	type fields struct {
		configStream io.Reader
		postgrescred *PostGresCred
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Test format",
			fields: fields{
				configStream: stream,
				postgrescred: &PostGresCred{
					Host:     "10.29.0.0",
					UserName: "postgres",
					Password: "password",
					DBName:   "databaseName",
					SSLmode:  "allow",
					Port:     5432,
				},
			},
			want: "dbname=databaseName host=10.29.0.0 port=5432 user=postgres password=password sslmode=allow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Configuration{
				configStream: tt.fields.configStream,
				Postgrescred: tt.fields.postgrescred,
			}
			got, err := conf.GetPostgresDSN()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration.GetPostgresDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Configuration.GetPostgresDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_GetPGBouncerHost(t *testing.T) {
	stream, err := createStream()
	if err != nil {
		t.Errorf("Error while generate a test stream %s", err)
		return
	}

	type fields struct {
		configStream   io.Reader
		postgrescred   *PostGresCred
		pgbouncerHosts []*PGBouncerHost
	}
	tests := []struct {
		name   string
		fields fields
		want   []*PGBouncerHost
	}{
		{
			name: "Test format",
			fields: fields{
				configStream: stream,
				pgbouncerHosts: []*PGBouncerHost{
					{
						Host: "localhost",
						Port: 22,
						PrivKey: `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAzl5oDetjducGHBtrh2PdUvJn7CZlb+wsU5zQ3ILwsVRUH6he9w4A
bEIlLwhYc0sBWcJiTcopNMPQrQjx5iU8Owp68QqkI/Vreky35pcMTUIPGE44j8aGX2CBGs
q7fUO8Wh9lO3fqGgLGJeI5CDEhoqZbHXt+zpCxIMtht798NP5evqW8Zp0tFq6PF1G69x4D
B+qhvs8htfxHtZx9DNfMqhRAlmQdmIhoKnJ46hVPQXIvc/9LnFlDLJ6oEr787jQGlYOao+
G8r/fzhC04HqBoey9z/s61YpLxGzoJFA7rXz/GoBkEhIkyitZhO/Is9g5b18k3ieqNYqU6
ws47HFDfEQAAA8hTSVMqU0lTKgAAAAdzc2gtcnNhAAABAQDOXmgN62N25wYcG2uHY91S8m
fsJmVv7CxTnNDcgvCxVFQfqF73DgBsQiUvCFhzSwFZwmJNyik0w9CtCPHmJTw7CnrxCqQj
9Wt6TLfmlwxNQg8YTjiPxoZfYIEayrt9Q7xaH2U7d+oaAsYl4jkIMSGiplsde37OkLEgy2
G3v3w0/l6+pbxmnS0Wro8XUbr3HgMH6qG+zyG1/Ee1nH0M18yqFECWZB2YiGgqcnjqFU9B
ci9z/0ucWUMsnqgSvvzuNAaVg5qj4byv9/OELTgeoGh7L3P+zrVikvEbOgkUDutfP8agGQ
SEiTKK1mE78iz2DlvXyTeJ6o1ipTrCzjscUN8RAAAAAwEAAQAAAQBOzdoqRpLK2tmIbigX
oVjozcxFbzwZCzS6EQ3oxs+mx68AD8mDygL7VB7i4Or1y9SONB5Z2jL2BThwexP0cI+ZdB
0SYp/fY15Ra25mTZPTBMDC5UvQC11Qmodydaw232DTgV2k4duxZxHHcaWZrTlM5P2yOnBn
7PTWsxNzmVmS2kdhf1bmAunjcbla6pXxUvHcfynGm3k2H8blmeUo/ZfuZMDEmBGwxuU4Xr
LjVRtJVsknrm/2c5aowLdWjo03YZzLs9MX0bFSCxfGZ8MtWXNxOKIUC3iCAgchCX27Jm3Z
h9QgyB71KMVdtv2w/9skl831TFcvmP3V+qsTRNuV1PBpAAAAgQCai/Yq2mjetvYmGTaXMH
T8INyA3RL+z5UtLlPaTWtT4xT3ul6pQNBt6SJwTPcYqLnZXDZU5OYW35esiwDwiUS/KYY4
nio+83colesLyntG4AQDnQojowevYQWao+RCuIO/angxrqqqN6EIKXVt6oLG+bsJ5+Odw6
HHQWy8yVwVHgAAAIEA/ifZlqSCMixy19tfDKGmKof6uYyf033q4xL8BvWkICIIxnY9wkIj
kAqFu/+59uYXRNUKHt91jXL8I5X2ZgQW2cTGiEVZdY+akDr9SDu71y0aqb9D9W7cf6tunM
OGyKyKgkttJ0dsbKL2JgDk/Jc6/1etbueXm9WibOP6EyQz15sAAACBAM/dyCdsCZdb8mIu
Qv4/Hj6tI/ODFhTwTyCI1zi0lUXpHU4gJCov3TkwOy8tVrcTgNzD/f8A5hzVZ1mXyko3Jy
105Jw+ht//J/OWHeWQUKSXOGp0GuV1SvD8Zq0uwLbR1D6yx2eBke9PKxiNoKdqRde7hIHJ
QVBd+KAa2vXMHyzDAAAAEGpvYmRAUkQtTElOVVgtMDEBAg==
-----END OPENSSH PRIVATE KEY-----
`,
					},
				},
			},
			want: []*PGBouncerHost{
				{
					Host: "localhost",
					Port: 22,
					PrivKey: `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAzl5oDetjducGHBtrh2PdUvJn7CZlb+wsU5zQ3ILwsVRUH6he9w4A
bEIlLwhYc0sBWcJiTcopNMPQrQjx5iU8Owp68QqkI/Vreky35pcMTUIPGE44j8aGX2CBGs
q7fUO8Wh9lO3fqGgLGJeI5CDEhoqZbHXt+zpCxIMtht798NP5evqW8Zp0tFq6PF1G69x4D
B+qhvs8htfxHtZx9DNfMqhRAlmQdmIhoKnJ46hVPQXIvc/9LnFlDLJ6oEr787jQGlYOao+
G8r/fzhC04HqBoey9z/s61YpLxGzoJFA7rXz/GoBkEhIkyitZhO/Is9g5b18k3ieqNYqU6
ws47HFDfEQAAA8hTSVMqU0lTKgAAAAdzc2gtcnNhAAABAQDOXmgN62N25wYcG2uHY91S8m
fsJmVv7CxTnNDcgvCxVFQfqF73DgBsQiUvCFhzSwFZwmJNyik0w9CtCPHmJTw7CnrxCqQj
9Wt6TLfmlwxNQg8YTjiPxoZfYIEayrt9Q7xaH2U7d+oaAsYl4jkIMSGiplsde37OkLEgy2
G3v3w0/l6+pbxmnS0Wro8XUbr3HgMH6qG+zyG1/Ee1nH0M18yqFECWZB2YiGgqcnjqFU9B
ci9z/0ucWUMsnqgSvvzuNAaVg5qj4byv9/OELTgeoGh7L3P+zrVikvEbOgkUDutfP8agGQ
SEiTKK1mE78iz2DlvXyTeJ6o1ipTrCzjscUN8RAAAAAwEAAQAAAQBOzdoqRpLK2tmIbigX
oVjozcxFbzwZCzS6EQ3oxs+mx68AD8mDygL7VB7i4Or1y9SONB5Z2jL2BThwexP0cI+ZdB
0SYp/fY15Ra25mTZPTBMDC5UvQC11Qmodydaw232DTgV2k4duxZxHHcaWZrTlM5P2yOnBn
7PTWsxNzmVmS2kdhf1bmAunjcbla6pXxUvHcfynGm3k2H8blmeUo/ZfuZMDEmBGwxuU4Xr
LjVRtJVsknrm/2c5aowLdWjo03YZzLs9MX0bFSCxfGZ8MtWXNxOKIUC3iCAgchCX27Jm3Z
h9QgyB71KMVdtv2w/9skl831TFcvmP3V+qsTRNuV1PBpAAAAgQCai/Yq2mjetvYmGTaXMH
T8INyA3RL+z5UtLlPaTWtT4xT3ul6pQNBt6SJwTPcYqLnZXDZU5OYW35esiwDwiUS/KYY4
nio+83colesLyntG4AQDnQojowevYQWao+RCuIO/angxrqqqN6EIKXVt6oLG+bsJ5+Odw6
HHQWy8yVwVHgAAAIEA/ifZlqSCMixy19tfDKGmKof6uYyf033q4xL8BvWkICIIxnY9wkIj
kAqFu/+59uYXRNUKHt91jXL8I5X2ZgQW2cTGiEVZdY+akDr9SDu71y0aqb9D9W7cf6tunM
OGyKyKgkttJ0dsbKL2JgDk/Jc6/1etbueXm9WibOP6EyQz15sAAACBAM/dyCdsCZdb8mIu
Qv4/Hj6tI/ODFhTwTyCI1zi0lUXpHU4gJCov3TkwOy8tVrcTgNzD/f8A5hzVZ1mXyko3Jy
105Jw+ht//J/OWHeWQUKSXOGp0GuV1SvD8Zq0uwLbR1D6yx2eBke9PKxiNoKdqRde7hIHJ
QVBd+KAa2vXMHyzDAAAAEGpvYmRAUkQtTElOVVgtMDEBAg==
-----END OPENSSH PRIVATE KEY-----
`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Configuration{
				configStream:   tt.fields.configStream,
				Postgrescred:   tt.fields.postgrescred,
				PGbouncerHosts: tt.fields.pgbouncerHosts,
			}
			if got, err := conf.GetPGBouncerHost(); !reflect.DeepEqual(got, tt.want) && err != nil {
				t.Errorf("Configuration.GetPGBouncerHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_GetPostgresCustomDSN(t *testing.T) {
	stream, err := createStream()
	if err != nil {
		t.Errorf("Error while generate a test stream %s", err)
		return
	}

	type fields struct {
		configStream   io.Reader
		Postgrescred   *PostGresCred
		PGbouncerHosts []*PGBouncerHost
	}
	type args struct {
		dbname  string
		host    string
		sslmode string
		port    int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test format",
			args: args{
				dbname:  "databaseName",
				host:    "10.29.0.0",
				sslmode: "allow",
				port:    5432,
			},
			fields: fields{
				configStream: stream,
			},
			want: "dbname=databaseName host=10.29.0.0 port=5432 user=postgres password=password sslmode=allow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Configuration{
				configStream: tt.fields.configStream,
			}
			got, err := conf.GetPostgresCustomDSN(tt.args.dbname, tt.args.host, tt.args.sslmode, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration.GetPostgresCustomDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Configuration.GetPostgresCustomDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_GenerateDefaultConfigFile(t *testing.T) {
	stream, err := createStream()
	if err != nil {
		t.Errorf("Error while generate a test stream %s", err)
		return
	}

	type fields struct {
		configStream   io.Reader
		Postgrescred   *PostGresCred
		PGbouncerHosts []*PGBouncerHost
	}
	tests := []struct {
		name   string
		fields fields
		want   *Configuration
	}{
		{
			fields: fields{
				configStream: stream,
				Postgrescred: &PostGresCred{
					Host:     "10.29.0.0",
					UserName: "postgres",
					Password: "password",
					DBName:   "databaseName",
					SSLmode:  "disable",
					Port:     5432,
				},
				PGbouncerHosts: func(hostnames ...string) []*PGBouncerHost {
					hosts := []*PGBouncerHost{}
					for _, h := range hostnames {
						pghost := PGBouncerHost{
							Host:     h,
							Port:     DefaultSSHPort,
							UserName: "ansible",
						}
						hosts = append(hosts, &pghost)
					}
					return hosts
				}("pgbouncer-01", "pgbouncer-02"),
			},
			want: &Configuration{
				Postgrescred: &PostGresCred{
					Host:     "10.29.0.0",
					UserName: "postgres",
					Password: "password",
					DBName:   "databaseName",
					SSLmode:  "disable",
					Port:     5432,
				},
				PGbouncerHosts: []*PGBouncerHost{
					{
						Host:     "pgbouncer-01",
						Port:     DefaultSSHPort,
						UserName: "ansible",
						PrivKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAzl5oDetjducGHBtrh2PdUvJn7CZlb+wsU5zQ3ILwsVRUH6he9w4A
bEIlLwhYc0sBWcJiTcopNMPQrQjx5iU8Owp68QqkI/Vreky35pcMTUIPGE44j8aGX2CBGs
q7fUO8Wh9lO3fqGgLGJeI5CDEhoqZbHXt+zpCxIMtht798NP5evqW8Zp0tFq6PF1G69x4D
B+qhvs8htfxHtZx9DNfMqhRAlmQdmIhoKnJ46hVPQXIvc/9LnFlDLJ6oEr787jQGlYOao+
G8r/fzhC04HqBoey9z/s61YpLxGzoJFA7rXz/GoBkEhIkyitZhO/Is9g5b18k3ieqNYqU6
ws47HFDfEQAAA8hTSVMqU0lTKgAAAAdzc2gtcnNhAAABAQDOXmgN62N25wYcG2uHY91S8m
fsJmVv7CxTnNDcgvCxVFQfqF73DgBsQiUvCFhzSwFZwmJNyik0w9CtCPHmJTw7CnrxCqQj
9Wt6TLfmlwxNQg8YTjiPxoZfYIEayrt9Q7xaH2U7d+oaAsYl4jkIMSGiplsde37OkLEgy2
G3v3w0/l6+pbxmnS0Wro8XUbr3HgMH6qG+zyG1/Ee1nH0M18yqFECWZB2YiGgqcnjqFU9B
ci9z/0ucWUMsnqgSvvzuNAaVg5qj4byv9/OELTgeoGh7L3P+zrVikvEbOgkUDutfP8agGQ
SEiTKK1mE78iz2DlvXyTeJ6o1ipTrCzjscUN8RAAAAAwEAAQAAAQBOzdoqRpLK2tmIbigX
oVjozcxFbzwZCzS6EQ3oxs+mx68AD8mDygL7VB7i4Or1y9SONB5Z2jL2BThwexP0cI+ZdB
0SYp/fY15Ra25mTZPTBMDC5UvQC11Qmodydaw232DTgV2k4duxZxHHcaWZrTlM5P2yOnBn
7PTWsxNzmVmS2kdhf1bmAunjcbla6pXxUvHcfynGm3k2H8blmeUo/ZfuZMDEmBGwxuU4Xr
LjVRtJVsknrm/2c5aowLdWjo03YZzLs9MX0bFSCxfGZ8MtWXNxOKIUC3iCAgchCX27Jm3Z
h9QgyB71KMVdtv2w/9skl831TFcvmP3V+qsTRNuV1PBpAAAAgQCai/Yq2mjetvYmGTaXMH
T8INyA3RL+z5UtLlPaTWtT4xT3ul6pQNBt6SJwTPcYqLnZXDZU5OYW35esiwDwiUS/KYY4
nio+83colesLyntG4AQDnQojowevYQWao+RCuIO/angxrqqqN6EIKXVt6oLG+bsJ5+Odw6
HHQWy8yVwVHgAAAIEA/ifZlqSCMixy19tfDKGmKof6uYyf033q4xL8BvWkICIIxnY9wkIj
kAqFu/+59uYXRNUKHt91jXL8I5X2ZgQW2cTGiEVZdY+akDr9SDu71y0aqb9D9W7cf6tunM
OGyKyKgkttJ0dsbKL2JgDk/Jc6/1etbueXm9WibOP6EyQz15sAAACBAM/dyCdsCZdb8mIu
Qv4/Hj6tI/ODFhTwTyCI1zi0lUXpHU4gJCov3TkwOy8tVrcTgNzD/f8A5hzVZ1mXyko3Jy
105Jw+ht//J/OWHeWQUKSXOGp0GuV1SvD8Zq0uwLbR1D6yx2eBke9PKxiNoKdqRde7hIHJ
QVBd+KAa2vXMHyzDAAAAEGpvYmRAUkQtTElOVVgtMDEBAg==
-----END OPENSSH PRIVATE KEY-----`,
					},
					{
						Host:     "pgbouncer-02",
						Port:     DefaultSSHPort,
						UserName: "ansible",
						PrivKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAzl5oDetjducGHBtrh2PdUvJn7CZlb+wsU5zQ3ILwsVRUH6he9w4A
bEIlLwhYc0sBWcJiTcopNMPQrQjx5iU8Owp68QqkI/Vreky35pcMTUIPGE44j8aGX2CBGs
q7fUO8Wh9lO3fqGgLGJeI5CDEhoqZbHXt+zpCxIMtht798NP5evqW8Zp0tFq6PF1G69x4D
B+qhvs8htfxHtZx9DNfMqhRAlmQdmIhoKnJ46hVPQXIvc/9LnFlDLJ6oEr787jQGlYOao+
G8r/fzhC04HqBoey9z/s61YpLxGzoJFA7rXz/GoBkEhIkyitZhO/Is9g5b18k3ieqNYqU6
ws47HFDfEQAAA8hTSVMqU0lTKgAAAAdzc2gtcnNhAAABAQDOXmgN62N25wYcG2uHY91S8m
fsJmVv7CxTnNDcgvCxVFQfqF73DgBsQiUvCFhzSwFZwmJNyik0w9CtCPHmJTw7CnrxCqQj
9Wt6TLfmlwxNQg8YTjiPxoZfYIEayrt9Q7xaH2U7d+oaAsYl4jkIMSGiplsde37OkLEgy2
G3v3w0/l6+pbxmnS0Wro8XUbr3HgMH6qG+zyG1/Ee1nH0M18yqFECWZB2YiGgqcnjqFU9B
ci9z/0ucWUMsnqgSvvzuNAaVg5qj4byv9/OELTgeoGh7L3P+zrVikvEbOgkUDutfP8agGQ
SEiTKK1mE78iz2DlvXyTeJ6o1ipTrCzjscUN8RAAAAAwEAAQAAAQBOzdoqRpLK2tmIbigX
oVjozcxFbzwZCzS6EQ3oxs+mx68AD8mDygL7VB7i4Or1y9SONB5Z2jL2BThwexP0cI+ZdB
0SYp/fY15Ra25mTZPTBMDC5UvQC11Qmodydaw232DTgV2k4duxZxHHcaWZrTlM5P2yOnBn
7PTWsxNzmVmS2kdhf1bmAunjcbla6pXxUvHcfynGm3k2H8blmeUo/ZfuZMDEmBGwxuU4Xr
LjVRtJVsknrm/2c5aowLdWjo03YZzLs9MX0bFSCxfGZ8MtWXNxOKIUC3iCAgchCX27Jm3Z
h9QgyB71KMVdtv2w/9skl831TFcvmP3V+qsTRNuV1PBpAAAAgQCai/Yq2mjetvYmGTaXMH
T8INyA3RL+z5UtLlPaTWtT4xT3ul6pQNBt6SJwTPcYqLnZXDZU5OYW35esiwDwiUS/KYY4
nio+83colesLyntG4AQDnQojowevYQWao+RCuIO/angxrqqqN6EIKXVt6oLG+bsJ5+Odw6
HHQWy8yVwVHgAAAIEA/ifZlqSCMixy19tfDKGmKof6uYyf033q4xL8BvWkICIIxnY9wkIj
kAqFu/+59uYXRNUKHt91jXL8I5X2ZgQW2cTGiEVZdY+akDr9SDu71y0aqb9D9W7cf6tunM
OGyKyKgkttJ0dsbKL2JgDk/Jc6/1etbueXm9WibOP6EyQz15sAAACBAM/dyCdsCZdb8mIu
Qv4/Hj6tI/ODFhTwTyCI1zi0lUXpHU4gJCov3TkwOy8tVrcTgNzD/f8A5hzVZ1mXyko3Jy
105Jw+ht//J/OWHeWQUKSXOGp0GuV1SvD8Zq0uwLbR1D6yx2eBke9PKxiNoKdqRde7hIHJ
QVBd+KAa2vXMHyzDAAAAEGpvYmRAUkQtTElOVVgtMDEBAg==
-----END OPENSSH PRIVATE KEY-----`,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Configuration{
				configStream:   tt.fields.configStream,
				Postgrescred:   tt.fields.Postgrescred,
				PGbouncerHosts: tt.fields.PGbouncerHosts,
			}

			got, err := conf.GenerateDefaultConfig().WithPrivKeyFromFile(DefaultPrivKeyPath)
			if err != nil {
				t.Errorf("Configuration.GenerateDefaultConfigFile().WithPrivKeyFromFile(DefaultPrivKeyPath) error %v", err)
			}
			if !reflect.DeepEqual(got.Postgrescred, tt.want.Postgrescred) {
				t.Errorf("Configuration.GenerateDefaultConfigFile() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.PGbouncerHosts, tt.want.PGbouncerHosts) {
				t.Errorf("Configuration.GenerateDefaultConfigFile() = %v, want %v", got, tt.want)
			}

		})
	}
}
