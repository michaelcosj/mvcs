package models

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"path/filepath"
	"strings"
)

type Tree struct {
	basename string
	path     string
	hash     string
	content  string
	treeList []*Tree
	blobList []*blob
}

func NewTree(path string) *Tree {
	return &Tree{
		basename: filepath.Base(path),
		path:     path,
	}
}

// tree should include its own path?
func getTreeFromHash(path, hash string) (*Tree, error) {
	file := filepath.Join(constants.OBJ_DIR, hash)
	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}

	data := parseTree(content)

	treeList := make([]*Tree, 0)
	for _, trData := range data.treeList {
		trPath := filepath.Join(path, trData["basename"])
		subTree, err := getTreeFromHash(trPath, trData["hash"])
		if err != nil {
			return nil, err
		}
		treeList = append(treeList, subTree)
	}

	blobList := make([]*blob, 0)
	for _, blData := range data.blobList {
		blPath := filepath.Join(path, blData["name"])
		bl, err := getBlobFromHash(blPath, blData["hash"])
		if err != nil {
			return nil, err
		}
		blobList = append(blobList, bl)
	}

	return &Tree{filepath.Base(path), path, hash, content, treeList, blobList}, nil
}

func (t Tree) Files() map[string]string {
	files := make(map[string]string)

	for _, tree := range t.treeList {
		for path, content := range tree.Files() {
			files[path] = content
		}
	}

	for _, blob := range t.blobList {
		files[blob.path] = blob.content
	}

	return files
}

func (t Tree) compressAndSave() error {
	for _, tr := range t.treeList {
		if err := tr.compressAndSave(); err != nil {
			return err
		}
	}

	for _, file := range t.blobList {
		if err := file.compressAndSave(); err != nil {
			return err
		}
	}

	dstPath := filepath.Join(constants.OBJ_DIR, t.hash)
	return helpers.CompressStrToFile(dstPath, t.content)
}

func (t *Tree) findTree(name string) (bool, *Tree) {
	for _, tr := range t.treeList {
		if tr.basename == name {
			return true, tr
		}
	}
	return false, nil
}

func (t *Tree) findBlob(name string) (bool, *blob) {
	for _, blob := range t.blobList {
		if blob.name == name {
			return true, blob
		}
	}
	return false, nil
}

func (t *Tree) AddTree(subTr *Tree) {
	if found, tr := t.findTree(subTr.basename); found {
		tr.blobList = append(tr.blobList, subTr.blobList...)
		return
	}

	t.treeList = append(t.treeList, subTr)
}

func (t *Tree) addChild(child *blob) {
	parentDir := filepath.Dir(child.path)

	if t.path != parentDir && parentDir != "." {
		subTree := NewTree(parentDir)
		subTree.addChild(child)

		dirs := strings.Split(parentDir, "/")
		for i := len(dirs) - 1; i > 0; i-- {
			tmpTree := NewTree(filepath.Join(dirs[:i]...))
			tmpTree.treeList = append(tmpTree.treeList, subTree)
			subTree = tmpTree
		}

		t.AddTree(subTree)
		return
	}

	if found, blob := t.findBlob(child.name); found {
		child.hash = blob.hash
	}
	t.blobList = append(t.blobList, child)
}

func (t *Tree) AddChildren(files []string) error {
	for _, file := range files {
		blob, err := NewBlob(file)
		if err != nil {
			return err
		}
		t.addChild(blob)
	}
	return nil
}

func (t *Tree) generateHash() error {
	var content strings.Builder

	for _, tr := range t.treeList {
		tr.generateHash()
		_, err := content.WriteString("tree: " + tr.hash + " " + tr.basename + "\n")
		if err != nil {
			return err
		}
	}

	for _, file := range t.blobList {
		_, err := content.WriteString("blob: " + file.hash + " " + file.name + "\n")
		if err != nil {
			return err
		}
	}

	t.content = content.String()
	t.hash = helpers.HashStr(t.content)

	return nil
}
