package app

import (
	"path/filepath"
	"strings"
	"time"
)

type fileData struct {
	name    string
	path    string
	content string
	hash    string
}

func newFileData(path string) (*fileData, error) {
	content, err := getFileContent(path)
	if err != nil {
		return nil, err
	}
	return &fileData{filepath.Base(path), path, content, hashStr(content)}, nil
}

func (fm fileData) compressData(dstPath string) error {
	return compressStrToFile(filepath.Join(dstPath, fm.hash), fm.content)
}

type treeData struct {
	basename string
	path     string
	treeList []*treeData
	fileList []*fileData
	content  string
	hash     string
}

func newTreeData(path string) *treeData {
	return &treeData{
		basename: filepath.Base(path),
		path:     path,
	}
}

func (td treeData) compressData(dstPath string) error {
	for _, tree := range td.treeList {
		if err := tree.compressData(dstPath); err != nil {
			return err
		}
	}

	for _, file := range td.fileList {
		if err := file.compressData(dstPath); err != nil {
			return err
		}
	}

	return compressStrToFile(filepath.Join(dstPath, td.hash), td.content)
}

func (td *treeData) findTree(name string) (bool, *treeData) {
	for _, tree := range td.treeList {
		if tree.basename == name {
			return true, tree
		}
	}
	return false, nil
}

func (td *treeData) addChild(fd *fileData) {
	parentDir := filepath.Dir(fd.path)

	if td.path != parentDir && parentDir != "." {
		subTree := newTreeData(parentDir)
		subTree.addChild(fd)

		dirs := strings.Split(parentDir, "/")
		for i := len(dirs) - 1; i > 0; i-- {
			tmpTree := newTreeData(filepath.Join(dirs[:i]...))
			tmpTree.treeList = append(tmpTree.treeList, subTree)
			subTree = tmpTree
		}

		if found, tree := td.findTree(subTree.basename); found {
			tree.fileList = append(tree.fileList, subTree.fileList...)
			return
		}

		td.treeList = append(td.treeList, subTree)
		return
	}

	td.fileList = append(td.fileList, fd)
}

func (td *treeData) generateHash() error {
	var content strings.Builder

	for _, tree := range td.treeList {
		tree.generateHash()
		_, err := content.WriteString("tree: " + tree.hash + " " + filepath.Base(tree.basename) + "\n")
		if err != nil {
			return err
		}
	}

	for _, file := range td.fileList {
		_, err := content.WriteString("file: " + file.hash + " " + filepath.Base(file.name) + "\n")
		if err != nil {
			return err
		}
	}

	td.content = content.String()
	td.hash = hashStr(td.content)

	return nil
}

type commitData struct {
	message    string
	content    string
	tree       *treeData
	hash       string
	parentHash string
}

func newCommitData(parent, msg string, tree *treeData) *commitData {
	var content strings.Builder
	content.WriteString("Parent: " + parent + "\n")
	content.WriteString("Timestamp: " + time.Now().Format(TIME_FORMAT) + "\n")
	content.WriteString("Message: " + msg + "\n")
	content.WriteString("Tree: " + tree.hash + "\n")

	return &commitData{
		message: msg,
		content: content.String(),
		tree:    tree,
		hash:    hashStr(content.String()),
	}
}

func (cd commitData) compressData(dstPath string) error {
	if err := cd.tree.compressData(dstPath); err != nil {
		return err
	}
	return compressStrToFile(filepath.Join(dstPath, cd.hash), cd.content)
}
