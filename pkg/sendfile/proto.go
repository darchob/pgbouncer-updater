package sendfile

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type ResponseType = uint8

const (
	Ok      ResponseType = 0
	Warning ResponseType = 1
	Error   ResponseType = 2
)

type Response struct {
	Type    ResponseType
	Message string
}

// ParseResponse reads from the given reader (assuming it is the output of the remote) and parses it into a Response structure.
func ParseResponse(reader io.Reader) (Response, error) {
	buffer := make([]uint8, 1)
	_, err := reader.Read(buffer)
	if err != nil {
		return Response{}, err
	}

	responseType := buffer[0]
	message := ""
	if responseType > 0 {
		bufferedReader := bufio.NewReader(reader)
		message, err = bufferedReader.ReadString('\n')
		if err != nil {
			return Response{}, err
		}
	}

	return Response{responseType, message}, nil
}

func (r *Response) IsOk() bool {
	return r.Type == Ok
}

func (r *Response) IsWarning() bool {
	return r.Type == Warning
}

// IsError returns true when the remote responded with an error.
func (r *Response) IsError() bool {
	return r.Type == Error
}

// IsFailure returns true when the remote answered with a warning or an error.
func (r *Response) IsFailure() bool {
	return r.IsWarning() || r.IsError()
}

// GetMessage returns the message the remote sent back.
func (r *Response) GetMessage() string {
	return r.Message
}

type FileInfos struct {
	Message     string
	Filename    string
	Permissions string
	Size        int64
}

func (r *Response) ParseFileInfos() (*FileInfos, error) {
	message := strings.ReplaceAll(r.Message, "\n", "")
	parts := strings.Split(message, " ")
	if len(parts) < 3 {
		return nil, errors.New("unable to parse message as file infos")
	}

	size, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return &FileInfos{
		Message:     r.Message,
		Permissions: parts[0],
		Size:        int64(size),
		Filename:    parts[2],
	}, nil
}

// Ack writes an `Ack` message to the remote, does not await its response, a seperate call to ParseResponse is
// therefore required to check if the acknowledgement succeeded.
func Ack(writer io.Writer) error {
	var msg = []byte{0}
	n, err := writer.Write(msg)
	if err != nil {
		return err
	}
	if n < len(msg) {
		return errors.New("failed to write ack buffer")
	}
	return nil
}
