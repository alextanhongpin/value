// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

func main() {
	var d *Dimension
	fmt.Println(d.Validate())
	fmt.Println(d.Valid())
	fmt.Println(d)

	dd, err := NewDimension(100, UnitMM)
	if err != nil {
		panic(err)
	}
	fmt.Println(dd)
}

type Unit string

var (
	UnitMM Unit = "mm"
	UnitCM Unit = "cm"
	UnitM  Unit = "m"
)

func (u Unit) Validate() error {
	switch u {
	case UnitMM, UnitCM, UnitM:
		return nil
	default:
		return errors.New("unit not found")
	}
}
func (u Unit) Valid() bool {
	return u.Validate() == nil
}

type Dimension struct {
	Unit  Unit
	Value uint
	dirty bool
}

func NewDimension(value uint, unit Unit) (*Dimension, error) {
	dim := &Dimension{Value: value, Unit: unit, dirty: true}
	if err := dim.Validate(); err != nil {
		return nil, err
	}
	return dim, nil
}

func (d *Dimension) IsZero() bool {
	return d == nil || *d == Dimension{}
}

func (d *Dimension) IsSet() bool {
	return !d.IsZero() && d.dirty
}

func (d *Dimension) Validate() error {
	if d.IsZero() {
		return errors.New("not set")
	}
	if err := d.Unit.Validate(); err != nil {
		return err
	}
	if d.Value == 0 {
		return errors.New("zero dimension")
	}
	return nil
}

func (d *Dimension) Valid() bool {
	return d.Validate() == nil
}

func (d Dimension) String() string {
	if d.IsZero() {
		return "no dimension specified"
	}
	return fmt.Sprintf("%d%s", d.Value, d.Unit)
}
