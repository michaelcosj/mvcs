package models

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/michaelcosj/mvcs/app/constants"
	"github.com/michaelcosj/mvcs/app/helpers"
)

type Commit struct {
	Hash      string
	Content   string
	Author    string
	Message   string
	Timestamp string
	RootTree  *Tree

	parentHash string
}

func NewCommit(parentHash, msg, username, email string, tree *Tree) (*Commit, error) {
	author := fmt.Sprintf("%s <%s>", username, email)

	commit := &Commit{
		Author:     author,
		Message:    msg,
		RootTree:   tree,
		parentHash: parentHash,
		Timestamp:  time.Now().Format(constants.TIME_FORMAT),
	}

	return commit, nil
}

func GetCommitFromHash(hash string) (*Commit, error) {
	file := filepath.Join(constants.OBJ_DIR, hash)

	content, err := helpers.DecompressFile(file)
	if err != nil {
		return nil, err
	}

	data := parseCommit(content)

	tree, err := getTreeFromHash(".", data.treeHash)
	if err != nil {
		return nil, err
	}

	commit := &Commit{
		Hash:       hash,
		Message:    data.message,
		Author:     data.author,
		parentHash: data.parentHash,
		Content:    content,
		RootTree:   tree,
	}

	return commit, nil
}

func GetHeadCommit() (*Commit, error) {
	headHash, err := helpers.GetFileContent(constants.HEAD_FILE)
	if err != nil {
		return nil, err
	}

	headHash = strings.TrimSpace(headHash)
	if len(headHash) != constants.HASH_LEN {
		return nil, nil
	}

	return GetCommitFromHash(headHash)
}

func (cm *Commit) generateHash() error {
	var content strings.Builder

	if len(cm.parentHash) == hashLength {
		content.WriteString("parent: " + cm.parentHash + "\n")
	}

	if err := cm.RootTree.generateHash(); err != nil {
		return err
	}

	content.WriteString("tree: " + cm.RootTree.hash + "\n")
	content.WriteString("author: " + cm.Author + "\n")
	content.WriteString("timestamp: " + cm.Timestamp + "\n")
	content.WriteString("message: " + cm.Message + "\n")

	cm.Content = content.String()
	cm.Hash = helpers.HashStr(cm.Content)

	return nil
}

func (cm *Commit) CompressAndSave() error {
	if err := cm.generateHash(); err != nil {
		return err
	}

	if err := cm.RootTree.compressAndSave(); err != nil {
		return err
	}

	dstPath := filepath.Join(constants.OBJ_DIR, cm.Hash)
	return helpers.CompressStrToFile(dstPath, cm.Content)
}

func (cm Commit) GetParent() (*Commit, error) {
	if len(cm.parentHash) != constants.HASH_LEN {
		return nil, nil
	}

	return GetCommitFromHash(cm.parentHash)
}
