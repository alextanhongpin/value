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
	value     T
	dirty     bool
	validator func(T) error
}

type Option[T any] func(*Value[T]) *Value[T]

func WithValidator[T any](validator func(T) error) Option[T] {
	return func(v *Value[T]) *Value[T] {
		v.validator = validator

		return v
	}
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

type validatable interface {
	Validate() error
}

func Validate[T validatable](vals ...T) error {
	for _, val := range vals {
		if err := val.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func MustDeref[T any](t *T, err error) T {
	if err != nil {
		return Deref(new(T))
	}

	return Deref(t)
}

func Deref[T any](t *T) T {
	if t == nil {
		var t T

		return t
	}

	return *t
}

type optional interface {
	IsZero() bool
	Valid() bool
}

func Optional[T optional](t T) bool {
	return t.IsZero() || t.Valid()
}

func New[T any](t T, options ...Option[T]) (*Value[T], error) {
	val := &Value[T]{
		value:     t,
		dirty:     true,
		validator: nil,
	}

	for _, opt := range options {
		opt(val)
	}

	// Allow returning invalid value object, to allow deferred validation.
	return val, val.Validate()
}

func (v *Value[T]) IsZero() bool {
	return v == nil || !v.dirty
}

func (v *Value[T]) IsSet() bool {
	return !v.IsZero() && v.dirty
}

func (v *Value[T]) SetValidator(fn func(T) error) {
	v.validator = fn
}

func (v *Value[T]) With(t T) (*Value[T], error) {
	if err := v.validate(t); err != nil {
		return v, err
	}

	return New(t, WithValidator(v.validator))
}

func (v *Value[T]) Set(t T) error {
	if err := v.validate(t); err != nil {
		return err
	}

	v.value = t
	v.dirty = true

	return nil
}

func (v *Value[T]) Get() (t T, isSet bool) {
	if !v.IsSet() {
		return
	}

	return v.value, v.dirty
}

func (v *Value[T]) MustGet() T {
	if !v.IsSet() {
		panic(ErrNotSet)
	}

	return v.value
}

func (v *Value[T]) validate(t T) error {
	if validate := v.validator; validate != nil {
		return validate(t)
	}

	return nil
}

func (v *Value[T]) Validate() error {
	if v.IsZero() {
		return ErrNotSet
	}

	return v.validate(v.value)
}

func (v *Value[T]) Valid() bool {
	return v.Validate() == nil
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

// UnmarshalJSON does not add back the validator - figure
// out how to add it back through reflection (NOTE:
// manually add in value object, see Age example).
func (v *Value[T]) UnmarshalJSON(raw []byte) error {
	if v == nil {
		return errors.New("unmarshal to nil Value[T]")
	}

	if bytes.Equal(raw, []byte("null")) {
		return nil
	}

	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}

	*v = *Must(New(t, WithValidator(v.validator)))
	return nil
}
