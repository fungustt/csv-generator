package generator

import (
	"encoding/csv"
	"fmt"
	"github.com/fungustt/generator/random"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	DefaultBatchSize = 10000
	headerFieldLen   = 3
	fileLen          = 7
)

var (
	minTime = time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxTime = time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
)

type File struct {
	filesInDir    int
	measurements  int
	stringsInFile int
	batchSize     int
	os            Os
	fileRand      *random.StrRandomizer

	headers []string
}

func NewFile(os Os, filesInDir, stringsInFile, measurements int) *File {
	headerGen := random.NewStrRandomizer(headerFieldLen)
	defer headerGen.Stop()

	// Generating headers
	headers := []string{"date"}
	for i := 1; i <= measurements; i++ {
		headers = append(headers, headerGen.Get())
	}

	// Use strings count as batch size if batch size is less that first one
	batchSize := DefaultBatchSize
	if stringsInFile < DefaultBatchSize {
		batchSize = stringsInFile
	}

	return &File{
		os:            os,
		filesInDir:    filesInDir,
		stringsInFile: stringsInFile,
		measurements:  measurements,
		batchSize:     batchSize,
		headers:       headers,
		fileRand:      random.NewStrRandomizer(fileLen),
	}
}

func (f *File) listen(fileCh <-chan *target, closeCh <-chan int, errCh chan<- error, wg *sync.WaitGroup) {
	fileWg := sync.WaitGroup{}
	for {
		select {
		case t := <-fileCh:
			fileWg.Add(f.filesInDir)
			// TODO add goroutine pool to prevent process struggle
			for i := 1; i <= f.filesInDir; i++ {
				go func() {
					f.generate(t, errCh)
					fileWg.Done()
				}()
			}
		case <-closeCh:
			fileWg.Wait()
			f.fileRand.Stop()
			wg.Done()
			break
		}
	}
}

func (f *File) generate(target *target, errCh chan<- error) {
	path := fmt.Sprintf("%s/%s.csv", target.path, f.fileRand.Get())
	file, err := os.Create(path)
	if err != nil {
		errCh <- err
		return
	}

	// Close file when function exit
	defer func(f *os.File) {
		if err := file.Close(); err != nil {
			errCh <- err
			return
		}
	}(file)

	// New csv file
	writer := csv.NewWriter(file)
	// Flush data to file on function exit
	defer writer.Flush()

	// Write headers to file
	if err := writer.Write(f.headers); err != nil {
		errCh <- err
		return
	}

	// Data object size
	var size int

	// Rand is not thread-safe, to prevent mutex usage, which slows runtime we will use new rand for every file
	randObj := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate write object
	allData := make([][]string, f.batchSize)

	batches := int(math.Ceil(float64(f.stringsInFile) / float64(f.batchSize)))
	for batch := 1; batch <= batches; batch++ {
		// Calculate size of batch object
		size = f.batchSize
		if batch == batches && batch != 1 {
			size = f.stringsInFile - f.batchSize*(batch-1)
		}

		// Reuse memory
		allData = allData[:size]

		// Collect random data
		for str := 0; str < size; str++ {
			// Add random date
			data := []string{time.Unix(randObj.Int63n(maxTime-minTime)+minTime, 0).Format("2006-01-02")}

			// Add float val for all measurements
			for measurement := 1; measurement <= f.measurements; measurement++ {
				data = append(data, fmt.Sprintf("%f", randObj.NormFloat64()))
			}

			allData[str] = data
		}

		// Write data to file
		if err := writer.WriteAll(allData); err != nil {
			errCh <- err
			return
		}

		allData = allData[:0]
	}
}
