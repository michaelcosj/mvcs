package commands

import (
	"fmt"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"michaelcosj/mvcs/app/models"
	"path/filepath"
	"sort"
)

func RunAdd(paths ...string) error {
	stage, err := models.GetStage()
	if err != nil {
		return err
	}

	for _, path := range paths {
		exists, err := helpers.CheckPathExists(path)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("path %s did not match any files", path)
		}

		isDir, err := helpers.CheckPathIsDir(path)
		if err != nil {
			return err
		}

		if isDir {
			dirFiles, err := helpers.GetFilesInDir(path)
			if err != nil {
				return err
			}

			for _, file := range dirFiles {
        parentDir := filepath.Dir(file)
				if parentDir != constants.MVCS_DIR && parentDir != constants.OBJ_DIR {
          if err := stage.AddFile(file); err != nil {
            return err
          }
				}
			}
		} else {
			if err := stage.AddFile(path); err != nil {
				return err
			}
		}

	}

	sort.Sort(stage)
	return stage.Save()
}
