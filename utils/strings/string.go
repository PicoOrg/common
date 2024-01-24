package strings

import (
	"io"
)

type String interface {
	MD5(data io.ReadCloser) (md5 string, err error)
	MD5String(data string) (md5 string, err error)
	StructToJson(data any) (buffer io.ReadCloser, err error)
	AsciiLetters() string
	AsciiNumbers() string
	AsciiLowerLetters() string
	AsciiUpperLetters() string
	FString(template string, vars map[string]any) (ret string, err error)
}
