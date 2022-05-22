package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/value"
)

func main() {
	u := User{Age: value.NewObject(NewAge(10))}
	b, err := json.Marshal(u)
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
		panic(err)
	}
	fmt.Println("john", john)
	fmt.Println("age valid", john.Age.Valid())
	fmt.Println(john.Errors())

	var a *Age
	fmt.Println(a.Validate())

	var a2 *value.Object[*Age]
	fmt.Println(a2.IsZero())
	fmt.Println(a2.Valid())
	fmt.Println(a2.Validate())
}

type User struct {
	Age *value.Object[*Age] `json:"age"`
}

func (u *User) Valid() bool {
	return u.Age.Valid()
}

func (u *User) Errors() map[string]error {
	return map[string]error{
		"age": u.Age.Validate(),
	}
}

var (
	ErrAgeNotSet       = errors.New("age not set")
	ErrInvalidAgeRange = errors.New("invalid age range")
)

// Age value object.
type Age int

func NewAge(age int) *Age {
	v := Age(age)
	return &v
}

func (a *Age) Validate() error {
	if a == nil {
		return ErrAgeNotSet
	}

	n := int(*a)
	if n <= 0 || n > 150 {
		return fmt.Errorf("%w: %d", ErrInvalidAgeRange, n)
	}

	return nil
}
