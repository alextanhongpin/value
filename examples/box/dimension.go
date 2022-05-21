package box

import (
	"errors"
	"fmt"

	"github.com/alextanhongpin/value"
)

var ErrEmptyDimension = errors.New("empty dimension")

type DimensionTuple struct {
	Unit  Unit
	Value value.Value[uint]
}

func NewDimensionTuple(val uint, unit Unit) *DimensionTuple {
	return &DimensionTuple{
		Value: *value.Must(value.New(val)),
		Unit:  unit,
	}
}

func (tup DimensionTuple) String() string {
	return fmt.Sprintf("%s%s", tup.Value.String(), tup.Unit)
}

func ValidateDimensionTuple(tuple *DimensionTuple) error {
	if tuple == nil {
		return ErrEmptyDimension
	}

	if err := tuple.Unit.Validate(); err != nil {
		return err
	}

	if err := tuple.Value.Validate(); err != nil {
		return err
	}

	return nil
}

type Dimension struct {
	value.Value[*DimensionTuple]
}

func NewDimension(val uint, unit Unit) (*Dimension, error) {
	tup := value.Must(value.New(&DimensionTuple{
		Value: *value.Must(value.New(val)),
		Unit:  unit,
	}, value.WithValidator(ValidateDimensionTuple)))

	dim := &Dimension{*tup}

	return dim, dim.Validate()
}
