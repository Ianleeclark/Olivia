# The Network layer

# How the Network Layer works for incoming connections.

There is a thin layer in each Olivia node which allows remote connections.
Whenever a new connection is established, the network layer will divert each
connection to its own operating goroutine. Incoming connections may be from
other olivia nodes or applications requesting cached values.

Inside each goroutine, the connection will be handled by a connection
processing finite state machine (FSM). The FSM operates in several states
allowing varying states and processing paths in each goroutine.

Upon a remote node connecting, a `REQUEST connect` command will be sent and any
operations which are necessary will happen: currently (0.1.x) we just
send/request bloom filters.

# How the Network Layer works for outgoing connections.

There's two ways. We have the ability to just send a normal command through
a peer connection, or we can use the message_handler + receiver to send a
non-blocking request.
