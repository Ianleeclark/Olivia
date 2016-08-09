package networkHandler

import (
	"testing"
	"time"
)

func TestExecuteRepeatedly(t *testing.T) {
	countChan := make(chan int)
	killChan := make(chan bool)

	go executeRepeatedly(
		5*time.Millisecond,
		func() { return },
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
		break
	default:
		return
	}

	if count != 10 {
		t.Errorf("Expected 10, got %v", count)
	}
}
