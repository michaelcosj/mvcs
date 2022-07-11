package helpers

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
	"strings"
)

func CheckPathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking %s exists: %w", path, err)
	}
	return true, nil
}

func CheckPathIsDir(path string) (bool, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fInfo.IsDir(), nil
}

func CreateDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func CreateDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := CreateDir(dir); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}
	return nil
}

func GetFilesInDir(dir string) ([]string, error) {
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

func CreateFile(file string) error {
	_, err := os.Create(file)
	return err
}

func CreateFiles(files ...string) error {
	for _, file := range files {
		if err := CreateFile(file); err != nil {
			return fmt.Errorf("error creating file %s: %w", file, err)
		}
	}
	return nil
}

func WriteToFile(file, data string) error {
	return os.WriteFile(file, []byte(data), os.ModePerm)
}

func ClearFiles(paths ...string) error {
	for _, path := range paths {
		if err := CreateFile(path); err != nil {
			return err
		}
	}
	return nil
}

func GetFileModTime(file string, formatString string) (string, error) {
	fInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	return fInfo.ModTime().Format(formatString), nil
}

func GetFileContent(file string) (string, error) {
	data, err := os.ReadFile(file)
	return string(data), err
}

// hash string (md5)
func HashStr(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// compress string (zlib)
func CompressStr(data string) (string, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	w.Write([]byte(data))
	w.Close()

	return buf.String(), nil
}

// compress data and write it to dstfile
func CompressStrToFile(dstFile, data string) error {
	compressedData, err := CompressStr(data)
	if err != nil {
		return err
	}
	return WriteToFile(dstFile, compressedData)
}

func DecompressStr(data string) (string, error) {
	var content strings.Builder
	buf := bytes.NewBufferString(data)

	r, err := zlib.NewReader(buf)
	if err != nil {
		return "", nil
	}

	io.Copy(&content, r)
	r.Close()

	return content.String(), nil
}

func DecompressFile(file string) (string, error) {
  data, err := GetFileContent(file)
  if err != nil {
    return "", err
  }
  return DecompressStr(data)
}
