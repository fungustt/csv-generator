package generator

import (
	"github.com/fungustt/generator/config"
	"sync"
)

type Generator struct {
	csvDir  string
	dir     *Dir
	file    *File
	nesting int
}

type target struct {
	nesting int
	path    string
}

func (g *Generator) Generate() error {
	fileCh := make(chan *target)
	closeCh := make(chan int)

	errCh := make(chan error)
	wgChan := make(chan bool)

	wg := sync.WaitGroup{}
	initialTarget := &target{
		path:    g.csvDir,
		nesting: g.nesting,
	}

	wg.Add(2)
	go g.file.listen(fileCh, closeCh, errCh, &wg)

	go func() {
		g.dir.cleanup(initialTarget.path, errCh)
		fileCh <- initialTarget
		g.dir.generate(initialTarget, fileCh, errCh)
		wg.Done()

		closeCh <- 1
	}()

	go func() {
		wg.Wait()
		close(wgChan)
	}()

	for {
		select {
		case <-wgChan:
			return nil
		case err := <-errCh:
			close(errCh)
			return err
		}
	}
}

func NewGenerator(configuration *config.Configuration, os Os) (*Generator, error) {
	g := &Generator{
		csvDir: configuration.CsvDir,
		file: NewFile(
			os,
			configuration.FilesInDir,
			configuration.StringsInFile,
			configuration.Measurements,
		),
		dir:     NewDir(configuration.DirsInDir, os),
		nesting: configuration.DirNesting,
	}

	return g, nil
}
