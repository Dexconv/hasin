package common

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const MB = 1000000

func DirSize(dirpath string) (dirsize int64, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)

	err = os.Chdir(dirpath)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			dirsize += file.Size()
		}
	}
	err = os.Chdir(exPath)
	return
}
