package handler

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed swagger-ui/*
var swaggerUI embed.FS

//go:embed swagger-ui.json
var swaggerSpec []byte

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// Serve swagger UI
func (h *SwaggerHandler) ServeSwaggerUI(c echo.Context) error {
	// Extract path without the /docs prefix
	path := strings.TrimPrefix(c.Request().URL.Path, "/docs")

	// If path is empty, serve the index.html
	if path == "" {
		path = "index.html"
	}

	// Open the file from embedded filesystem
	file, err := fs.Sub(swaggerUI, "swagger-ui")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading swagger UI")
	}

	contentBytes, err := file.Open(path)
	if err != nil {
		return c.String(http.StatusNotFound, "Swagger UI not found")
	}
	defer contentBytes.Close()

	// Read content into bytes
	content, err := io.ReadAll(contentBytes)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading swagger UI content")
	}

	// Get file info for content type (we don't need size since we're using c.Blob)
	_, err = contentBytes.Stat()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting file info")
	}

	// Set content type based on file extension
	switch {
	case strings.HasSuffix(path, ".html"):
		c.Response().Header().Set("Content-Type", "text/html")
	case strings.HasSuffix(path, ".js"):
		c.Response().Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(path, ".css"):
		c.Response().Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(path, ".png"):
		c.Response().Header().Set("Content-Type", "image/png")
	}

	// Serve content using c.Blob instead of http.ServeContent
	return c.Blob(http.StatusOK, "", content)
}

// Serve OpenAPI spec
func (h *SwaggerHandler) ServeOpenAPISpec(c echo.Context) error {
	return c.Blob(http.StatusOK, "application/json", swaggerSpec)
}