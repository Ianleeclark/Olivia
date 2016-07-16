package incomingNetwork

import (
	"testing"
)

func TestNewCPFSM(t *testing.T) {
	csfsm := NewProcessorFSM(UNAUTHENTICATED)

	if csfsm.State != UNAUTHENTICATED {
		t.Fatalf("Failed create a new Connection Processor")
	}
}

func TestChangeState(t *testing.T) {
	csfsm := NewProcessorFSM(UNAUTHENTICATED)
	csfsm.ChangeState(PROCESSING)
}

func TestAuthenticate(t *testing.T) {
	csfsm := NewProcessorFSM(UNAUTHENTICATED)
	csfsm.Authenticate("TestBcryptPassword")

	if csfsm.State != PROCESSING {
		t.Fatalf("Failed to authenticate connection.")
	}
}
