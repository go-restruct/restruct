// +build ignore

package main

import "fmt"

var head = `// +build !go1.5

package restruct

import "runtime"

func init() {
`
var wrhead = "\t{\n"
var wrfoot = "\t}\n\truntime.GC()\n"
var tpl = "\t\tRegisterArrayType([%d]%s{})\n"
var tpl2 = "\t\tRegisterArrayType([%d][%d]%s{})\n"
var foot = "}\n"

var types = []string{
	"bool",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"float32",
	"float64",
	"complex64",
	"complex128",
	"string",
}

func main() {
	fmt.Print(head)

	// All primitives up to 32 elements
	for _, typ := range types {
		fmt.Print(wrhead)
		for i := 0; i <= 32; i++ {
			fmt.Printf(tpl, i, typ)
		}
		fmt.Print(wrfoot)
	}

	// Longer byte arrays
	fmt.Print(wrhead)
	for i := 33; i <= 128; i++ {
		fmt.Printf(tpl, i, "uint8")
	}
	fmt.Print(wrfoot)

	// Common array-of-array types.
	fmt.Print(wrhead)
	for i := 1; i <= 4; i++ {
		for j := 1; j <= 4; j++ {
			fmt.Printf(tpl2, i, j, "float32")
			fmt.Printf(tpl2, i, j, "float64")
			fmt.Printf(tpl2, i, j, "uint8")
			fmt.Printf(tpl2, i, j, "uint16")
			fmt.Printf(tpl2, i, j, "uint32")
			fmt.Printf(tpl2, i, j, "uint64")
		}
	}
	fmt.Print(wrfoot)

	fmt.Print(foot)
}
