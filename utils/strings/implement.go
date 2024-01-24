package strings

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func New() String {
	return &implement{}
}

type implement struct {
}

func (m *implement) MD5(data io.ReadCloser) (string, error) {
	instance := md5.New()
	_, err := io.Copy(instance, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", instance.Sum(nil)), nil
}

func (m *implement) MD5String(data string) (string, error) {
	return m.MD5(io.NopCloser(bytes.NewBufferString(data)))
}

func (*implement) StructToJson(data any) (io.ReadCloser, error) {
	buf := bytes.NewBuffer([]byte{})
	return io.NopCloser(buf), json.NewEncoder(buf).Encode(data)
}

func (m *implement) AsciiLetters() string {
	return m.AsciiLowerLetters() + m.AsciiUpperLetters()
}

func (m *implement) AsciiNumbers() string {
	return "0123456789"
}

func (m *implement) AsciiLowerLetters() string {
	return "abcdefghijklmnopqrstuvwxyz"
}

func (m *implement) AsciiUpperLetters() string {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

func (m *implement) FString(template string, vars map[string]any) (string, error) {
	vars["$"] = "$"
	varStringMap := make(map[string]string)
	for k, v := range vars {
		varStringMap[k] = fmt.Sprintf("%v", v)
	}
	var err error
	ret := os.Expand(template, func(name string) string {
		value, ok := varStringMap[name]
		if !ok {
			err = fmt.Errorf("var %s not found", name)
			return ""
		}
		return value
	})
	if err != nil {
		return "", err
	}
	return ret, nil
}
