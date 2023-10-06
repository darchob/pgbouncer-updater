package sendfile

import (
	"context"
	"io"
)

type Scp interface {
	Copy(context context.Context, sourceFile, destinationPath string) error
	CompareFiles(currentFile, oldFile io.Reader) error
	SaveOld(ctx context.Context, filePath, remotePath string) error
	Close()
}

func NewScpClient(host, username string, port int64, privKey []byte, sudo bool) Scp {
	return &Host{
		Username: username,
		Host:     host,
		PrivKey:  privKey,
		remoteBinary: func(sudo bool) string {
			if sudo {
				return "sudo /usr/bin/scp"
			}
			return "/usr/bin/scp"
		}(sudo),
	}
}
