# sizepool
A fix size of pool for Go

[![Build Status](https://www.travis-ci.org/FelixSeptem/sizepool.svg?branch=master)](https://www.travis-ci.org/FelixSeptem/sizepool)
[![Coverage Status](https://coveralls.io/repos/github/FelixSeptem/sizepool/badge.svg?branch=master)](https://coveralls.io/github/FelixSeptem/sizepool?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/FelixSeptem/sizepool)](https://goreportcard.com/report/github.com/FelixSeptem/sizepool)
[![GoDoc](http://godoc.org/github.com/FelixSeptem/sizepool?status.svg)](http://godoc.org/github.com/FelixSeptem/sizepool)
[![GolangCI](https://golangci.com/badges/github.com/FelixSeptem/sizepool.svg)](https://golangci.com/r/github.com/FelixSeptem/sizepool)

### The difference between `sync.Pool`
[`sync.Pool`](https://godoc.org/sync#Pool) is aimed at to to cache allocated but unused items for later reuse, relieving pressure on the garbage collector.
And the [sizepool](https://github.com/FelixSeptem/sizepool) is aimed at preallocate a fix size of object size in order to cache a pool of the object(in a list).
You can `Get` and `Put` object into pool just like use `sync.Pool`,the difference is the `sync.Pool` won't block when the pool is empty, it will allocated a new object.
But the `sizepool` will be blocked (or return ErrNoEnoughItem immediately, depends on you use `BGet` or `Get`).
Also, object in `sync.Pool` may be garbage collect,and object in `sizepool` won't be garbage collect.
The thing to be notice is that the size of `sizepool` may not be constant, if your call more time `Put` than `Get` and `BGet`,
your `sizepool` will have the bigger size than your given init size.