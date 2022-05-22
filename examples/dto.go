package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alextanhongpin/value"
)

func main() {
	var dto CreateUserDto

	if err := json.Unmarshal([]byte(`
		{
			"name": "John Appleseed",
			"email": "john.appleseed@mail.com"
		}
	`), &dto); err != nil {
		panic(err)
	}

	createUser(value.Validate(&dto))
}

// Force validation for the dto when passing between layers.
func createUser(dtoContainer value.ToValidate[*CreateUserDto]) {
	dto := dtoContainer.MustValidate()

	fmt.Printf("dto shape: %+v\n", dto)
	fmt.Println("is dto valid?", dto.Valid())
	fmt.Println("dto validation errors:", dto.Errors())
	fmt.Println("name:", dto.Name.MustGet())
	fmt.Println("email:", dto.Email.MustGet())
	fmt.Println("is email valid:", dto.Email.Valid())
	fmt.Println("set email?", dto.Email.Set(NewEmail("john.doemail.com")))
	fmt.Println("get email", dto.Email.MustGet())
}

var (
	ErrEmailNotSet        = errors.New("email not set")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type Email string

func NewEmail(email string) *Email {
	val := Email(email)
	return &val
}

func (e Email) String() string {
	return string(e)
}

func (e *Email) Validate() error {
	if e == nil || *e == "" {
		return ErrEmailNotSet
	}

	if !strings.Contains(e.String(), "@") {
		return ErrInvalidEmailFormat
	}

	return nil
}

var (
	ErrDtoNotSet = errors.New("dto not set")
)

type CreateUserDto struct {
	Name  *value.Value[string]  `json:"name"`
	Email *value.Object[*Email] `json:"email"`
	//Email *Email `json:"email"`
}

func (dto *CreateUserDto) Validate() error {
	if dto == nil {
		return ErrDtoNotSet
	}

	return value.ValidateAll(dto.Name, dto.Email)
}

func (dto *CreateUserDto) Valid() bool {
	return dto.Validate() == nil
}

func (dto *CreateUserDto) Errors() map[string]error {
	return map[string]error{
		"name":  dto.Name.Validate(),
		"email": dto.Email.Validate(),
	}
}
