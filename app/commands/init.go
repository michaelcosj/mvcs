package commands

import (
	"errors"
	"fmt"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
)

func RunInit() error {
	pathExists, err := helpers.CheckPathExists(constants.MVCS_DIR)
	if err != nil {
		return fmt.Errorf("error checking mvcs dir: %w", err)
	}

	if pathExists {
		return errors.New(".mvcs directory already exists")
	}

	if err := helpers.CreateDirs(constants.MVCS_DIR, constants.OBJ_DIR); err != nil {
		return fmt.Errorf("error creating mvcs directories: %w", err)
	}

	if err := helpers.CreateFiles(constants.HEAD_FILE, constants.STAGE_FILE, constants.CFG_FILE, constants.LOGS_FILE); err != nil {
		return fmt.Errorf("error creating mvcs files: %w", err)
	}

	return nil
}
