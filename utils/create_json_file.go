package utils

import (
	"bytes"
	"io"
	"os"
)

func CreateJSONFile(fileName string, data []byte) error {
	buff := bytes.NewBuffer(data)

	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	io.Copy(file, buff)

	return nil
}
