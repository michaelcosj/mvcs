package commands

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"michaelcosj/mvcs/app/models"
	"strings"
)

func RunCommit(msg string) error {
	head, err := helpers.GetFileContent(constants.HEAD_FILE)
	if err != nil {
		return err
	}
  parentHash := strings.TrimSpace(head)

  commit, err := models.NewCommit(parentHash, msg)
  if err != nil {
    return err
  }

	stg, err := models.GetStage()
	if err != nil {
		return err
	}
  if err := commit.RootTree.AddChildren(stg.Files()); err != nil {
    return err
  }

  if err := commit.GenerateHash(); err != nil {
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
