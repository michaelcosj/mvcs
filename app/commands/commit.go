package commands

import (
	"errors"
	"strings"

	"github.com/michaelcosj/mvcs/app/constants"
	"github.com/michaelcosj/mvcs/app/helpers"
	"github.com/michaelcosj/mvcs/app/models"
)

func RunCommit(msg string) error {
	parent := ""
	rootTree := models.NewTree(".")

	head, err := models.GetHeadCommit()
	if err != nil {
		return err
	}

	if head != nil {
		parent = head.Hash
		rootTree.AddFromTree(head.RootTree)
	}

	stg, err := models.GetStage()
	if err != nil {
		return err
	}

	if err := rootTree.AddChildren(stg.Files()); err != nil {
		return err
	}

	username, email, err := parseConfig()
	if err != nil {
		return err
	}

	commit, err := models.NewCommit(parent, msg, username, email, rootTree)
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

func parseConfig() (string, string, error) {
	var username, email string

	configData, err := helpers.GetFileContent(constants.CFG_FILE)
	if err != nil {
		return "", "", err
	}

	fields := strings.Split(configData, "\n")

	for _, field := range fields {
		if strings.HasPrefix(field, "username=") {
			username = strings.TrimSpace(strings.TrimPrefix(field, "username="))
		} else if strings.HasPrefix(field, "email=") {
			email = strings.TrimSpace(strings.TrimPrefix(field, "email="))
		}
	}

	if !helpers.IsValidUsername(username) {
		return "", "", errors.New("invalid username in config file")
	}

	if !helpers.IsValidEmail(email) {
		return "", "", errors.New("invalid email in config file")
	}

	return username, email, nil
}
