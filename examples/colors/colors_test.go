package colors_test

import (
	"encoding/json"
	"testing"

	"github.com/alextanhongpin/value/examples/colors"
)

func TestNewRGB(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		rgb := colors.NewRGB(0, 10, 11)
		if err := rgb.Validate(); err != nil {
			t.Fatalf("expected valid constructor, got %s", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		rgb := colors.NewRGB(-1, 0, 0)
		if err := rgb.Validate(); err == nil {
			t.Fatalf("expected invalid constructor, got %s", err)
		}
	})
}

func TestMarshalRGB(t *testing.T) {
	rgb := colors.NewRGB(255, 255, 255)

	b, err := json.Marshal(rgb)
	if err != nil {
		t.Fatalf("failed to marshal rgb: %s", err)
	}

	expected := `"rgb(255, 255, 255)"`
	if got := string(b); expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}

	var rgb2 colors.RGB
	if err := json.Unmarshal(b, &rgb); err != nil {
		t.Fatalf("failed to unmarshal rgb: %s", err)
	}

	if !rgb2.Equals(rgb2) {
		t.Fatalf("expected %s to match %s", rgb, rgb2)
	}
}
