# Olivia
[![Build Status](https://travis-ci.org/GrappigPanda/Olivia.svg?branch=master)](https://travis-ci.org/GrappigPanda/Olivia)
[![Go Report Card](https://goreportcard.com/badge/github.com/GrappigPanda/Olivia)](https://goreportcard.com/report/github.com/GrappigPanda/Olivia)

Olivia is essentially a distributed hash table that I built to test out some
weird ideas I've had about Go and distributed programming that I've had.

## Is it production ready?
I'd definitely consider this **not** production ready. I typically approach
personal projects in two phases: the initial naive phase, where I implement
things in the most natural (for me) way; and the improvement phase, where I
exchange naive solutions for more optimal solutions. I'm heavily entrenched
into the naive phase, so while I'm getting to the optimal solution phase, I'm
not there yet.

## What is Olivia
Olivia is essentially just a distributed hash table with added goodies.
From early on, I wanted to use bloomfilters to enhance lookup speed
between different remote nodes and I'm fairly happy with how that's turned
out. Bloom filters allow us to prioritize which nodes we'll request keys from,
rather than blindly sending out requests to all known nodes.

## What Olivia could become
I'm still not 100% sure on the end vision. When I originally started out
building this thing, I just wanted three things:
  - Distribution
  - Remote key lookups (consistency between each node has never been huge)
  - Key Expiration

As I build it, however, my desired are changing. I now want consistency between
nodes, but I haven't yet decided how I'll go about it. Olivia was originally
intended to be part of a distributed torrent tracker network, but it is no
longer a worthwhile pursuit, as distributed hash tables (as a part of the
Bittorrent network) are able to be edited (upon BEP acceptance).

## Contact Maintainer

[open an issue](https://github.com/GrappigPanda/notorious/issues/new)

[tweet me](http://twitter.com/GrappigPanda)

[or email me](mailto:ian@ianleeclark.com)

## License

The MIT License (MIT)

Copyright (c) 2016 Ian Clark

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
