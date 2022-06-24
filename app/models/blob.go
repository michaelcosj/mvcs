package models

import (
	"fmt"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"path/filepath"
)

const hashLength = 32

type blob struct {
	name    string
	path    string
	content string
	hash    string
}

func NewBlob(file string) (*blob, error) {
	content, err := helpers.GetFileContent(file)
	if err != nil {
		return nil, err
	}
	return &blob{filepath.Base(file), file, content, helpers.HashStr(content)}, nil
}

func getBlobFromFile(path, file string) (*blob, error) {
	hash := filepath.Base(file)
	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)
  fmt.Println("here: ", name)

	return &blob{name, path, content, hash}, nil
}

func (fm blob) compressAndSave() error {
	dstPath := filepath.Join(constants.OBJ_DIR, fm.hash)
	return helpers.CompressStrToFile(dstPath, fm.content)
}