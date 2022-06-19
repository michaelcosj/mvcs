package app

import (
	"fmt"
	"path/filepath"
	"sort"
)

func runAdd(paths ...string) error {
	stage, err := getStage(STAGE_FILE)
	if err != nil {
		return err
	}

	for _, path := range paths {
		exists, err := checkPathExists(path)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("path %s did not match any files", path)
		}

		isDir, err := checkPathIsDir(path)
		if err != nil {
			return err
		}

		if isDir {
			dirFiles, err := getFilesInDir(path)
			if err != nil {
				return err
			}

			for _, file := range dirFiles {
        parentDir := filepath.Dir(file)
				if parentDir != MVCS_DIR && parentDir != OBJ_DIR {
          if err := stage.addFile(file); err != nil {
            return err
          }
				}
			}
		} else {
			if err := stage.addFile(path); err != nil {
				return err
			}
		}

	}

	sort.Sort(stage)
	if err := stage.saveToFile(STAGE_FILE); err != nil {
		return err
	}

	return nil
}
