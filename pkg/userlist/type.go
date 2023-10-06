package userlist

import (
	"fmt"
	"io"
	"reflect"
)

type User struct {
	file     io.Writer
	UserName string
	Md5      string
}

const (
	formatedList = "\"%s\" \"%s\"\n"
)

func (u *User) WriteMany(users interface{}) error {
	switch list := users.(type) {
	case map[string]string:
		for k, v := range list {
			user := &User{
				file:     u.file,
				UserName: k,
				Md5:      v,
			}

			if err := user.write(); err != nil {
				return err
			}
		}

	case []*User:
		for _, user := range list {
			user.file = u.file
			if err := user.write(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("users list Type %v not implemented yet", reflect.TypeOf(users).String())
	}
	return nil
}

func (u *User) write() error {
	_, err := fmt.Fprintf(u.file, formatedList, u.UserName, u.Md5)
	if err != nil {
		return err
	}
	return nil
}
