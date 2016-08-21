## Message Handler

I'm confused by this code. I got too experimental and I'm passing channels of
channels around so that I can do callbacks from (mostly) the receiver.

Essentially what happens is we receive a response in the receiver, look up
the hash value in the the message receiver by using a channel, and we get a
callback channel. We can send the response through the callback channel.

This allows us to send requests and _eventually_ handle the response rather
than wait. There's a high likelihood I rip this code out in the future, but at
the same time, it's extremely useful for sending requests which I expect to
take an indefinite amount of time.
