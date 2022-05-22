package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/value"
	"github.com/alextanhongpin/value/examples/colors"
)

var ErrHomePageNotFound = errors.New("homepage not found")

type HomePage struct {
	Width           int                        `json:"width"`
	Height          int                        `json:"height"`
	BackgroundColor *value.Object[*colors.RGB] `json:"backgroundColor"`
}

func (p *HomePage) Validate() error {
	if p == nil {
		return ErrHomePageNotFound
	}

	return p.BackgroundColor.ValidateOptional()
}

func (p *HomePage) Valid() bool {
	return p.Validate() == nil
}

func main() {
	var homepage HomePage
	raw := []byte(`
		{
			"width": 1096,
			"height": 360,
			"backgroundColor": "rgb(255, 255, 255)"
		}
	`)
	if err := json.Unmarshal(raw, &homepage); err != nil {
		panic(err)
	}

	fmt.Println("is homepage valid?", homepage.Valid())
	fmt.Println("homepage background color:")
	fmt.Println(homepage.BackgroundColor.Get())
	fmt.Println("is homepage background color valid?:", homepage.BackgroundColor.Valid())
	fmt.Println("is homepage background color optional?", homepage.BackgroundColor.Optional())

	// Set to red background.
	fmt.Println()
	fmt.Println("set to red background?:", homepage.BackgroundColor.Set(&colors.RGB{R: 255}))
	fmt.Println("is homepage valid?", homepage.Valid())
	fmt.Println("homepage background color:")
	fmt.Println(homepage.BackgroundColor.Get())
	fmt.Println("is homepage background color valid?:", homepage.BackgroundColor.Valid())
	fmt.Println("is homepage background color optional?", homepage.BackgroundColor.Optional())
}
