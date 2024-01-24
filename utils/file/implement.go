package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func New() File {
	return &implement{}
}

type implement struct {
}

func (*implement) Save(path string, buffer io.ReadCloser) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, buffer)
	if err != nil {
		return err
	}
	return nil
}

func (*implement) IsExist(path string) (bool, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return true, nil
	}
	return true, nil
}

func (*implement) GetCurrentAbsPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return path, nil
}
