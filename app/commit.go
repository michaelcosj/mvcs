package app

import "os"

func runCommit(msg string) error {
	stg, err := getStage(STAGE_FILE)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	tree := newTreeData(workingDir)
	for _, stagedFile := range *stg {
		fileData, err := newFileData(stagedFile.path)
		if err != nil {
			return err
		}

		tree.addChild(fileData)
	}
	if err := tree.generateHash(); err != nil {
		return err
	}

	parentCommit, err := getFileContent(HEAD_FILE)
	if err != nil {
		return err
	}

	commit := newCommitData(parentCommit, msg, tree)
	commit.compressData(OBJ_DIR)

	if err := clearFiles(STAGE_FILE, HEAD_FILE); err != nil {
		return err
	}

	return writeToFile(HEAD_FILE, commit.hash)
}
