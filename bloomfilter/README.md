# Bloom Filter

A quick primer if you don't know what a bloom filter is:
https://en.wikipedia.org/wiki/Bloom_filter

# How Bloom Filters Help Olivia

Bloom filters are used in Olivia to aide in selecting a peer to retrieve a key
not found in the current node. So, if a key is requested from an node and the
node doesn't contain the value, the operating node will search all known peers
for the key. Since latency is a thing which exists, bloom filters are used to
improve search times so we can only query nodes which **probably** have the
key.

Once a node is made aware of another node, each node will send a copy of their
bloom filter (which will continued to be updated). Moreover, on a timed
interval (default of 30-seconds), a new updated bloom filter will be
transferred between each node. 

This does mean there's a 30 second window where any updates to bloom filters
are not known by remote nodes. This is a **problem** and potential solutions 
are being thought out. The currently courted (as of writing this) idea/solution
is to send a list of updated hashes with each heartbeat. This allows quicker
updating the bloom filters by simple O(1) insertions.

The bloomfilter is backed by a third-party library
(https://github.com/willf/bitset) and will continue to be so for the
foreseeable future. To cut down on network burden, all bloom filters are
marshalled to JSON and then run-length encoded. This tends to heavily cut down
on total size of data being transmitted.
