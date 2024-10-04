package history

import (
	"errors"
)

type InMemoryHistory struct {
	entries []HistoryEntry
}

func (h *InMemoryHistory) SaveChat(entry HistoryEntry) (HistoryEntry, error) {
	//Find existing
	for i, e := range h.entries {
		if e.Name == entry.Name {
			h.entries[i] = entry
			return entry, nil
		}
	}

	h.entries = append(h.entries, entry)
	return entry, nil
}

func (h *InMemoryHistory) LoadChat(lookback int) (HistoryEntry, error) {
	length := len(h.entries)
	if length == 0 {
		return HistoryEntry{}, errors.New("no history")
	}

	if lookback >= length {
		return HistoryEntry{}, errors.New("lookback is further than available entries")
	}

	//Get the looback index
	index := length - lookback - 1
	return h.entries[index], nil
}

func (h *InMemoryHistory) List() ([]string, error) {
	names := []string{}
	for _, item := range h.entries {
		names = append(names, item.Name)
	}

	return names, nil
}
