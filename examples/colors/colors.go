package colors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/value"
)

var ErrChannelOutOfRange = fmt.Errorf("%w: rgb must be between 0 and 255", value.ErrInvalidValue)

func ValidateChannelRange(n int) error {
	if n < 0 || n > 255 {
		return ErrChannelOutOfRange
	}

	return nil
}

func ValidateRGBTuple(rgb RGBTuple) error {
	if err := ValidateChannelRange(rgb.R); err != nil {
		return err
	}

	if err := ValidateChannelRange(rgb.G); err != nil {
		return err
	}

	if err := ValidateChannelRange(rgb.B); err != nil {
		return err
	}

	return nil
}

type RGBTuple struct {
	R int
	G int
	B int
}

func (rgb RGBTuple) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", rgb.R, rgb.G, rgb.B)
}

type RGB struct {
	// Ideally avoid embedding pointer, as you have to initialize both the RGB struct as well as the embedded struct.
	value.Value[RGBTuple]
}

func NewRGB(r, g, b int) (*RGB, error) {
	val, _ := value.New(RGBTuple{r, g, b}, value.WithValidator(ValidateRGBTuple))
	rgb := &RGB{*val}

	return rgb, rgb.Validate()
}

func (r *RGB) Validate() error {
	if r == nil {
		return value.ErrNotSet
	}

	return r.Value.Validate()
}

func (r *RGB) Valid() bool {
	return r.Validate() == nil
}

func (r RGB) String() string {
	if !r.Valid() {
		return "INVALID RGB"
	}

	rgb := r.MustGet()

	return rgb.String()
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
