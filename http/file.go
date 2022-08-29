package http

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func ReceiveFileToLocal(formFileKey string, fileDest string, r *http.Request) error {
	file, _, err := r.FormFile(formFileKey)
	if err != nil {
		return err
	}
	defer file.Close()

	dst, err := os.Create(fileDest)
	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func ReceiveFileToBytes(formFileKey string, r *http.Request) (fileByte []byte, err error) {
	file, _, err := r.FormFile(formFileKey)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)

	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	fileByte = buf.Bytes()

	return
}
