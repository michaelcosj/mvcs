package commands

import (
	"fmt"
	"michaelcosj/mvcs/app/models"
	"path/filepath"
)

// staged / unstaged
// modified/created/deleted


func RunStatus() error {
  stage, err := models.GetStage()
  if err != nil {
    return err
  }

  for _, file := range *stage {
    fmt.Printf("Filepath: %s\tModified: %s\n", filepath.Base(file.Path), file.ModTime)
  }
  return nil
}
