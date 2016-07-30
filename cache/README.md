## Cache

This is where operations for the cache reside. Mostly handles setting/getting
values. Honestly, it's a super simple and ugly way of doing it and needs to be
handled better in the future.

In the future, I want our COW to be more robust and virile. Right now, our COW
is made to wait on each write as we iterate over every N elements in the
underlying cache. I'd prefer to construct and change the pointer as quickly as
possible. I could probably just lock, set value, unlock, goroutine update the
cache and return. While I'm updating the cache, I don't need to lock, I can just
wait until the new cache is completely constructed and then lock, change
pointer, unlock. Yeah, I should do that.

Beyond that, I want to allow key expirations. I can use my binary heap that I
created to order by expiration time and on each heartbeat update, check the
root node. If it's at our current time or before, we can expire the key.
