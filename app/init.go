package app

import (
	"errors"
	"fmt"
)

func runInit() error {
	pathExists, err := checkPathExists(MVCS_DIR)
	if err != nil {
		return fmt.Errorf("error checking mvcs dir: %w", err)
	}

	if pathExists {
		return errors.New(".mvcs directory already exists")
	}

	if err := createDirs(MVCS_DIR, OBJ_DIR); err != nil {
		return fmt.Errorf("error creating mvcs directories: %w", err)
	}

	if err := createFiles(HEAD_FILE, STAGE_FILE, CFG_FILE, LOGS_FILE); err != nil {
		return fmt.Errorf("error creating mvcs files: %w", err)
	}

	return nil
}
