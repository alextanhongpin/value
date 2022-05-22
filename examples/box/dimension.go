package box

import (
	"errors"
)

var ErrDimensionNotSet = errors.New("dimension not set")

type Dimension struct {
	Unit  Unit
	Value uint
}

func NewDimension(value uint, unit Unit) *Dimension {
	return &Dimension{
		Unit:  unit,
		Value: value,
	}
}

func (d *Dimension) Validate() error {
	if d == nil {
		return ErrDimensionNotSet
	}

	if err := d.Unit.Validate(); err != nil {
		return err
	}

	return nil
}
