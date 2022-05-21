package box

import "github.com/alextanhongpin/value"

type BoxDimension struct {
	Length Dimension
	Width  Dimension
	Height Dimension
}

func ValidateBoxDimension(b *BoxDimension) error {
	if b == nil {
		return ErrEmptyDimension
	}

	return b.Length.Validate()
}

type Box struct {
	value.Value[*BoxDimension]
}
