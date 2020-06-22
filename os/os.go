package os

import "os"

type Wrapper struct{}

func (w Wrapper) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (w Wrapper) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (w Wrapper) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func NewWrapper() Wrapper {
	return Wrapper{}
}
