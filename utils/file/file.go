package file

import (
	"io"
)

type File interface {
	Save(path string, buffer io.ReadCloser) (err error)
	IsExist(path string) (exist bool, err error)
	GetCurrentAbsPath() (path string, err error)
}
