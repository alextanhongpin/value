package value

import (
	"bytes"
	"encoding/json"
	"errors"
)

var ErrObjectNotSet = errors.New("value object not set")

type Object[T validatable] struct {
	_     struct{}
	value T
	dirty bool
	// Allow setting null error here (might not work when serializing/deserializatin though).
}

func NewObject[T validatable](t T) *Object[T] {
	// Allow returning invalid object, so that the validation can be deferred in parent Validate method.
	return &Object[T]{
		value: t,
		dirty: true,
	}
}

func (o *Object[T]) IsZero() bool {
	return o == nil || !o.dirty
}

func (o *Object[T]) Validate() error {
	if o.IsZero() {
		return ErrObjectNotSet
	}

	return o.value.Validate()
}

func (o *Object[T]) MustValidate() {
	if err := o.Validate(); err != nil {
		panic(err)
	}
}

func (o *Object[T]) Valid() bool {
	return o.Validate() == nil
}

func (o *Object[T]) MustValid() {
	o.MustValidate()
}

func (o *Object[T]) ValidateOptional() error {
	if o.IsZero() {
		return nil
	}

	return o.value.Validate()
}

func (o *Object[T]) Optional() bool {
	return o.IsZero() || o.Valid()
}

func (o *Object[T]) Get() (t T, set bool) {
	if o.IsZero() {
		return
	}

	return o.value, o.dirty
}

func (o *Object[T]) MustGet() T {
	if o.IsZero() {
		panic(ErrObjectNotSet)
	}

	return o.value
}

func (o *Object[T]) Set(t T) error {
	if err := t.Validate(); err != nil {
		return err
	}

	o.value = t
	o.dirty = true

	return nil
}

func (o *Object[T]) MarshalJSON() ([]byte, error) {
	if o.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(o.value)
}

func (o *Object[T]) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}

	var t T
	if err := json.Unmarshal(raw, &t); err != nil {
		return err
	}

	o.value = t
	o.dirty = true

	return nil
}
