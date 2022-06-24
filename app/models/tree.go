package models

import (
	"fmt"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"path/filepath"
	"strings"
)

type Tree struct {
	basename string
	path     string
	treeList []*Tree
	blobList []*blob
	content  string
	hash     string
}

func NewTree(path string) *Tree {
	return &Tree{
		basename: filepath.Base(path),
		path:     path,
	}
}

// tree should include its own path?
func getTreeFromFile(path, file string) (*Tree, error) {
	hash := filepath.Base(file)
	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}

	data := parseTree(content)

	treeList := make([]*Tree, 0)
	blobList := make([]*blob, 0)

	for _, trData := range data.treeList {
		trFile := filepath.Join(constants.OBJ_DIR, trData["hash"])
		trPath := filepath.Join(path, trData["basename"])

		subTree, err := getTreeFromFile(trPath, trFile)
		if err != nil {
			return nil, err
		}

		treeList = append(treeList, subTree)
	}

	for _, blData := range data.blobList {
		blPath := filepath.Join(path, blData["name"])
		blFile := filepath.Join(constants.OBJ_DIR, blData["hash"])

		bl, err := getBlobFromFile(blPath, blFile)
		if err != nil {
			return nil, err
		}

		blobList = append(blobList, bl)
	}

  fmt.Println("Here: ", len(treeList), len(blobList))

	return &Tree{
		basename: filepath.Base(path),
		path:     path,
		treeList: treeList,
		blobList: blobList,
		hash:     hash,
		content:  content,
	}, nil
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

func (t *Tree) addBlob(bl *blob) {
	if found, blob := t.findBlob(bl.name); found {
		bl.hash = blob.hash
	}

	t.blobList = append(t.blobList, bl)
}

func (t *Tree) addTree(subTr *Tree) {
	if found, tr := t.findTree(subTr.basename); found {
		tr.blobList = append(tr.blobList, subTr.blobList...)
		return
	}

	t.treeList = append(t.treeList, subTr)
}

func (t *Tree) AddChild(child *blob) {
	parentDir := filepath.Dir(child.path)

	if t.path != parentDir && parentDir != "." {
		subTree := NewTree(parentDir)
		subTree.AddChild(child)

		dirs := strings.Split(parentDir, "/")
		for i := len(dirs) - 1; i > 0; i-- {
			tmpTree := NewTree(filepath.Join(dirs[:i]...))
			tmpTree.treeList = append(tmpTree.treeList, subTree)
			subTree = tmpTree
		}

		t.addTree(subTree)
		return
	}

	t.addBlob(child)
}

func (t *Tree) AddChildren(files []string) error {
  for _, file := range files {
    blob, err := NewBlob(file)
    if err != nil {
      return err
    }
    t.AddChild(blob)
  }
  return nil
}

func (t *Tree) GenerateHash() error {
	var content strings.Builder

	for _, tr := range t.treeList {
		tr.GenerateHash()
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
