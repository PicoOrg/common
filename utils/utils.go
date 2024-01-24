package utils

import (
	"github.com/picoorg/common/utils/cert"
	"github.com/picoorg/common/utils/file"
	"github.com/picoorg/common/utils/http"
	"github.com/picoorg/common/utils/strings"
)

type Utils interface {
	File() file.File
	Http() http.Http
	Strings() strings.String
	Cert() cert.Cert
	Exit(err error)
	Randint(l, r int) (random int) // [l, r]
	RandString(collection string, length int) (random string)
}
