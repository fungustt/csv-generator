package generator

import "os"

type Os interface {
	Create(name string) (*os.File, error)
	Mkdir(name string, perm os.FileMode) error
	RemoveAll(path string) error
}
