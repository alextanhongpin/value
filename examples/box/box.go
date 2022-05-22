package box

import (
	"errors"
	"fmt"
)

var (
	ErrNoBox                    = errors.New("no box")
	ErrInconsistentBoxDimension = errors.New("inconsistent box dimension")
)

type Box struct {
	Length *Dimension
	Width  *Dimension
	Height *Dimension
}

func New(length, width, height *Dimension) *Box {
	return &Box{
		Length: length,
		Width:  width,
		Height: height,
	}
}

func (b *Box) Validate() error {
	if b == nil {
		return ErrNoBox
	}

	if b.Length == nil {
		return fmt.Errorf("%w: box length", ErrDimensionNotSet)
	}

	if err := b.Length.Validate(); err != nil {
		return fmt.Errorf("%w: box length", err)
	}

	if b.Width == nil {
		return fmt.Errorf("%w: box width", ErrDimensionNotSet)
	}

	if err := b.Width.Validate(); err != nil {
		return fmt.Errorf("%w: box width", err)
	}

	if b.Height == nil {
		return fmt.Errorf("%w: box height", ErrDimensionNotSet)
	}

	if err := b.Height.Validate(); err != nil {
		return fmt.Errorf("%w: box height", err)
	}

	units := make(map[Unit]int)
	units[b.Length.Unit]++
	units[b.Height.Unit]++
	units[b.Width.Unit]++
	if len(units) != 1 && units[b.Length.Unit] != 3 {
		return ErrInconsistentBoxDimension
	}

	return nil
}

func (b *Box) Valid() bool {
	return b.Validate() == nil
}

func (b *Box) Volume() *Dimension {
	if !b.Valid() {
		return nil
	}

	length, unit := b.Length.Value, b.Length.Unit
	width := b.Width.Value
	height := b.Height.Value

	return NewDimension(length*width*height, unit)
}
