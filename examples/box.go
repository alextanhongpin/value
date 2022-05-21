package main

import (
	"fmt"

	"github.com/alextanhongpin/value/examples/box"
)

func main() {
	dim, err := box.NewDimension(100, "mm")
	if err != nil {
		panic(err)
	}
	fmt.Println(dim.String())
	fmt.Println(dim.Get())
	fmt.Println("isZero:", dim.IsZero())
	fmt.Println("isSet:", dim.IsSet())
	fmt.Println("valid:", dim.Valid())
	fmt.Println("validate:", dim.Validate())
	fmt.Println(dim)

	if err := dim.Set(box.NewDimensionTuple(250, "cm")); err != nil {
		panic(err)
	}
	fmt.Println(dim)
}
