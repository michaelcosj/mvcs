package app

import (
	"errors"
	"fmt"
	"michaelcosj/mvcs/app/commands"
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"os"
	"path/filepath"
)

// Run help command
func runHelp(program string) {
	fmt.Println("Usage:", program, "[command]")
	fmt.Println(
		"MVCS - Mediocre Version Control System\n\n" +
			"    help                   display this help and exit \n" +
			"    init                   initialze mvcs in a repo \n" +
			"    status                 show the status of working tree \n" +
			"    add <file|directory>   add a file or directory to staging area \n" +
			"    commit <message>       commit all files or directories in the staging area \n" +
			"    revert <commit hash>   revert back to the specified commit \n" +
			"    history                show all commits and their info \n" +
			"    clean                  clean all orphaned commits",
	)
}

func Run() error {
	if len(os.Args) < 2 {
		return errors.New("Not enough arguments")
	}

	program := os.Args[0]
	command := os.Args[1]

	switch command {
	case "help":
		runHelp(program)
	case "init":
		if err := commands.RunInit(); err != nil {
			return fmt.Errorf("mvcs init failed: %s", err.Error())
		}
		fmt.Println("mvcs init done")
	case "add":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs add")
		}
		paths := os.Args[2:]
		if err := commands.RunAdd(paths...); err != nil {
			return fmt.Errorf("mvcs add failed: %s", err.Error())
		}
		fmt.Println("mvcs add done")
	case "commit":
		if len(os.Args) < 3 {
			return fmt.Errorf("not enough arguments for mvcs commit")
		}
		msg := os.Args[2]
		if err := commands.RunCommit(msg); err != nil {
			return fmt.Errorf("mvcs commit failed: %s", err.Error())
		}
		fmt.Println("mvcs commit done")
	case "status":
	case "revert":
	case "history":
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
		runHelp(program)
		return fmt.Errorf("Invalid command '%s'", command)
	}

	return nil
}
