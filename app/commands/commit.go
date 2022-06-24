package commands

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"michaelcosj/mvcs/app/models"
	"strings"
)

func RunCommit(msg string) error {
  // get commit at HEAD_FILE
	head, err := helpers.GetFileContent(constants.HEAD_FILE)
	if err != nil {
		return err
	}
  parentHash := strings.TrimSpace(head)
  // if commit exists, get tree from commit, else create a new tree
  tree := models.NewTree(".")
  if len(parentHash) == 32 {
    parentCommit, err := models.GetCommitFromHash(parentHash)
    if err != nil {
      return err
    }
    tree = parentCommit.RootTree
  }

  // add files to the tree
	stg, err := models.GetStage()
	if err != nil {
		return err
	}
  if err := tree.AddChildren(stg.Files()); err != nil {
    return err
  }

  // create a new commit and add the tree to it
  commit, err := models.NewCommit(parentHash, msg, tree)
  if err != nil {
    return err
  }

  // save commit
	commit.CompressAndSave()

  // clear files
	if err := helpers.ClearFiles(constants.STAGE_FILE, constants.HEAD_FILE); err != nil {
		return err
	}

  // add commit to HEAD
	return helpers.WriteToFile(constants.HEAD_FILE, commit.Hash)
}
