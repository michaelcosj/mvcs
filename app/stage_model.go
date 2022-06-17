package app

import (
	"fmt"
	"os"
	"strings"
)

type stageEntry struct {
	path    string
	modTime string
}

type stage []stageEntry

func (sd stage) saveToFile(filepath string) error {
	file, err := openFileWriteOnly(filepath)
	if err != nil {
		return err
	}

	for _, entry := range sd {
		fmt.Fprintf(file, "%s\t%s\n", entry.path, entry.modTime)
	}

	return nil
}

func (sd *stage) addFile(filepath string) error {
	modTime, err := getFileModTime(filepath, TIME_FORMAT)
	if err != nil {
		return err
	}
  
  entry := stageEntry{filepath, modTime}

  if exists, entryIndex := sd.checkEntryExists(entry); exists {
    (*sd)[entryIndex] = entry;
    return nil
  }

  *sd = append(*sd, entry)
  return nil
}

func (sd stage) checkEntryExists(entry stageEntry) (bool, int) {
  for i, e := range sd {
    if e.path == entry.path {
      return true, i
    }
  }

  return false, 0
}

func (sd stage) Len() int {
	return len(sd)
}

func (sd stage) Swap(i, j int) {
	sd[i], sd[j] = sd[j], sd[i]
}
func (sd stage) Less(i, j int) bool {
	return sd[i].path > sd[j].path
}

func getStageFromFile(file string) (*stage, error) {
  data, err := os.ReadFile(file)
  if err != nil {
    return nil, err
  }

  lines := strings.Split(string(data), "\n");
  stage := make(stage, 0)

  for _, line := range lines {
    entryData := strings.Split(line, "\t")
    if len(entryData) > 1 && len(entryData[0]) > 0 {
      stage = append(stage, stageEntry{entryData[0], entryData[1]})
    }
  }

  return &stage, nil
}
