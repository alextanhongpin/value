package value

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNotSet       = errors.New("not set")
	ErrInvalidValue = errors.New("invalid value")
)

// Value represents a generic value object.
type Value[T any] struct {
	value T
	dirty bool
}

func New[T any](t T) *Value[T] {
	return &Value[T]{
		value: t,
		dirty: true,
	}
}

func (v *Value[T]) IsZero() bool {
	return v == nil || !v.dirty
}

func (v *Value[T]) Set(t T) error {
	v.value = t
	v.dirty = true

	return nil
}

func (v *Value[T]) Get() (t T, isSet bool) {
	if v.IsZero() {
		return
	}

	return v.value, v.dirty
}

func (v *Value[T]) MustGet() T {
	if v.IsZero() {
		panic(ErrNotSet)
	}

	return v.value
}

func (v *Value[T]) Validate() error {
	if v.IsZero() {
		return ErrNotSet
	}

	return nil
}

func (v *Value[T]) MustValidate() {
	if err := v.Validate(); err != nil {
		panic(err)
	}
}

func (v *Value[T]) Valid() bool {
	return v.Validate() == nil
}

func (v *Value[T]) MustValid() {
	v.MustValidate()
}

func (v *Value[T]) Optional() bool {
	return v.IsZero() || v.Valid()
}

func (v *Value[T]) ValidateOptional() error {
	return nil
}

func (v *Value[T]) String() string {
	if v.IsZero() {
		return "NOT SET"
	}

	return fmt.Sprint(v.value)
}

func (v Value[T]) MarshalJSON() ([]byte, error) {
	if v.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(v.value)
}

func (v *Value[T]) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}

	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}

	v.value = t
	v.dirty = true

	return nil
}
