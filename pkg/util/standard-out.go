package util

import "fmt"

func StructPrintln(goStruct interface{})  {
	fmt.Printf("(%%#v) %#v\n", goStruct)
}
