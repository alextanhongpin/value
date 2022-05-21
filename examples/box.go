package main

import "github.com/alextanhongpin/value"
import "errors"
import "fmt"

// NOTE: In short, don't use embedding. Keeping the code clean takes priority over embedding.
func main() {
	//var dim *Dimension
	dim, err := NewDimension(100, "cm")
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
}

var ErrUnitNotFound = errors.New("unit not found")

func ValidateUnit(v string) error {
	for _, m := range []string{"mm", "cm", "m"} {
		if m == v {
			return nil
		}
	}
	return ErrUnitNotFound
}

//type Unit value.Value[string] // This does not have embedded methods.
type Unit struct {
	value.Value[string]
}

func NewUnit(unit string) (*Unit, error) {
	val, err := value.New(unit, value.WithValidator(ValidateUnit))
	if err != nil {
		return nil, err
	}
	res := Unit{*val}
	return &res, nil
}

type DimensionValue struct {
	Unit  Unit
	Value value.Value[int] // Should allow only positive values.
}

func (v DimensionValue) String() string {
	return fmt.Sprintf("%s%s", v.Value.String(), v.Unit.String())
}

func ValidateDimensionValue(dim DimensionValue) error {
	return AnyError(dim.Unit.Validate, dim.Value.Validate)
}

type Dimension struct {
	value.Value[DimensionValue]
}

func (d *Dimension) Validate() error {
	if d == nil {
		return value.ErrNotSet
	}
	return d.Value.Validate()
}

func (d *Dimension) IsZero() bool {
	return d == nil || d.Value.IsZero()
}

func (d *Dimension) IsSet() bool {
	return !d.IsZero() && d.Value.IsSet()
}

func (d *Dimension) Valid() bool {
	if d == nil {
		return false
	}
	return d.Value.Valid()
}

func (d Dimension) String() string {
	if d.IsZero() {
		return "NOT SET"
	}
	v, _ := d.Get()
	return v.String()
}

func NewDimension(val int, unit string) (*Dimension, error) {
	u, err := NewUnit(unit)
	if err != nil {
		return nil, err
	}
	v, err := value.New(val)
	if err != nil {
		return nil, err
	}

	dim, err := value.New(DimensionValue{Unit: *u, Value: *v}, value.WithValidator(ValidateDimensionValue))
	if err != nil {
		return nil, err
	}

	return &Dimension{*dim}, nil
}

func AnyError(errFns ...func() error) error {
	for _, fn := range errFns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
