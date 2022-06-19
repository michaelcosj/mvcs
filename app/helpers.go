package app

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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

func getFilesInDir(dir string) ([]string, error) {
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

func openFile(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
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

func writeToFile(path, data string) error {
	return os.WriteFile(path, []byte(data), os.ModePerm)
}

func clearFile(path string) error {
	fp, err := openFile(path)
	if err != nil {
		return err
	}

	if err := fp.Truncate(0); err != nil {
		return err
	}
	_, err = fp.Seek(0, io.SeekStart)
	return err
}

func clearFiles(paths... string) error {
  for _, path := range paths {
    if err := clearFile(path); err != nil {
      return err
    }
  }
  return nil
}

func getFileModTime(file string, formatString string) (string, error) {
	fInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	return fInfo.ModTime().Format(formatString), nil
}

func getFileContent(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	return string(data), err
}

// hash string (md5)
func hashStr(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// compress string (zlib)
func compressStr(data string) (string, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	w.Write([]byte(data))
	w.Close()

	return buf.String(), nil
}

// compress data and write it to dstfile
func compressStrToFile(dstFile, data string) error {
	compressedData, err := compressStr(data)
	if err != nil {
		return err
	}
	return writeToFile(dstFile, compressedData)
}
