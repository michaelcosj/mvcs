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
	program := os.Args[0]
	if len(os.Args) < 2 {
		commands.RunHelp(program)
		return errors.New("Not enough arguments")
	}
	command := os.Args[1]

	if len(os.Args) < 3 {
		switch command {
		case "help":
			commands.RunHelp(program)
		case "init":
			return commands.RunInit()
		case "status":
			return commands.RunStatus()
		case "history":
			return commands.RunHistory()
		default:
			return fmt.Errorf("not enough arguments for mvcs %s", command)
		}
	}

	switch command {
	case "add":
		paths := os.Args[2:]
		return commands.RunAdd(paths...)
	case "commit":
		msg := os.Args[2]
		return commands.RunCommit(msg)
	case "checkout":
		commitHash := strings.TrimSpace(os.Args[2])
		return commands.RunCheckout(commitHash)
	case "read-hash":
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

	fmt.Printf("mvcs %s done\n", command)
	return nil
}
