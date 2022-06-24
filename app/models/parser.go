package models

import (
	"strings"
)

type commitData struct {
	treeHash   string
	parentHash string
	message    string
}

type treeData struct {
	treeList []map[string]string
	blobList []map[string]string
}

func parseCommit(content string) commitData {
	var treeHash, parentHash, msg string

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "tree") {
			treeHash = strings.TrimSpace(strings.TrimPrefix(line, "tree: "))
		} else if strings.HasPrefix(line, "message") {
			msg = strings.TrimSpace(strings.TrimPrefix(line, "msg: "))
		} else if strings.HasPrefix(line, "parent") {
			parentHash = strings.TrimSpace(strings.TrimPrefix(line, "parent: "))
		}
	}

	return commitData{treeHash, parentHash, msg}
}

func parseTree(content string) treeData {
	var treeList, blobList []map[string]string

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "tree") {
			treeData := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "tree: ")), " ")

			treeHash := treeData[0]
			treeName := treeData[1]

			treeList = append(treeList, map[string]string{"hash": treeHash, "basename": treeName})
		} else if strings.HasPrefix(line, "blob") {
			blobData := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "blob: ")), " ")

			blobHash := blobData[0]
			blobName := blobData[1]

			blobList = append(blobList, map[string]string{"hash": blobHash, "name": blobName})
		}
	}

	return treeData{treeList, blobList}
}
