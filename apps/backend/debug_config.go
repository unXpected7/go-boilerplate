package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/env"
)

type SimpleConfig struct {
	PrimaryEnv    string `koanf:"primary_env"`
	ServerPort    string `koanf:"server_port"`
	DatabaseHost  string `koanf:"database_host"`
	DatabasePort  int    `koanf:"database_port"`
	DatabaseUser  string `koanf:"database_user"`
	DatabaseName  string `koanf:"database_name"`
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

	fmt.Println("All loaded keys:")
	for key, value := range k.All() {
		fmt.Printf("%s: %v\n", key, value)
	}

	config := &SimpleConfig{}
	err = k.Unmarshal("", config)
	if err != nil {
		fmt.Printf("Error unmarshalling: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Unmarshalled config: %+v\n", config)
}