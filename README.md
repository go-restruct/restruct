# restruct [![Build Status](https://travis-ci.org/johnwchadwick/restruct.svg)](https://travis-ci.org/johnwchadwick/restruct) [![codecov.io](http://codecov.io/github/johnwchadwick/restruct/coverage.svg?branch=master)](http://codecov.io/github/johnwchadwick/restruct?branch=master)
`restruct` is a work-in-progress library for reading and writing binary data in
Go. Similar to lunixbochs `struc` and `encoding/binary`, this library reads data
based on the layout of structures and, like `struc`, based on what is contained
in struct tags.

**Heads up!** This code relies on Go 1.5, because it relies on a new function
added to the reflect package (`reflect.ArrayOf`.)

## Priorities

  * __Features first__: Performance is a secondary concern. First, the program
    needs to work well.
  * __Test early, test often__: This project aims for 100% coverage.
  * __Flexibility__: Like struc, it is important that we support variable-length
    strings and slices. It is also important that we can use embedded structs
	and slices of structs.

## Example (WIP)
> This code does not work yet.

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
