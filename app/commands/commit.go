package commands

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"michaelcosj/mvcs/app/models"
)

func RunCommit(msg string) error {
  parentHash := ""
  rootTree := models.NewTree(".")

  head, err := models.GetHeadCommit()
  if err != nil {
    return err
  }

	if head != nil {
    parentHash = head.Hash
    rootTree.AddTree(head.RootTree)
	}

	stg, err := models.GetStage()
	if err != nil {
		return err
	}

  if err := rootTree.AddChildren(stg.Files()); err != nil {
    return err
  }

  commit, err := models.NewCommit(parentHash, msg, rootTree)
  if err != nil {
    return err
  }

  if err := commit.CompressAndSave(); err != nil {
    return err
  }

	if err := helpers.ClearFiles(constants.STAGE_FILE, constants.HEAD_FILE); err != nil {
		return err
	}

	return helpers.WriteToFile(constants.HEAD_FILE, commit.Hash)
}
