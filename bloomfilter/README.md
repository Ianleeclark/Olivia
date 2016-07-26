# Bloom Filter

A quick primer if you don't know what a bloom filter is:
https://en.wikipedia.org/wiki/Bloom_filter

# How Bloom Filters Help Olivia

Bloom filters are used in Olivia to aide in selecting a peer to retrieve a key
not found in the current node. So, if a key is requested from an node and the
node doesn't contain the value, the operating node will search all known peers
for the key. Since latency is a thing which exists, bloom filters are used to
improve searches. 

Once a node is made aware of another node, each node will send a copy of their
bloom filter (which will continued to be updated).

