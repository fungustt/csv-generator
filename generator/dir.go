package generator

import (
	"fmt"
	"github.com/fungustt/generator/random"
)

const dirLen = 5

type Dir struct {
	dirsInDir int
	os        Os
	dirRand   *random.StrRandomizer
}

func NewDir(dirsInDir int, os Os) *Dir {
	return &Dir{
		dirsInDir: dirsInDir,
		os:        os,
		dirRand:   random.NewStrRandomizer(dirLen),
	}
}

func (d *Dir) cleanup(rootDir string, errCh chan<- error) {
	if err := d.os.RemoveAll(rootDir); err != nil {
		errCh <- err
	}

	if err := d.os.Mkdir(rootDir, 0755); err != nil {
		errCh <- err
	}
}

func (d *Dir) generate(t *target, fileCh chan<- *target, errCh chan<- error) {
	for dirsToCreate := 1; dirsToCreate <= d.dirsInDir; dirsToCreate++ {
		// Generating path to create dir in given dir
		newPath := fmt.Sprintf("%s/%s", t.path, d.dirRand.Get())

		// Create new dir
		err := d.os.Mkdir(newPath, 0755)
		if err != nil {
			errCh <- err
			return
		}

		// Use this dir as target for subdirs
		subTarget := &target{
			nesting: t.nesting - 1,
			path:    newPath,
		}

		// Pass newly created dir target to channel to start file generation
		fileCh <- subTarget

		// If more nesting needed, recursively generate subdirs
		if t.nesting > 1 {
			d.generate(subTarget, fileCh, errCh)
		}
	}
}
