package goreq

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Save(name string, res *http.Response) error {
	suffix := filepath.Ext(name)
	if suffix == "" {
		suffix = GetResSuffix(res)
		name += suffix
	}
	reader, err := NewReaderFromRes(res)
	if err != nil {
		return err
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}
