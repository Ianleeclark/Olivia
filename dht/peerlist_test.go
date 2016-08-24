package dht

import (
	"github.com/GrappigPanda/Olivia/config"
	"testing"
)

var CONFIG = config.ReadConfig()

// This is just to ensure that nothing blows up.
func TestNewPeerList(t *testing.T) {
	NewPeerList(nil, *CONFIG)
}
