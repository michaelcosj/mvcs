package app

import (
	"errors"
	"fmt"
	"michaelcosj/mvcs/app/commands"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"os"
	"path/filepath"
	"strings"
)

func Run() error {
	if len(os.Args) < 2 {
		return errors.New("Not enough arguments")
	}

	program := os.Args[0]
	command := os.Args[1]

  var err error = nil

	switch command {
	case "help":
		commands.RunHelp(program)
	case "init":
		err = commands.RunInit()
	case "status":
    err = commands.RunStatus()
	case "history":
    err = commands.RunHistory()
	case "add":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs add")
		}
		paths := os.Args[2:]
		err = commands.RunAdd(paths...)
	case "commit":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs commit")
		}
		msg := os.Args[2]
		err = commands.RunCommit(msg)
	case "checkout":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs commit")
		}
		commitHash := strings.TrimSpace(os.Args[2])
    err = commands.RunCheckout(commitHash)
	case "read-hash":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs cat-file")
		}
		hash := os.Args[2]
    data, err := helpers.DecompressFile(filepath.Join(constants.OBJ_DIR, hash))
    if err != nil {
      return err
    }
    fmt.Println(data)
	default:
		commands.RunHelp(program)
		return fmt.Errorf("Invalid command '%s'", command)
	}
  
  if err != nil {
		return fmt.Errorf("mvcs %s failed: %s", command, err.Error())
  }
  fmt.Printf("mvcs %s done\n", command)

	return nil
}
