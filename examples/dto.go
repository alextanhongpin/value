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
			"email": "john.appleseed@mail.com",
			"address": {
				"street1": "street 1",
				"street2": "street 2",
				"city": "Singapore",
				"state": "N/A",
				"country": "Singapore",
				"postalCode": "-"
			}
		}
	`), &dto); err != nil {
		panic(err)
	}
	fmt.Printf("dto shape: %+v\n", dto)

	createUser(value.Validate(&dto))
}

// Force validation for the dto when passing between layers.
func createUser(dtoContainer value.ToValidate[*CreateUserDto]) {
	dto := dtoContainer.MustValidate()

	fmt.Println("is dto valid?", dto.Valid())
	fmt.Println("dto validation errors:", dto.Errors())
	fmt.Println("name:", dto.Name.MustGet())
	fmt.Println("email:", dto.Email.MustGet())
	fmt.Println("is email valid:", dto.Email.Valid())
	fmt.Println("set email?", dto.Email.Set(NewEmail("john.doemail.com")))
	fmt.Println("get email", dto.Email.MustGet())
	fmt.Println("get address", dto.Address.MustGet())
}

// Email value object.

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

// Address value object.

var (
	ErrAddressNotSet     = errors.New("address not set")
	ErrIncompleteAddress = errors.New("incomplete address")
)

type Address struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

func (a *Address) Validate() error {
	if a == nil {
		return ErrAddressNotSet
	}
	if a.Street1 == "" {
		return fmt.Errorf("%w: street1", ErrIncompleteAddress)
	}
	if a.Street2 == "" {
		return fmt.Errorf("%w: street2", ErrIncompleteAddress)
	}
	if a.City == "" {
		return fmt.Errorf("%w: city", ErrIncompleteAddress)
	}
	if a.State == "" {
		return fmt.Errorf("%w: state", ErrIncompleteAddress)
	}
	if a.PostalCode == "" {
		return fmt.Errorf("%w: postal code", ErrIncompleteAddress)
	}
	if a.Country == "" {
		return fmt.Errorf("%w: country", ErrIncompleteAddress)
	}

	return nil
}

// Dto.

var (
	ErrDtoNotSet = errors.New("dto not set")
)

type CreateUserDto struct {
	Name  *value.Value[string]  `json:"name"`
	Email *value.Object[*Email] `json:"email"`
	// Embedding this fails...
	// Cannot embed more than 2 Object, must create another type with different
	// names.
	Address *value.Object[*Address] `json:"address"`
}

func (dto *CreateUserDto) Validate() error {
	if dto == nil {
		return ErrDtoNotSet
	}

	return value.ValidateAll(dto.Name, dto.Email, dto.Address)
}

func (dto *CreateUserDto) Valid() bool {
	return dto.Validate() == nil
}

func (dto *CreateUserDto) Errors() map[string]error {
	return map[string]error{
		"name":    dto.Name.Validate(),
		"email":   dto.Email.Validate(),
		"address": dto.Address.Validate(),
	}
}
