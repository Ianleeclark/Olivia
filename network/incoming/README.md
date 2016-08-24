## Incoming Network

For more info on how the incoming network works, see the base README.md
in the network/ folder.

## Commands

1. GET
  - Get allows requests for key/value pairs from a remote node.
2. SET
  - Set allows a remote node/client to set a value in an olivia node.
3. SETEX
  - Setex allows setting a key on an expiration timer. The expiration time
    **must** be in seconds (e.g., "key1:value1:30").
4. REQUEST
  - Request allows requests for different bits of information.
  - Bloomfilter:
    - Allows a remote node/client to request a bloom filter from a remote node.
  - Connect:
    - Allows a remote node/client to request a connection from a remote node.
    - Upon acceptance, both nodes will exchange bloom filters.
  - Peers:
    - Allows a remote node/client to request a peer list from a remote node.
  - Disconnect:
    - Allows a remote node/client to gracefully shutdown.

