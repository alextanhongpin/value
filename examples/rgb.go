package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/value"
)

func main() {
	var rgbValue Channel
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

	var rgb2 *RGB
	if !rgb2.Valid() {
		rgb, err := NewRGB(255, 255, 255)
		if err != nil {
			panic(err)
		}
		rgb2 = rgb
	}
	if err := rgb2.R.Set(100); err != nil {
		panic(err)
	}
	fmt.Println(rgb2.Valid())
	fmt.Println(rgb2)
}

var ErrChannelOutOfRange = fmt.Errorf("%w: rgb must be between 0 and 255", value.ErrInvalidValue)

func ValidateChannelRange(n int) error {
	if n < 0 || n > 255 {
		return ErrChannelOutOfRange
	}
	return nil
}

type Channel struct {
	value.Value[int]
}

func NewChannel(rgb int) (*Channel, error) {
	value, err := value.New(rgb, value.WithValidator(ValidateChannelRange))
	if err != nil {
		return nil, err
	}

	return &Channel{*value}, nil
}

func (r *Channel) Validate() error {
	if r == nil {
		return value.ErrNotSet
	}
	return r.Value.Validate()
}

func (r *Channel) Valid() bool {
	return r.Validate() == nil
}

type RGB struct {
	R *Channel
	G *Channel
	B *Channel
}

func NewRGB(r, g, b int) (*RGB, error) {
	rv, _ := NewChannel(r)
	gv, _ := NewChannel(g)
	bv, _ := NewChannel(b)
	rgb := &RGB{
		R: rv,
		G: gv,
		B: bv,
	}
	return rgb, rgb.Validate()
}

func (r *RGB) Validate() error {
	if r == nil {
		return value.ErrNotSet
	}
	return AnyError(r.R.Validate, r.R.Validate, r.B.Validate)
}

func (r *RGB) Valid() bool {
	if r == nil {
		return false
	}
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
