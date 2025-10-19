package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//go:embed swagger-ui/*
var swaggerUI embed.FS

func main() {
	// List all embedded files
	fmt.Println("Listing embedded files:")
	if files, err := fs.ReadDir(swaggerUI, "swagger-ui"); err != nil {
		fmt.Printf("Error reading dir: %v\n", err)
		os.Exit(1)
	} else {
		for _, file := range files {
			fmt.Printf("  %s\n", file.Name())
		}
	}

	// Try to open index.html directly
	fmt.Println("\nTrying to open index.html directly:")
	if content, err := swaggerUI.Open("internal/handler/swagger-ui/index.html"); err != nil {
		fmt.Printf("Error opening: %v\n", err)
	} else {
		fmt.Println("Successfully opened index.html directly!")
		content.Close()
	}

	// Try with sub
	fmt.Println("\nTrying with fs.Sub:")
	file, err := fs.Sub(swaggerUI, "swagger-ui")
	if err != nil {
		fmt.Printf("Error creating sub: %v\n", err)
		os.Exit(1)
	}

	if contentBytes, err := file.Open("index.html"); err != nil {
		fmt.Printf("Error opening index.html: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Successfully loaded index.html with sub!")
		contentBytes.Close()
	}
}