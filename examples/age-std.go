package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotSet          = errors.New("not set")
	ErrInvalidAgeRange = errors.New("invalid age range")
)

func main() {
	age, err := NewAge(100)
	if err != nil {
		panic(err)
	}
	fmt.Println(age)
	fmt.Println(age.Set(-100))
}

type Age struct {
	value int
	dirty bool
}

func NewAge(age int) (*Age, error) {
	a := &Age{value: age, dirty: true}
	return a, a.Validate()
}

func (a *Age) Validate() error {
	if a == nil || !a.dirty {
		return ErrNotSet
	}

	return a.validate(a.value)
}

func (a *Age) validate(age int) error {
	if age < 0 || age > 150 {
		return ErrInvalidAgeRange
	}

	return nil
}

func (a *Age) Valid() bool {
	return a.Validate() == nil
}

func (a *Age) Set(age int) error {
	if err := a.validate(age); err != nil {
		return err
	}

	a.value = age
	a.dirty = true

	return nil
}

func (a Age) String() string {
	if !a.Valid() {
		return "invalid age"
	}

	return fmt.Sprint(a.value)
}
