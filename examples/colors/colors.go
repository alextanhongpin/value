package colors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrChannelOutOfRange = errors.New("channel out of range")

type RGB struct {
	_       struct{}
	R, G, B int
}

func NewRGB(r, g, b int) *RGB {
	return &RGB{
		R: r,
		G: g,
		B: b,
	}
}

func (rgb *RGB) validateChannel(channel int) error {
	if channel < 0 || channel > 255 {
		return ErrChannelOutOfRange
	}

	return nil
}

func (rgb *RGB) Validate() error {
	if err := rgb.validateChannel(rgb.R); err != nil {
		return fmt.Errorf("%w: R", err)
	}

	if err := rgb.validateChannel(rgb.G); err != nil {
		return fmt.Errorf("%w: G", err)
	}

	if err := rgb.validateChannel(rgb.B); err != nil {
		return fmt.Errorf("%w: B", err)
	}

	return nil
}

func (r RGB) String() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", r.R, r.G, r.B)
}

func (r RGB) Equals(other RGB) bool {
	return true &&
		r.R == other.R &&
		r.G == other.G &&
		r.B == other.B
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

	*r = *NewRGB(ri, gi, bi)
	return nil
}
