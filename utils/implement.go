package utils

import (
	"math/rand"
	"os"

	"github.com/picoorg/common/utils/cert"
	"github.com/picoorg/common/utils/file"
	"github.com/picoorg/common/utils/http"
	"github.com/picoorg/common/utils/strings"
)

func New() Utils {
	return &implement{
		file:    file.New(),
		http:    http.New(),
		strings: strings.New(),
		cert:    cert.New(),
	}
}

type implement struct {
	file    file.File
	http    http.Http
	strings strings.String
	cert    cert.Cert
}

func (u *implement) File() file.File {
	return u.file
}

func (u *implement) Http() http.Http {
	return u.http
}

func (u *implement) Strings() strings.String {
	return u.strings
}

func (u *implement) Cert() cert.Cert {
	return u.cert
}

func (u *implement) Exit(err error) {
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func (u *implement) Randint(l, r int) int {
	return rand.Intn(r-l+1) + l
}

func (u *implement) RandString(collection string, length int) string {
	ret := ""
	for i := 0; i < length; i++ {
		index := u.Randint(0, len(collection)-1)
		ret += collection[index : index+1]
	}
	return ret
}
