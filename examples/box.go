package main

import (
	"fmt"

	"github.com/alextanhongpin/value"
	"github.com/alextanhongpin/value/examples/box"
)

func calculateBoxVolume(boxContainer value.ToValidate[*box.Box]) (volume value.ToValidate[*box.Dimension]) {
	box, err := boxContainer.Validate()
	if err != nil {
		return
	}

	return value.Validate(box.Volume())
}

type Cargo struct {
	box value.Object[*box.Box]
}

func main() {
	dim := box.New(
		box.NewDimension(10, box.UnitCM),
		box.NewDimension(11, box.UnitCM),
		box.NewDimension(12, box.UnitCM),
	)

	volume := calculateBoxVolume(value.Validate(dim))
	vol, err := volume.Validate()
	if err != nil {
		panic(err)
	}
	fmt.Println(vol)

	cargo := &Cargo{}
	fmt.Println("is cargo box present?", !cargo.box.IsZero())
	fmt.Println("is cargo box valid?", cargo.box.Valid())
	cargo.box.Set(dim)
	fmt.Println("is cargo box present?", !cargo.box.IsZero())
	fmt.Println("is cargo box valid?", cargo.box.Valid())
}
