package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestLoadFromPath_ReadsNodeSpeedLimit(t *testing.T) {
	path := writeTempConfig(t, `{
		"Log": { "Level": "warning" },
		"Nodes": [{
			"ApiHost": "https://example.com/",
			"NodeID": 1,
			"ApiKey": "k",
			"Timeout": 15,
			"SpeedLimit": 50
		}]
	}`)

	c := New()
	if err := c.LoadFromPath(path); err != nil {
		t.Fatalf("LoadFromPath: %v", err)
	}
	if len(c.NodeConfigs) != 1 {
		t.Fatalf("got %d nodes, want 1", len(c.NodeConfigs))
	}
	if c.NodeConfigs[0].SpeedLimit != 50 {
		t.Fatalf("SpeedLimit = %d, want 50", c.NodeConfigs[0].SpeedLimit)
	}
}

func TestLoadFromPath_MissingSpeedLimitDefaultsToZero(t *testing.T) {
	path := writeTempConfig(t, `{
		"Nodes": [{
			"ApiHost": "https://example.com/",
			"NodeID": 1,
			"ApiKey": "k"
		}]
	}`)

	c := New()
	if err := c.LoadFromPath(path); err != nil {
		t.Fatalf("LoadFromPath: %v", err)
	}
	if c.NodeConfigs[0].SpeedLimit != 0 {
		t.Fatalf("SpeedLimit = %d, want 0", c.NodeConfigs[0].SpeedLimit)
	}
}

func TestLoadFromPath_NegativeSpeedLimitClampedToZero(t *testing.T) {
	path := writeTempConfig(t, `{
		"Nodes": [{
			"ApiHost": "https://example.com/",
			"NodeID": 1,
			"ApiKey": "k",
			"SpeedLimit": -10
		}]
	}`)

	c := New()
	if err := c.LoadFromPath(path); err != nil {
		t.Fatalf("LoadFromPath: %v", err)
	}
	if c.NodeConfigs[0].SpeedLimit != 0 {
		t.Fatalf("SpeedLimit = %d, want 0", c.NodeConfigs[0].SpeedLimit)
	}
}
