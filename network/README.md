# The Network layer

# How the Network Layer works for incoming connections.

There is a thin layer in each Olivia node which allows remote connections.
Whenever a new connection is established, the network layer will divert each
connection to its own operating goroutine. Incoming connections may be from
other olivia nodes or applications requesting cached values.

Inside each goroutine, the connection will be handled by a connection
processing finite state machine (FSM). The FSM operates in several states
allowing varying states and processing paths in each goroutine.

# How the Network Layer works for outgoing connections.

Heh. TBD
