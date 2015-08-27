# restruct [![Build Status](https://travis-ci.org/johnwchadwick/restruct.svg)](https://travis-ci.org/johnwchadwick/restruct) [![codecov.io](http://codecov.io/github/johnwchadwick/restruct/coverage.svg?branch=master)](http://codecov.io/github/johnwchadwick/restruct?branch=master)
`restruct` is a library for reading and writing binary data in Go. Similar to
lunixbochs `struc` and `encoding/binary`, this library reads data based on the
layout of structures and, like `struc`, based on what is contained in struct
tags.

`restruct` aims to provide a clean, flexible, robust implementation of struct
packing. In the future, through fast-path optimizations and code generation, it
also aims to be quick, but it is currently very slow.

**Heads up!** This code relies on Go 1.5, because it relies on a new function
added to the reflect package (`reflect.ArrayOf`.)

## Status

  * As of writing, coverage is 100%. This means every line can work in some
    cases, but many more tests and assertions are needed to cover edge cases.
  * Unpacking and packing are fully functional.
  * Performance is poor, because the library is unoptimized.
  * The library needs more documentation.

## Example

```go
package main

import (
	"encoding/binary"
	"io/ioutil"
	"os"

	"github.com/johnwchadwick/restruct"
)

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

	file, _ := os.Open("records")
	defer file.Close()
	data, _ := ioutil.ReadAll(file)

	restruct.Unpack(data, binary.LittleEndian, &c)
}
```
