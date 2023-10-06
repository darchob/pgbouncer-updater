package sendfile

import (
	"crypto/md5"
	"io"
)

func copyN(writer io.Writer, src io.Reader, size int64) (int64, error) {
	var total int64
	total = 0
	for total < size {
		n, err := io.CopyN(writer, src, size)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}

func md5sum(file io.Reader) ([]byte, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
