## dht

Stuff you'll find here are dht implementation details. Peer information,
peer lists, &c. There isn't a specific DHT algorithm which we're using, but the
general workflow for how our peers communicate can be found in the network
folder.

We essentially hold 3 peers as important nodes: any other peers we encounter
and receive from our 3 important nodes are added to a backuplist. The three
important nodes are set on a quick heartbeat, whereas each other node will have
an artery clogged hearbeat every minute. Each peer operates on a FSM with
multiple states. If a peer is continually not responding to queries, it will be
set to a timed out state, but will continued to be beaten with our heart.

Peer lists are typically used for heartbeats and peer-wide operations: for
instance, whenever a request for a key not found in the current node is made,
we'll iterate through each peer in the peerlist and see if that probably has
the key.
