package history

import (
	"testing"
)

func TestInMemoryHistorySuite(t *testing.T) {

	RunSuite(t, func() HistoryKeeper {
		return &InMemoryHistory{}
	})
}
