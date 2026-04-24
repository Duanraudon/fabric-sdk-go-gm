package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// DirExists ...
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// WriteFile ...
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	dirPath := filepath.Dir(filename)
	exists := DirExists(dirPath)
	if !exists {
		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, data, perm)
}

// TranslatePaths ...
func TranslatePaths(configDir string, p *string) string {
	if filepath.IsAbs(*p) {
		return *p
	}
	return filepath.Join(configDir, *p)
}
