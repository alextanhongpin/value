package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/value"
)

func main() {
	var rgbValue RGBValue
	if err := json.Unmarshal([]byte("200"), &rgbValue); err != nil {
		panic(err)
	}
	fmt.Println(rgbValue.Get())
	fmt.Println(rgbValue.Validate())
	fmt.Println(rgbValue.Valid())

	var rgb RGB
	if err := json.Unmarshal(json.RawMessage(`"rgb(0, 100, 255)"`), &rgb); err != nil {
		panic(err)
	}
	fmt.Println(rgb)
	fmt.Println(rgb.Valid())
	fmt.Println(rgb.Validate())
}

var ErrRGBValueOutOfRange = fmt.Errorf("%w: rgb must be between 0 and 255", value.ErrInvalidValue)

func ValidateRGBValueRange(n int) error {
	if n < 0 || n > 255 {
		return ErrRGBValueOutOfRange
	}
	return nil
}

type RGBValue struct {
	value.Value[int]
}

func NewRGBValue(rgb int) (*RGBValue, error) {
	value, err := value.New(rgb, value.WithValidator(ValidateRGBValueRange))
	if err != nil {
		return nil, err
	}

	return &RGBValue{*value}, nil
}

func (r *RGBValue) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}

	v := new(value.Value[int])
	if err := json.Unmarshal(raw, v); err != nil {
		return err
	}
	v.SetValidator(ValidateRGBValueRange)
	r.Value = *v
	return nil
}

type RGB struct {
	R *RGBValue
	G *RGBValue
	B *RGBValue
}

func NewRGB(r, g, b int) (*RGB, error) {
	rv, err := NewRGBValue(r)
	if err != nil {
		return nil, err
	}
	gv, err := NewRGBValue(g)
	if err != nil {
		return nil, err
	}
	bv, err := NewRGBValue(b)
	if err != nil {
		return nil, err
	}
	return &RGB{
		R: rv,
		G: gv,
		B: bv,
	}, nil
}

func (r *RGB) Validate() error {
	return AnyError(r.R.Validate, r.R.Validate, r.B.Validate)
}

func (r *RGB) Valid() bool {
	return AllTrue(r.R.Valid, r.G.Valid, r.B.Valid)
}

func (r *RGB) String() string {
	if !r.Valid() {
		return "INVALID RGB"
	}
	return fmt.Sprintf("rgb(%s, %s, %s)", r.R, r.G, r.B)
}

func (r RGB) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RGB) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}
	var rgb string
	if err := json.Unmarshal(raw, &rgb); err != nil {
		return err
	}
	var ri, gi, bi int
	_, err := fmt.Sscanf(rgb, "rgb(%d,%d,%d)", &ri, &gi, &bi)
	if err != nil {
		return err
	}
	r2, err := NewRGB(ri, gi, bi)
	if err != nil {
		return err
	}
	*r = *r2
	return nil
}

func AnyError(errFns ...func() error) error {
	for _, fn := range errFns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func AllTrue(boolFns ...func() bool) bool {
	for _, fn := range boolFns {
		if valid := fn(); !valid {
			return false
		}
	}
	return true
}
