package box

import "errors"

var (
	ErrUnitNotFound = errors.New("unit not found")
)

type Unit string

var (
	UnitMM Unit = "mm"
	UnitCM Unit = "cm"
	UnitM  Unit = "m"
)

func (u *Unit) IsZero() bool {
	return u == nil || *u == ""
}

func (u *Unit) Validate() error {
	if u.IsZero() {
		return ErrUnitNotFound
	}

	switch *u {
	default:
		return ErrUnitNotFound
	case UnitMM, UnitCM, UnitM:
		return nil
	}
}
