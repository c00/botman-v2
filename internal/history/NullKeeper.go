package history

import "errors"

// Does not safe history at all. Just a stub to satisfy the interface.
// Is used when SaveHistory is set to false
type NullKeeper struct{}

func (h *NullKeeper) SaveChat(entry HistoryEntry) (HistoryEntry, error) {
	return entry, nil
}

func (h *NullKeeper) LoadChat(lookback int) (HistoryEntry, error) {
	return HistoryEntry{}, errors.New("nullkeeper does not keep history")

}

func (h *NullKeeper) List() ([]string, error) {
	return []string{}, nil
}
