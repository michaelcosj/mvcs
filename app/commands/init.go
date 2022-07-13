package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaelcosj/mvcs/app/constants"
	"github.com/michaelcosj/mvcs/app/helpers"
)

func RunInit(configArgs []string) error {
	pathExists, err := helpers.CheckPathExists(constants.MVCS_DIR)
	if err != nil {
		return fmt.Errorf("error checking mvcs dir: %w", err)
	}

	if pathExists {
		return errors.New(".mvcs directory already exists")
	}

	username, email, err := parseConfigArgs(configArgs)
	if err != nil {
		return err
	}

	if err := helpers.CreateDirs(constants.MVCS_DIR, constants.OBJ_DIR); err != nil {
		return fmt.Errorf("error creating mvcs directories: %w", err)
	}

	if err := helpers.CreateFiles(constants.HEAD_FILE, constants.STAGE_FILE, constants.CFG_FILE, constants.LOGS_FILE); err != nil {
		return fmt.Errorf("error creating mvcs files: %w", err)
	}

	config := fmt.Sprintf("username=%s\nemail=%s\n", username, email)
	if err := helpers.WriteToFile(constants.CFG_FILE, config); err != nil {
		return err
	}

	return nil
}

func parseConfigArgs(args []string) (string, string, error) {
	var username, email string

	for _, field := range args {
		if strings.HasPrefix(field, "name=") {
			username = strings.TrimSpace(strings.TrimPrefix(field, "name="))
		} else if strings.HasPrefix(field, "email=") {
			email = strings.TrimSpace(strings.TrimPrefix(field, "email="))
		}
	}

	if !helpers.IsValidUsername(username) {
		return "", "", errors.New("invalid username")
	}

	if !helpers.IsValidEmail(email) {
		return "", "", errors.New("invalid email")
	}

	return username, email, nil
}
