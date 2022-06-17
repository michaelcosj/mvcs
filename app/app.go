package app

import (
	"errors"
	"fmt"
	"os"
)

const (
	// mvcs paths
	MVCS_DIR   = ".mvcs"
	COMMIT_DIR = MVCS_DIR + "/commits"
	TREE_DIR   = MVCS_DIR + "/trees"
	BLOB_DIR   = MVCS_DIR + "/blobs"

	// mvcs files
	HEAD_FILE  = MVCS_DIR + "/HEAD"
	STAGE_FILE = MVCS_DIR + "/STAGE"
	CFG_FILE   = MVCS_DIR + "/config"
	LOGS_FILE  = MVCS_DIR + "/logs"

	TIME_FORMAT = "2006-01-02|15:04:05"
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
		if err := runInit(); err != nil {
			return fmt.Errorf("mvcs init failed: %s", err.Error())
		}
		fmt.Println("mvcs init done")
	case "add":
		paths := os.Args[2:]
		if err := runAdd(paths...); err != nil {
			return fmt.Errorf("mvcs add failed: %s", err.Error())
		}
		fmt.Println("mvcs add done")
	case "commit":
		if err := runCommit(); err != nil {
			return fmt.Errorf("mvcs commit failed: %s", err.Error())
		}
		fmt.Println("mvcs commit done")
	case "status":
	case "revert":
	case "history":
	default:
		runHelp(program)
		return fmt.Errorf("Invalid command '%s'", command)
	}

	return nil
}
