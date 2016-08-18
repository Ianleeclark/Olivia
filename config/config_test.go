package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := ReadConfig()

	// Please note: If this is failing, it's probably due to the fact that the
	// config file was edited.
	if cfg.BaseNode != true {
		t.Errorf("Expected true, got %v", cfg.BaseNode)
	}

	if cfg.BloomfilterSize != 1000 {
		t.Errorf("Expected 1000, got %v", cfg.BloomfilterSize)
	}

	if cfg.HeartbeatInterval != 1000 {
		t.Errorf("Expected 1000, got %v", cfg.HeartbeatInterval)
	}

	if cfg.HeartbeatLoop != 30 {
		t.Errorf("Expected 30, got %v", cfg.HeartbeatLoop)
	}

	for _, peer := range cfg.RemotePeers {
		if peer != "127.0.0.1:5455" {
			t.Errorf("Expected 127.0.0.1:5455, got %v", peer)
		} else {
			break
		}
	}

	if cfg.ListenPort != 5454 {
		t.Errorf("Expected 5454, got %v", cfg.HeartbeatLoop)
	}

}
