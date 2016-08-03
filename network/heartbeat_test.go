package networkHandler

import (
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	Heartbeat(1*time.Second, 60*time.Second)
}

func TestExecuteRepeatedly(t *testing.T) {
	countChan := make(chan int)
	killChan := make(chan bool)

	go executeRepeatedly(
		1 * time.Millisecond,
		func() {return},
		killChan,
		countChan,
	)

	count := 0
	select {
	case x := <-countChan:
		count += x
		if count == 10 {
			killChan <- true
		}
		break;
	default:
		return
	}

	if count != 10 {
		t.Errorf("Expected 10, got %v", count)
	}
}
