package app

import (
	"fmt"
	"sort"
)

func runAdd(paths ...string) error {
  stage, err := getStageFromFile(STAGE_FILE);
  if err != nil {
    return err
  }

	for _, path := range paths {
		// Skip mvcs directory
		if path == MVCS_DIR {
			continue
		}

		// Check if path exists return error if it doesn't
		exists, err := checkPathExists(path)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("path %s did not match any files", path)
		}

		// Get all files recursively and save to stage
		isDir, err := checkPathIsDir(path)
		if err != nil {
			return err
		}

		if isDir {
			dirFiles, err := getAllFilesInDir(path)
			if err != nil {
				return err
			}

			for _, file := range dirFiles {
				if err := stage.addFile(file); err != nil {
					return err
				}
			}
		} else {
			if err := stage.addFile(path); err != nil {
				return err
			}
		}

	}

  // Sort stage and save stage to file
	sort.Sort(stage)
	if err := stage.saveToFile(STAGE_FILE); err != nil {
		return err
	}

	return nil
}
