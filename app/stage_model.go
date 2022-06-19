package app

import (
	"os"
	"strings"
)

// TODO maybe compress the stage

type stageEntry struct {
	path    string
	modTime string
}

type stage []stageEntry

func getStage(stageFile string) (*stage, error) {
	data, err := os.ReadFile(stageFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	stage := make(stage, 0)

	for _, line := range lines {
		entryData := strings.Split(line, "\t")
		if len(entryData) > 1 && len(entryData[0]) > 0 {
			stage = append(stage, stageEntry{entryData[0], entryData[1]})
		}
	}
	return &stage, nil
}

func (sd *stage) addFile(filepath string) error {
	modTime, err := getFileModTime(filepath, TIME_FORMAT)
	if err != nil {
		return err
	}

	entry := stageEntry{filepath, modTime}

	if exists, entryIndex := sd.checkEntryExists(entry); exists {
		(*sd)[entryIndex] = entry
		return nil
	}

	*sd = append(*sd, entry)
	return nil
}

func (sd stage) saveToFile(path string) error {
  var data strings.Builder
	for _, entry := range sd {
    data.WriteString(entry.path + "\t" + entry.modTime + "\n")
	}
	return writeToFile(path, data.String())
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
