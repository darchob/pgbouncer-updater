package userlist

import (
	"io"
	"os"
)

type UsersList interface {
	WriteMany(users interface{}) error
}

func NewUserList(file io.Writer) UsersList {
	return &User{
		file: file,
	}
}

func NewUserListToFile(filePath string) (UsersList, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return NewUserList(f), nil
}
