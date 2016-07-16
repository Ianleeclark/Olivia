package incomingNetwork

// FSMState represents The different states the conn processor will be at
// during a network transaction.
type FSMState int

const (
	UNAUTHENTICATED FSMState = iota
	PROCESSING
)

// ConnProcessor is what handles the processing of each connection coming from
// a remote source. It is a finite state machine which, depending on the state
// level, functions differently.
type ConnProcessor struct {
	State FSMState
}

// NewProcessorFSM Handles creation of a new connection processor.
func NewProcessorFSM(initialState FSMState) *ConnProcessor {
	return &ConnProcessor{
		State: initialState,
	}
}

// ChangeState handles upgrading/downgrading states in the fsm.
func (c *ConnProcessor) ChangeState(nextState FSMState) {
	c.State = nextState
}

// Authenticate handles authentication and upgrading a connection
func (c *ConnProcessor) Authenticate(password string) {
	c.ChangeState(UNAUTHENTICATED + 1)

}
