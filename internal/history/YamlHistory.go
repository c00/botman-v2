package history

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

func NewYamlHistory(path string) *YamlHistory {
	os.Mkdir(path, 0700)
	return &YamlHistory{Path: path}
}

type YamlHistory struct {
	Path string
}

func (h *YamlHistory) SaveChat(entry HistoryEntry) (HistoryEntry, error) {
	//Save that to disk as yaml
	err := h.saveFile(entry)
	if err != nil {
		log.Error("could not save history entry: %v", err)
		return HistoryEntry{}, err
	}

	return entry, nil
}

func (h *YamlHistory) saveFile(entry HistoryEntry) error {
	if len(entry.Messages) == 0 {
		log.Debug("No messages. Not saving entry")
		return nil
	}

	if entry.Name == "" {
		return errors.New("entry does not have a name")
	}

	savePath := filepath.Join(h.Path, toSafeFilename(entry.Name))
	log.Debug("Saving history entry to %v", savePath)

	dir := filepath.Dir(savePath)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%v should be a directory", dir)
	}

	bytes, err := yaml.Marshal(entry)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, bytes, 0700)
	if err != nil {
		return err
	}

	return nil
}

func (h *YamlHistory) LoadChat(lookback int) (HistoryEntry, error) {
	files, err := h.List()
	if err != nil {
		return HistoryEntry{}, err
	}

	length := len(files)
	if length == 0 {
		return HistoryEntry{}, errors.New("no history")
	}

	if lookback >= length {
		return HistoryEntry{}, errors.New("lookback is further than available entries")
	}

	//Get the looback index
	index := length - lookback - 1
	return h.loadFile(files[index])

}

func (h *YamlHistory) loadFile(path string) (HistoryEntry, error) {
	//Read the history dir
	filePath := filepath.Join(h.Path, path)
	log.Debug("Loading history entry from %v", filePath)

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return HistoryEntry{}, err
	}

	entry := HistoryEntry{}
	err = yaml.Unmarshal(bytes, &entry)
	if err != nil {
		return HistoryEntry{}, err
	}

	return entry, nil
}

func (h *YamlHistory) List() ([]string, error) {
	//Read the history dir
	historyPath := filepath.Join(h.Path)
	files := []string{}

	err := filepath.Walk(historyPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Check if the filename ends in .yaml, if not continue
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

func toSafeFilename(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, s) + ".yaml"
}
