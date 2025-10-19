package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/env"
	"github.com/go-playground/validator/v10"
)

type TestConfig struct {
	PrimaryEnv string `koanf:"primary_env" validate:"required"`
	ServerPort string `koanf:"server_port" validate:"required"`
}

type NestedConfig struct {
	Primary PrimaryTest `koanf:"primary"`
	Server ServerTest  `koanf:"server"`
}

type PrimaryTest struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerTest struct {
	Port string `koanf:"port" validate:"required"`
}

func main() {
	k := koanf.New(".")

	err := k.Load(env.Provider("BOILERPLATE_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BOILERPLATE_"))
	}), nil)
	if err != nil {
		fmt.Printf("Error loading env: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Environment variables loaded:")
	for key, value := range k.All() {
		fmt.Printf("%s: %v\n", key, value)
	}

	testConfig := &TestConfig{}
	err = k.Unmarshal("", testConfig)
	if err != nil {
		fmt.Printf("Error unmarshalling: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Unmarshalled config: %+v\n", testConfig)

	validate := validator.New()
	err = validate.Struct(testConfig)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Validation passed!")
}