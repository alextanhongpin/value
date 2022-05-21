package box

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/value"
)

var (
	ErrInconsistentBoxDimension = errors.New("inconsistent box dimension")
	ErrNoBox                    = errors.New("no box")
)

type Shape struct {
	Length Dimension
	Width  Dimension
	Height Dimension
}

func NewShape(length, width, height uint, unit Unit) *Shape {
	return &Shape{
		Length: value.MustDeref(NewDimension(length, unit)),
		Width:  value.MustDeref(NewDimension(width, unit)),
		Height: value.MustDeref(NewDimension(height, unit)),
	}
}

func (s Shape) String() string {
	return fmt.Sprintf("%s x %s x %s", s.Length.String(), s.Width.String(), s.Height.String())
}

func ValidateBoxDimension(shape *Shape) error {
	if shape == nil {
		return ErrEmptyDimension
	}

	if err := shape.Length.Validate(); err != nil {
		return err
	}

	if err := shape.Width.Validate(); err != nil {
		return err
	}

	if err := shape.Height.Validate(); err != nil {
		return err
	}

	length := shape.Length.MustGet()
	width := shape.Width.MustGet()
	height := shape.Height.MustGet()

	if length.Unit != width.Unit && width.Unit != height.Unit {
		return ErrInconsistentBoxDimension
	}

	return nil
}

type Box struct {
	value.Value[*Shape]
}

func New(shape *Shape) (*Box, error) {
	val, _ := value.New(
		shape,
		value.WithValidator(ValidateBoxDimension),
	)
	box := &Box{*val}

	return box, box.Validate()
}

func (b *Box) Validate() error {
	if b == nil {
		return ErrNoBox
	}

	return b.Value.Validate()
}

func (b *Box) Valid() bool {
	return b.Validate() == nil
}

func (b *Box) Volume() (*Dimension, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}

	dim := b.MustGet()
	length := dim.Length.MustGet().Value.MustGet()
	width := dim.Width.MustGet().Value.MustGet()
	height := dim.Height.MustGet().Value.MustGet()

	return NewDimension(length*width*height, dim.Length.MustGet().Unit)
}
