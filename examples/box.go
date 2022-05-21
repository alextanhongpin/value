package main

import (
	"fmt"

	"github.com/alextanhongpin/value"
	"github.com/alextanhongpin/value/examples/box"
)

func main() {
	dim, err := box.NewDimension(100, box.UnitMM)
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

	if err := dim.Set(box.NewDimensionTuple(250, box.UnitCM)); err != nil {
		panic(err)
	}
	fmt.Println(dim)

	b, err := box.New(&box.Shape{
		Length: value.MustDeref(box.NewDimension(10, box.UnitCM)),
		Width:  value.MustDeref(box.NewDimension(11, box.UnitCM)),
		Height: value.MustDeref(box.NewDimension(12, box.UnitCM)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(b.Volume())
	fmt.Println(b)
	fmt.Println(b.Valid())
	fmt.Println(b.Validate())
}
