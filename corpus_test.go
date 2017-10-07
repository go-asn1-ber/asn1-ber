package ber

import (
	"os"
	"path/filepath"
	"testing"
)

//TestCorpus runs the data in the fuzzing corpus to make sure it doesn't start with any
//data that causes crashes
func TestCorpus(t *testing.T) {

	corpusFiles := []string{}

	filepath.Walk("./corpus", func(path string, f os.FileInfo, err error) error {
		corpusFiles = append(corpusFiles, path)
		return nil
	})

	for _, corpusFile := range corpusFiles {
		done := make(chan bool)
		go func() {
			f := corpusFile
			defer func() {
				close(done)
				if r := recover(); r != nil {
					t.Errorf("panic with corpus %s", f)
				}
			}()
			file, err := os.Open(f)
			if err != nil {
				panic(err)
			}

			_, _, _ = readPacket(file)
			_ = file.Close()
		}()
		<-done
	}
}
