package models

import (
	"michaelcosj/mvcs/app/constants"
	"michaelcosj/mvcs/app/helpers"
	"os"
	"strings"
)

// TODO maybe compress the stage

type StageEntry struct {
	Path    string
	modTime string
}

type Stage []StageEntry

func GetStage() (*Stage, error) {
	data, err := os.ReadFile(constants.STAGE_FILE)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	stage := make(Stage, 0)

	for _, line := range lines {
		entryData := strings.Split(line, "\t")
		if len(entryData) > 1 && len(entryData[0]) > 0 {
			stage = append(stage, StageEntry{entryData[0], entryData[1]})
		}
	}
	return &stage, nil
}

func (sd *Stage) AddFile(file string) error {
	modTime, err := helpers.GetFileModTime(file, constants.TIME_FORMAT)
	if err != nil {
		return err
	}

	entry := StageEntry{file, modTime}

	if exists, entryIndex := sd.checkEntryExists(entry); exists {
		(*sd)[entryIndex] = entry
		return nil
	}

	*sd = append(*sd, entry)
	return nil
}

func (sd Stage) Save() error {
  var data strings.Builder
	for _, entry := range sd {
    data.WriteString(entry.Path + "\t" + entry.modTime + "\n")
	}
	return helpers.WriteToFile(constants.STAGE_FILE, data.String())
}

func (sd Stage) checkEntryExists(entry StageEntry) (bool, int) {
	for i, e := range sd {
		if e.Path == entry.Path {
			return true, i
		}
	}
	return false, 0
}

func (sd Stage) Files() []string {
  files := make([]string, 0)
	for _, stagedFile := range sd {
    files = append(files, stagedFile.Path)
	}
  return files
}

func (sd Stage) Len() int {
	return len(sd)
}

func (sd Stage) Swap(i, j int) {
	sd[i], sd[j] = sd[j], sd[i]
}

func (sd Stage) Less(i, j int) bool {
	return sd[i].Path > sd[j].Path
}
