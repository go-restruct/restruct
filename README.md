# restruct [![Build Status](https://travis-ci.org/johnwchadwick/restruct.svg)](https://travis-ci.org/johnwchadwick/restruct) [![codecov.io](http://codecov.io/github/johnwchadwick/restruct/coverage.svg?branch=master)](http://codecov.io/github/johnwchadwick/restruct?branch=master)
`restruct` is a work-in-progress library for reading and writing binary data in
Go. Similar to lunixbochs `struc` and `encoding/binary`, this library reads data
based on the layout of structures and, like `struc`, based on what is contained
in struct tags.

**Heads up!** This code relies on Go 1.5, because it relies on a new function
added to the reflect package (`reflect.ArrayOf`.)

## Status

  * All of the code needs more testing.
  * A preliminary implementation of unpacking was created. Most of it is
    covered by testing, but we need a lot more assertions made.
  * There is no packing implementation yet.
  * Struct tags specifying type overrides, sizeof fields, byte order and skip
    values are implemented and functional for decoding.
  * Performance is bad. This could be remedied with caching, careful profiling,
    and hopefully at some point, code generation for packing/unpacking. Still,
    if parsing binary data is not your bottleneck, this package should do just
    fine.

## Priorities

  * __Features first__: Performance is a secondary concern. First, the program
    needs to work well.
  * __Test early, test often__: This project aims for 100% coverage.
  * __Flexibility__: Like struc, it is important that we support variable-length
    strings and slices. It is also important that we can use embedded structs
	and slices of structs.

## Example (WIP)

```go
package main

import "os"
import "github.com/johnwchadwick/restruct"

type Record struct {
	Message string `struct:[128]byte`
}

type Container struct {
	Version   int `struct:int32`
	NumRecord int `struct:int32,sizeof=Records`
	Records   []Record
}

func main() {
	var c Container
	file, err := os.Open("records")
	restruct.Unpack(file, &c)
}
```
