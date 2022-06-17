package app

import "fmt"

func runCommit() error {
  // parse STAGE file and get list of files to be commited
  stage, err := getStageFromFile(STAGE_FILE)
  if err != nil {
    return err
  }
  fmt.Println(stage)
  // hash file content
  // create metadata for directories and hash it
  // create commit metadata and hash it
  // gzip files and store in BLOBS directory
  // gzip tree metadatas and store in trees directory
  // gzip commit metadatas and store in commit directory
  // clear STAGE
  // update HEAD with the hash of the latest commit

  return nil
}
