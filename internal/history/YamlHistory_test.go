package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestYamlHistorySuite(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "botman-yamlhistory")
	os.RemoveAll(tmp)
	os.Mkdir(tmp, 0700)

	RunSuite(t, func() HistoryKeeper {
		return &YamlHistory{Path: tmp}
	})
}
