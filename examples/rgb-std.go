package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrNotSet            = errors.New("not set")
	ErrChannelOutOfRange = errors.New("channel out of range")
)

func main() {
	var rgb *RGB
	fmt.Println(rgb.Valid())
	fmt.Println(rgb.Validate())

	rgb, err := NewRGB(0, 0, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(rgb.R, rgb.G, rgb.B)
	fmt.Println(rgb.Valid())
	fmt.Println(rgb.Validate())
	SaveRGB(NewValidate(rgb))
}

type Validatable interface {
	Validate() error
}
type ToValidate[T Validatable] interface {
	Validate() (T, error)
}

type Validate[T Validatable] struct {
	value T
	dirty bool
}

func NewValidate[T Validatable](t T) *Validate[T] {
	return &Validate[T]{value: t, dirty: true}
}

func (v *Validate[T]) Validate() (T, error) {
	if v == nil || !v.dirty {
		var t T
		return t, ErrNotFound
	}
	return v.value, v.value.Validate()
}

// Enforce that the repository layer must call the Validate() method before obtaining the value of RGB.
func SaveRGB(rgbVal ToValidate[*RGB]) {
	rgb, err := rgbVal.Validate()
	if err != nil {
		panic(err)
	}
	fmt.Println("saving rgb:", rgb)
}

type Channel struct {
	value int
	dirty bool
}

func NewChannel(value int) (*Channel, error) {
	channel := &Channel{value: value, dirty: true}
	// Return the invalid channel, so that the
	// validation can be delayed when we have multiple fields.
	return channel, channel.Validate()
}

func (c *Channel) IsEmpty() bool {
	return c == nil || !c.dirty
}

func (c *Channel) Validate() error {
	if c.IsEmpty() {
		return ErrNotSet
	}

	if c.value < 0 || c.value > 255 {
		return ErrChannelOutOfRange
	}

	return nil
}

func (c *Channel) Valid() bool {
	return c.Validate() == nil
}
func (c Channel) String() string {
	return fmt.Sprint(c.value)
}

type RGB struct {
	R     *Channel
	G     *Channel
	B     *Channel
	dirty bool
}

func NewRGB(r, g, b int) (*RGB, error) {
	rc, _ := NewChannel(r)
	gc, _ := NewChannel(g)
	bc, _ := NewChannel(b)
	rgb := &RGB{R: rc, G: gc, B: bc, dirty: true}
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

func (r *RGB) Validate() error {
	if r.IsZero() {
		return ErrNotSet
	}
	if err := r.R.Validate(); err != nil {
		return fmt.Errorf("%w: r", err)
	}
	if err := r.G.Validate(); err != nil {
		return fmt.Errorf("%w: g", err)
	}
	if err := r.B.Validate(); err != nil {
		return fmt.Errorf("%w: b", err)
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
	return fmt.Sprintf("rgb(%s, %s, %s)", r.R, r.G, r.B)
}
