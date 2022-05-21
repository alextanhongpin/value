package main

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/value/examples/colors"
)

func main() {
	var rgb colors.RGB
	if err := json.Unmarshal(json.RawMessage(`"rgb(0, 100, 255)"`), &rgb); err != nil {
		panic(err)
	}
	fmt.Println("rgb:", rgb)
	fmt.Println("valid:", rgb.Valid())
	fmt.Println("validate:", rgb.Validate())

	b, err := json.Marshal(rgb)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	var rgb2 *colors.RGB
	if !rgb2.Valid() {
		rgb2 = new(colors.RGB)
		err := rgb2.Set(colors.RGBTuple{0, 128, 255})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(rgb2.Valid())
	fmt.Println(rgb2.Get())
	rgbTuple := rgb2.MustGet()
	fmt.Println("r:", rgbTuple.R)
	fmt.Println("g:", rgbTuple.G)
	fmt.Println("b:", rgbTuple.B)
}
