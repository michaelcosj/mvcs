package commands

import (
	"fmt"

	"github.com/michaelcosj/mvcs/app/models"
)

func RunHistory() error {
	head, err := models.GetHeadCommit()
	if err != nil {
		return err
	}

	if head == nil {
		fmt.Println("This is emptiness")
		return nil
	}
	printCommit(*head, true)

	parentCommit, err := head.GetParent()
	if err != nil {
		return err
	}

	for parentCommit != nil {
		printCommit(*parentCommit, false)

		parentCommit, err = parentCommit.GetParent()
		if err != nil {
			return err
		}
	}

	return nil
}

func printCommit(c models.Commit, isHead bool) {
	headSuffix := ""
	if isHead {
		headSuffix = "(HEAD)"
	}

	fmt.Printf("commit %s %s\n", c.Hash, headSuffix)
	fmt.Printf("timestamp: %s\n", c.Timestamp)
	fmt.Printf("\t%s\n", c.Message)
}
