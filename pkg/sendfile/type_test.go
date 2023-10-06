package sendfile

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestHost_Copy(t *testing.T) {
	type fields struct {
		conn     *ssh.Client
		session  *ssh.Session
		Host     string
		Username string
		PrivKey  []byte
		file     *File
	}
	type args struct {
		sourceFile      string
		destinationFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				Username: "jobd",
				Host:     "localhost:22",
				PrivKey: []byte(`
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAsoKKybf4V2NEALrE/GqgSw+znYCXqyzhHcKdnR6wuqzdb7UDMZi4
P1lPJt5O68mGKx8qlb39e55svixeWccoZFpDnpy75YHxVD1pxJhNZxOFQ2JtYc/Eh+WYZ3
uh2ct1qzt2NRe55A/05u+l2/rWSa+dFT4yJd52Qn85yPENSgYrkxo4kfFfxaVzPVvsjHXk
HRvFZy0YDhAB5FaFt5cbKoVMnJSXLd7xJ8RzWW5VNycOQi7sYBa1itbMujhs6sKxf++M5k
1QIRsFMZnkSmxxtLn2WkBdc4Pjxg2W/NJ9Of7qKAoRPcyynsnsxRur6fdd26+QujoxE9wM
N20Iq5wdnwAAA8g6IuVtOiLlbQAAAAdzc2gtcnNhAAABAQCygorJt/hXY0QAusT8aqBLD7
OdgJerLOEdwp2dHrC6rN1vtQMxmLg/WU8m3k7ryYYrHyqVvf17nmy+LF5ZxyhkWkOenLvl
gfFUPWnEmE1nE4VDYm1hz8SH5Zhne6HZy3WrO3Y1F7nkD/Tm76Xb+tZJr50VPjIl3nZCfz
nI8Q1KBiuTGjiR8V/FpXM9W+yMdeQdG8VnLRgOEAHkVoW3lxsqhUyclJct3vEnxHNZblU3
Jw5CLuxgFrWK1sy6OGzqwrF/74zmTVAhGwUxmeRKbHG0ufZaQF1zg+PGDZb80n05/uooCh
E9zLKeyezFG6vp913br5C6OjET3Aw3bQirnB2fAAAAAwEAAQAAAQA2GlikMK0FF2Hp8rF3
a32vok+nAe12BQEpuu14TG/19CSdEbipFIdrM89IkYJL9mVCtox6m/2ytN5yeRITlcgJOk
5aSVitg8e353EiE6MKBaGTPca3KXiAU7bwTklMsFy2jCwUhV9i3u8z+xhC5vCBnsc2RAaA
8b7YAqVp4J1NfLT4myYDlUG93+GyrLWADl4XLY+FgZdROr7tGxnlifbcCeTHb+c7u3X4B1
spKRi+6zlqaQ1JkjKAhK43oho3eSEPp6v5mohP486WAVmLIHqM5JXGhx9V3o9JGnEXUuu7
FYAMTbuHGGiVQAAXuzgCUaxZTzfMmnwO+hY2iNjYCXrRAAAAgQC3TIWCWIHyUMHjt1Ei+5
mg4DTs6T03tr7Ns4s9gI3s6waNZYakHRg5VqfohslenPoD/TiqvOZxkfOBV+4y9LLL14vv
EwEb5Sn+jWcAYzM3dI2y8Cwiou0Cc8YprGWkwKljjteg58EVp3srlHg4yn52YdrX6bFTOY
DTZPonRRqU8wAAAIEA7US58QJwruBFDMQ2v3dTUYCFrGs6xMegHbULfjtLkTYu3tatvttr
10J61ifGNuqvVbTbQXJO2ikIZFH4XW5WC4KN8UjA5JkJL++w5Md9/s+jfBrW7bbbXYEMfq
0wZ4tanXtoD5hjHwf78aXdJobi+Kyhaov+pfxhQ41phxMFJycAAACBAMCaSVsuE6BLUJrR
ee1drAmJqG25SXPOt4xXEOPhPOfR3ViAzIFvF2ELqihhe1h9NidlYUgX60x+SXtCPRKFLx
mWa5NF7WsZy8Qk7NuFbnWcCQbEaEqKnUv67QnNR+eX9+QnI4kueupYMt2zmOWkh80dezIN
kezKtxKihlsqKKDJAAAAEGpvYmRAUkQtTElOVVgtMDEBAg==
-----END OPENSSH PRIVATE KEY-----
`),
			},
			args: args{
				destinationFile: "/tmp/test.yaml",
				sourceFile:      "test.yaml",
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			h := &Host{
				conn:         tt.fields.conn,
				session:      tt.fields.session,
				Host:         tt.fields.Host,
				Username:     tt.fields.Username,
				PrivKey:      tt.fields.PrivKey,
				file:         tt.fields.file,
				remoteBinary: "/usr/bin/scp",
				timeout:      60 * time.Second,
			}
			defer h.Close()

			if err := h.Copy(ctx, tt.args.sourceFile, tt.args.destinationFile); (err != nil) != tt.wantErr {
				t.Errorf("Host.Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHost_CompareFiles(t *testing.T) {

	type args struct {
		currentFile io.Reader
		oldFile     io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				currentFile: strings.NewReader("TestOne"),
				oldFile:     strings.NewReader("TestOne"),
			},
			wantErr: false,
		},
		{
			args: args{
				currentFile: strings.NewReader("TestOne"),
				oldFile:     strings.NewReader("TestTwo"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Host{}
			if err := h.CompareFiles(tt.args.currentFile, tt.args.oldFile); (err != nil) != tt.wantErr {
				t.Errorf("Host.CompareFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
