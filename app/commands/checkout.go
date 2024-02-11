package commands

import (
	"github.com/michaelcosj/mvcs/app/constants"
	"github.com/michaelcosj/mvcs/app/helpers"
	"github.com/michaelcosj/mvcs/app/models"
)

func RunCheckout(commitHash string) error {
	commit, err := models.GetCommitFromHash(commitHash)
	if err != nil {
		return err
	}

	for path, content := range commit.RootTree.Files() {
		if err := helpers.CreateFile(path); err != nil {
			return err
		}
		if err := helpers.WriteToFile(path, content); err != nil {
			return err
		}
	}

	if err := helpers.ClearFiles(constants.STAGE_FILE, constants.HEAD_FILE); err != nil {
		return err
	}
	return helpers.WriteToFile(constants.HEAD_FILE, commit.Hash)
}
