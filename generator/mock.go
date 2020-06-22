package generator

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type MockOs struct {
	mock.Mock
}

func (m MockOs) Create(name string) (*os.File, error) {
	args := m.Called(name)

	return &os.File{}, args.Error(1)
}

func (m MockOs) Mkdir(name string, perm os.FileMode) error {
	args := m.Called(name, perm)

	return args.Error(0)
}

func (m MockOs) RemoveAll(path string) error {
	args := m.Called(path)

	return args.Error(0)
}
