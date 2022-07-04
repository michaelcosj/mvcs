package commands

import (
	"fmt"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"michaelcosj/mvcs/app/models"
	"strings"
)


func RunHistory() error {
	head, err := helpers.GetFileContent(constants.HEAD_FILE)
	if err != nil {
		return err
	}
  headCommitHash := strings.TrimSpace(head)
  
  headCommit, err := models.GetCommitFromHash(headCommitHash) 
  if err != nil {
    return err
  }

  tmpCommit := headCommit
  fmt.Println(tmpCommit.Hash)
  for len(tmpCommit.ParentHash) == constants.HASH_LEN {
    tmpCommit, err = models.GetCommitFromHash(tmpCommit.ParentHash)
    if err != nil {
      return err
    }
    fmt.Println(tmpCommit.Hash)
  }

  return nil
}
