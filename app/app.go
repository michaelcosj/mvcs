package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaelcosj/mvcs/app/commands"
	"github.com/michaelcosj/mvcs/app/constants"
	"github.com/michaelcosj/mvcs/app/helpers"
)

func Run() error {
	program := os.Args[0]
	if len(os.Args) < 2 {
		commands.RunHelp(program)
		return errors.New("Not enough arguments")
	}

	command := os.Args[1]

	// Commands without arguments
	if len(os.Args) < 3 {
		switch command {
		case "help":
			commands.RunHelp(program)
			return nil
		case "status":
			return commands.RunStatus()
		case "history":
			return commands.RunHistory()
		default:
			return fmt.Errorf("not enough arguments for mvcs %s", command)
		}
	}

	// Commands with arguments
	switch command {
	case "init":
		args := os.Args[2:]
		return commands.RunInit(args)
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
		fmt.Print(data)
	default:
		commands.RunHelp(program)
		return fmt.Errorf("Invalid command '%s'", command)
	}

	return nil
}
