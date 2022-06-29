package models

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"path/filepath"
	"strings"
)

type Commit struct {
	Hash     string
	RootTree *Tree

	message    string
	content    string
	parentHash string
}

func NewCommit(parentHash, msg string) (*Commit, error) {
  tree := NewTree(".")

  if len(parentHash) == 32 {
    parentCommit, err := getCommitFromHash(parentHash)
    if err != nil {
      return nil, err
    }
    tree = parentCommit.RootTree
  }

	commit := &Commit{
		content:    "",
		Hash:       "",
		message:    msg,
		RootTree:   tree,
		parentHash: parentHash,
	}
	return commit, nil
}

func (cd *Commit) GenerateHash() error {
	var content strings.Builder

	if len(cd.parentHash) == hashLength {
		content.WriteString("parent: " + cd.parentHash + "\n")
	}

	if err := cd.RootTree.generateHash(); err != nil {
		return err
	}

	content.WriteString("tree: " + cd.RootTree.hash + "\n")
	content.WriteString("message: " + cd.message + "\n")

	cd.content = content.String()
	cd.Hash = helpers.HashStr(cd.content)
	return nil
}

func getCommitFromHash(hash string) (*Commit, error) {
	file := filepath.Join(constants.OBJ_DIR, strings.TrimSpace(hash))

	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}

	data := parseCommit(content)

	treeFile := filepath.Join(constants.OBJ_DIR, data.treeHash)
	tree, err := getTreeFromFile(".", treeFile)
	if err != nil {
		return nil, err
	}

	commit := &Commit{
		Hash:       hash,
		message:    data.message,
		parentHash: data.parentHash,
		content:    content,
		RootTree:   tree,
	}

	return commit, nil
}

func (cd Commit) CompressAndSave() error {
	if err := cd.RootTree.compressAndSave(); err != nil {
		return err
	}
	dstPath := filepath.Join(constants.OBJ_DIR, cd.Hash)
	return helpers.CompressStrToFile(dstPath, cd.content)
}
