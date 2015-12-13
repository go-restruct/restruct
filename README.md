# restruct [![Build Status](https://travis-ci.org/johnwchadwick/restruct.svg)](https://travis-ci.org/johnwchadwick/restruct) [![codecov.io](http://codecov.io/github/johnwchadwick/restruct/coverage.svg?branch=master)](http://codecov.io/github/johnwchadwick/restruct?branch=master) [![godoc.org](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/johnwchadwick/restruct)
`restruct` is a library for reading and writing binary data in Go. Similar to
lunixbochs `struc` and `encoding/binary`, this library reads data based on the
layout of structures and, like `struc`, based on what is contained in struct
tags.

`restruct` aims to provide a clean, flexible, robust implementation of struct
packing. In the future, through fast-path optimizations and code generation, it
also aims to be quick, but it is currently very slow.

**Heads up!** This code works best on Go 1.5 and above, because it makes use of
the `reflect.ArrayOf` function added in Go 1.5. See "Arrays" below.

## Status

  * As of writing, coverage is hovering around 95%, but more thorough testing
    is always useful and desirable.
  * Unpacking and packing are fully functional.
  * More optimizations are probably possible.

## Arrays
Restruct supports array types without limitations on all versions of Go it
supports. However, in Go 1.4 and below, there are limitations on overriding
with array types in struct tags because Go 1.4 and below do not provide the
`reflect.ArrayOf` function necessary to dynamically get an array type.

When compiled on Go 1.4 and below, Restruct will use a workaround to support
a limited number of array types in struct tags. You can specify the following
kinds of arrays by default:

  * uint8/byte arrays ranging from 0...128 in length.
  * arrays of any other primitive ranging from 0...32 in length.
  * any array of array from 1...4 lengths of type uint8, uint16, uint32, uint64,
    float32, or float64.

In addition, you can statically register more array types by calling the
[`RegisterArrayType`](https://godoc.org/github.com/johnwchadwick/restruct#RegisterArrayType)
function (this is a no-op on Go 1.5.)

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
