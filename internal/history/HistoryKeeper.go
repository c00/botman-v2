package history

import "github.com/c00/botman-v2/internal/logger"

var log = logger.New("historyKeeper")

type HistoryKeeper interface {
	SaveChat(entry HistoryEntry) (HistoryEntry, error)
	LoadChat(lookback int) (HistoryEntry, error)
	List() ([]string, error)
}
