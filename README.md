# Olivia
[![Build Status](https://travis-ci.org/GrappigPanda/Olivia.svg?branch=master)](https://travis-ci.org/GrappigPanda/Olivia)
[![Go Report Card](https://goreportcard.com/badge/github.com/GrappigPanda/Olivia)](https://goreportcard.com/report/github.com/GrappigPanda/Olivia)

Olivia is essentially a distributed hash table that I built to test out some
weird ideas I've had about Go and distributed programming that I've had.
These are just high-level.

```
 1. Implement some sort of gossip protocol to improve node discovery 
    Current node discovery takes O(n^2), but it looks like swim offers O(N)
    Chord DHT offers O(logN) and doesn't entirely change olivia.
 2. Implement some sort of redundancy + consensus. I'm thinking Raft
```

## Deployment

Docker stuff and non-docker stuff

## Is it production ready?
only if ur dumb lul

I can't say for sure if I ever meant for this to be complete and it was
just something which I did after work to de-stress and so that I could
solve actually interesting problems.

## What is Olivia
Olivia is essentially just a distributed hash table with added goodies.
From early on, I wanted to use bloomfilters to enhance lookup speed
between different remote nodes and I'm fairly happy with how that's turned
out. Bloom filters allow us to prioritize which nodes we'll request keys from,
rather than blindly sending out requests to all known nodes.

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
