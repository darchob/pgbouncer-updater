package configuration

import "fmt"

var FileNotFound error = fmt.Errorf("file not found")

func FileNotFoundFunc() error {
	return FileNotFound
}
