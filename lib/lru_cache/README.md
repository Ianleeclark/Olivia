# LRU Cache

A quick primer: https://en.wikipedia.org/wiki/Cache_algorithms under section
`Least Recently Used (LRU)`

# How LRU Caches are used in Olivia

As of the current time in writing this, LRU Caches are used solely to memoize
key hashes so that every insertion/retrieval of a key from the bloom filters
don't need to always hash the key `x` times.
