package commands

import (
	"fmt"

	"github.com/michaelcosj/mvcs/app/models"
)

// staged / unstaged
// modified/created/deleted

func RunStatus() error {
	stage, err := models.GetStage()
	if err != nil {
		return err
	}

	if len(*stage) > 0 {
		for _, file := range *stage {
			fmt.Printf("Filepath: %s\tModified: %s\n", file.Path, file.ModTime)
		}
	} else {
		fmt.Println("Nothing added to the staging area")
	}

	return nil
}
