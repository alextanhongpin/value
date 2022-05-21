package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/value"
)

func main() {
	age, err := NewAge(10)
	if err != nil {
		panic(err)
	}
	fmt.Println(age.Valid())
	fmt.Println(age.IsZero())
	fmt.Println(age.Validate())
	fmt.Println(age.Set(-10))
	fmt.Println(age.Get())
	fmt.Println(age.MustGet())
	fmt.Println(age)

	b, err := json.Marshal(age)
	if err != nil {
		panic(err)
	}
	fmt.Println("marshall", string(b))

	u := User{Age: age}
	b, err = json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("user", string(b))

	u = User{}
	b, err = json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println("user", string(b))

	var john User
	if err := json.Unmarshal([]byte(`{"age": -1}`), &john); err != nil {
		if errors.Is(err, value.ErrInvalidValue) {
			// Check if this is error due to value object, which
			// will be a validation error.
			fmt.Println("unmarshal error", err)
		} else {
			panic(err)
		}
	}
	fmt.Println("john", john)
	fmt.Println("age valid", john.Age.Valid())
	fmt.Println(john.Errors())

	// Unfortunately embedding the value object doesn't help much.
	// We still need to check all nil values, which makes no difference
	// then implementing each value object individually.
	var age2 *Age
	fmt.Println(age2.Valid())
}

type User struct {
	Age *Age `json:"age"`
}

func (u *User) Valid() bool {
	return u.Age.Valid()
}

func (u *User) Errors() map[string]error {
	return map[string]error{
		"age": u.Age.Validate(),
	}
}

// Age value object.
type Age struct {
	// *Value[int] Don't use pointer, there will be issue with unmarshalling
	value.Value[int]
}

func (a *Age) UnmarshalJSON(raw []byte) error {
	if a == nil {
		// TODO
	}
	var v value.Value[int]
	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}
	// Set back the age validator manually here.
	v.SetValidator(ValidateAge)
	a.Value = v

	// Additionally perform validation here
	//return a.Validate()
	return nil
}

var ErrInvalidAgeRange = fmt.Errorf("%w: invalid age", value.ErrInvalidValue)

func ValidateAge(age int) error {
	if age < 0 {
		return ErrInvalidAgeRange
	}
	return nil

}
func NewAge(age int) (*Age, error) {
	value, err := value.New(age, value.WithValidator(ValidateAge))
	if err != nil {
		return nil, err
	}
	return &Age{*value}, nil
}
