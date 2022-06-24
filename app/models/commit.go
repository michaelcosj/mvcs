package models

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"os"
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

func NewCommit(parent, msg string, tree *Tree) (*Commit, error) {
	var content strings.Builder

	if len(parent) == hashLength {
		content.WriteString("parent: " + parent + "\n")
	}

	if err := tree.GenerateHash(); err != nil {
		return nil, err
	}

	content.WriteString("tree: " + tree.hash + "\n")
	content.WriteString("message: " + msg + "\n")

	commit := &Commit{
		content:    content.String(),
		Hash:       helpers.HashStr(content.String()),
		message:    msg,
		RootTree:   tree,
		parentHash: parent,
	}

	return commit, nil
}

func GetCommitFromHash(hash string) (*Commit, error) {
	file := filepath.Join(constants.OBJ_DIR, strings.TrimSpace(hash))

	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}

	data := parseCommit(content)
	tree, err := getTreeFromHash(data.treeHash)
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

func getTreeFromHash(hash string) (*Tree, error) {
	treeFile := filepath.Join(constants.OBJ_DIR, hash)

	repoPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return getTreeFromFile(repoPath, treeFile)
}

func (cd Commit) CompressAndSave() error {
	if err := cd.RootTree.compressAndSave(); err != nil {
		return err
	}
	dstPath := filepath.Join(constants.OBJ_DIR, cd.Hash)
	return helpers.CompressStrToFile(dstPath, cd.content)
}
