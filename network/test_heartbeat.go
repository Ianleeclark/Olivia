package networkHandler

import (
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	Heartbeat(1*time.Second, 60*time.Second)
}
