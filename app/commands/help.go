package commands

import "fmt"

// Run help command
func RunHelp(program string) {
	fmt.Println("Usage:", program, "<command>")
	fmt.Println(
		"MVCS - Mediocre Version Control System\n\n" +
			"    help                               display this help and exit \n" +
			"    init name=<username> email=<email> initialze mvcs in a repo \n" +
			"    status                             show the status of working tree \n" +
			"    add <file|directory>               add a file or directory to staging area \n" +
			"    commit <message>                   commit all files or directories in the staging area \n" +
			"    checkout <commit hash>             revert back to the specified commit \n" +
			"    history                            show all commits and their info \n" +
			"    clean                              clean all orphaned commits",
	)
}

