package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func AbsolutePath(Datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(Datadir, filename)
}

func BackPwd() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("fail in back pwd: %v", err)
	}
	dir := filepath.Dir(path)
	return dir, nil
}
