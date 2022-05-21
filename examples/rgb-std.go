package main

import (
	"errors"
	"fmt"
)

func main() {
	var rgb *RGB
	fmt.Println(rgb.Valid(), rgb.Validate())
	rgb, err := NewRGB(0, 0, 1)
	if err != nil {
		panic(err)
	}
	rgb.R = -100
	fmt.Println(rgb)
	fmt.Println(rgb.Validate())
}

type RGB struct {
	R     int
	G     int
	B     int
	dirty bool
}

func NewRGB(r, g, b int) (*RGB, error) {
	rgb := &RGB{R: r, G: g, B: b, dirty: true}
	if err := rgb.Validate(); err != nil {
		return nil, err
	}
	return rgb, nil
}

func (r *RGB) IsZero() bool {
	return r == nil || !r.dirty
}
func (r *RGB) IsSet() bool {
	return !r.IsZero() && r.dirty
}

func (r *RGB) validateChannel(channel int) error {
	if channel < 0 || channel > 255 {
		return errors.New("out of range")
	}
	return nil
}

func (r *RGB) Validate() error {
	if r.IsZero() {
		return errors.New("not set")
	}
	if err := r.validateChannel(r.R); err != nil {
		return err
	}
	if err := r.validateChannel(r.G); err != nil {
		return err
	}
	if err := r.validateChannel(r.B); err != nil {
		return err
	}
	return nil
}
func (r *RGB) Valid() bool {
	return r.Validate() == nil
}

func (r RGB) String() string {
	if r.IsZero() {
		return "NOT SET"
	}
	return fmt.Sprintf("rgb(%d, %d, %d)", r.R, r.G, r.B)
}
