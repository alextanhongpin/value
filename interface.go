package value

type validatable interface {
	Validate() error
}

type ToValidate[T validatable] interface {
	Validate() (T, error)
	MustValidate() T
}

type validate[T validatable] struct {
	value T
}

func (v *validate[T]) Validate() (T, error) {
	return v.value, v.value.Validate()
}

func (v *validate[T]) MustValidate() T {
	if err := v.value.Validate(); err != nil {
		panic(err)
	}
	return v.value
}

func Validate[T validatable](t T) *validate[T] {
	return &validate[T]{value: t}
}

func ValidateAll(val ...validatable) error {
	for _, v := range val {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
