## Message Handler

I'm confused by this code. I got too experimental and I'm passing channels of
channels around so that I can do callbacks from (mostly) the receiver.

Essentially what happens is we receive a response in the receiver, look up
the hash value in the the message receiver by using a channel, and we get a
callback channel. We can send the response through the callback channel.
