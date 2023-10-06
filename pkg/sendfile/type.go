package sendfile

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	ErrorDiff = fmt.Errorf("Content are not same")
)

type Host struct {
	conn         *ssh.Client
	session      *ssh.Session
	Host         string
	Username     string
	PrivKey      []byte
	file         *File
	remoteBinary string
	timeout      time.Duration
}

type File struct {
	fileName string
	mode     os.FileMode
	content  io.Reader
	size     int64
}

func (h *Host) CompareFiles(currentFilePath, oldFilePath io.Reader) error {
	old, err := md5sum(oldFilePath)
	if err != nil {
		return err
	}

	new, err := md5sum(currentFilePath)
	if err != nil {
		return err
	}

	if !bytes.Equal(old, new) {
		return ErrorDiff
	}

	return nil
}

func (h *Host) SaveOld(ctx context.Context, filePath, remotePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	if err := h.connect(); err != nil {
		return err
	}

	return h.copyFrom(ctx, f, remotePath)
}

func (h *Host) Copy(ctx context.Context, sourceFile, destinationFile string) error {
	var err error

	file, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer file.Close()

	s, err := file.Stat()
	if err != nil {
		return err
	}

	if s.Size() <= 0 {
		return fmt.Errorf("empty file")
	}

	h.file = &File{
		fileName: filepath.Base(destinationFile),
		mode:     s.Mode().Perm(),
		size:     s.Size(),
		content:  file,
	}

	if err := h.connect(); err != nil {
		return err
	}

	return h.copy(ctx, destinationFile)
}
func (h *Host) auth() (*ssh.ClientConfig, error) {
	signer, err := ssh.ParsePrivateKey(h.PrivKey)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User:            h.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}, nil
}

func (h *Host) connect() error {
	clientConfig, err := h.auth()
	if err != nil {
		return err
	}

	h.conn, err = ssh.Dial("tcp", h.Host, clientConfig)
	if err != nil {
		return err
	}

	h.session, err = h.conn.NewSession()
	if err != nil {
		return err
	}

	return nil
}

func (h *Host) copyFrom(ctx context.Context, w io.Writer, remotePath string) error {
	wg := sync.WaitGroup{}
	errCh := make(chan error, 4)

	wg.Add(1)
	go func() {
		var err error

		defer func() {
			// NOTE: this might send an already sent error another time, but since we only receive opne, this is fine. On the "happy-path" of this function, the error will be `nil` therefore completing the "err<-errCh" at the bottom of the function.
			errCh <- err
			// We must unblock the go routine first as we block on reading the channel later
			wg.Done()

		}()

		r, err := h.session.StdoutPipe()
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		in, err := h.session.StdinPipe()
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}
		defer in.Close()

		err = h.session.Start(fmt.Sprintf("%s -f %q", h.remoteBinary, remotePath))
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		err = Ack(in)
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		res, err := ParseResponse(r)
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}
		if res.IsFailure() {
			log.Error(errors.New(res.GetMessage()))
			errCh <- errors.New(res.GetMessage())
			return
		}

		infos, err := res.ParseFileInfos()
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		err = Ack(in)
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		_, err = copyN(w, r, infos.Size)
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		err = Ack(in)
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}

		err = h.session.Wait()
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		}
	}()

	if h.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
		defer cancel()
	}

	if err := wait(&wg, ctx); err != nil {
		return err
	}
	finalErr := <-errCh
	close(errCh)
	return finalErr

}
func (h *Host) copy(ctx context.Context, dstPath string) error {
	stdout, err := h.session.StdoutPipe()
	if err != nil {
		return err
	}

	w, err := h.session.StdinPipe()
	if err != nil {
		return err
	}
	defer w.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()
		defer w.Close()

		_, err = fmt.Fprintf(w, "C%#o %d %s\n", h.file.mode, h.file.size, h.file.fileName)
		if err != nil {
			errCh <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errCh <- err
			return
		}

		_, err = io.Copy(w, h.file.content)
		if err != nil {
			errCh <- err
			return
		}

		_, err = fmt.Fprint(w, "\x00")
		if err != nil {
			errCh <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		err := h.session.Run(fmt.Sprintf("%s -qt %q", h.remoteBinary, dstPath))
		if err != nil {
			errCh <- err
			return
		}
	}()

	if h.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
		defer cancel()
	}

	if err := wait(&wg, ctx); err != nil {
		return err
	}

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) Close() {
	if h.conn.Close() != nil {
		h.conn.Close()
	}

	if h.session != nil {
		h.session.Close()
	}
}

func wait(wg *sync.WaitGroup, ctx context.Context) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func checkResponse(r io.Reader) error {
	response, err := ParseResponse(r)
	if err != nil {
		return err
	}

	if response.IsFailure() {
		return errors.New(response.GetMessage())
	}

	return nil

}
