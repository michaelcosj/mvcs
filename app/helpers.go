package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func checkPathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking %s exists: %w", path, err)
	}

	return true, nil
}

func checkPathIsDir(path string) (bool, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fInfo.IsDir(), nil
}

func getFileModTime(file string, formatString string) (string, error) {
	fInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	return fInfo.ModTime().Format(formatString), nil
}

func createDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func createDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDir(dir); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	return nil
}

func openFile(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
}

func openFileReadOnly(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_RDONLY, os.ModePerm)
}

func openFileWriteOnly(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

func createFile(file string) error {
	_, err := openFile(file)
	return err
}

func createFiles(files ...string) error {
	for _, file := range files {
		if err := createFile(file); err != nil {
			return fmt.Errorf("error creating file %s: %w", file, err)
		}
	}

	return nil
}

func getAllFilesInDir(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
		  files = append(files, path)
		}

		return nil
	})

	return files, err
}

// hash string (md5)

// compress string (gzip)
