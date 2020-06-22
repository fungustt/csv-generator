package generator

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestDirGenerateWillThrowOnMkDirError(t *testing.T) {
	err := errors.New("mkdir mock error")

	mockOs := new(MockOs)
	mockOs.On("Mkdir", mock.Anything, mock.Anything).Return(err)

	dir := NewDir(1, mockOs)

	fileCh := make(chan *target)
	errCh := make(chan error)

	go dir.generate(
		&target{
			nesting: 1,
			path:    "/tmp",
		},
		fileCh,
		errCh,
	)

	assert.EqualError(t, <-errCh, err.Error())
}

func TestAddSubTarget(t *testing.T) {
	mockOs := new(MockOs)
	mockOs.On("Mkdir", mock.Anything, mock.Anything).Return(nil)

	dir := NewDir(1, mockOs)

	fileCh := make(chan *target)
	errCh := make(chan error)

	go dir.generate(
		&target{
			nesting: 1,
			path:    "/tmp",
		},
		fileCh,
		errCh,
	)

	actual := <-fileCh
	assert.Regexp(t, regexp.MustCompile("^/tmp/[^/]+$"), actual.path)
	assert.Equal(t, 0, actual.nesting)
}

func TestRecursiveCall(t *testing.T) {
	mockOs := new(MockOs)
	mockOs.On("Mkdir", mock.Anything, mock.Anything).Return(nil)

	dir := NewDir(1, mockOs)

	fileCh := make(chan *target)
	errCh := make(chan error)

	go dir.generate(
		&target{
			nesting: 2,
			path:    "/tmp",
		},
		fileCh,
		errCh,
	)

	actual := <-fileCh
	assert.Regexp(t, regexp.MustCompile("^/tmp/[^/]+$"), actual.path)
	assert.Equal(t, 1, actual.nesting)

	actual = <-fileCh
	assert.Regexp(t, regexp.MustCompile("^/tmp/[^/]+/[^/]+$"), actual.path)
	assert.Equal(t, 0, actual.nesting)
}
