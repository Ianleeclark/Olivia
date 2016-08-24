## Receiver

Essentially all this does is wait for a response on a channel, look up the
hash value in the response channel, and then call the message handler to
handle callbacks for the response.

This is the other half of the message handler.
