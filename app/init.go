package app

import (
	"errors"
	"fmt"
)

// Creates mvcs files and directories
func runInit() error {
	// Check if mvcs dir exists and return error if it does
	pathExists, err := checkPathExists(MVCS_DIR)

	if err != nil {
		return err
	}

	if pathExists {
		return errors.New(".mvcs directory already exists")
	}

	// Create mvcs files and directories
	if err := createDirs(MVCS_DIR, COMMIT_DIR, TREE_DIR, BLOB_DIR); err != nil {
		return err
	}

	if err := createFiles(HEAD_FILE, STAGE_FILE, CFG_FILE, LOGS_FILE); err != nil {
		return err
	}

	return nil
}
