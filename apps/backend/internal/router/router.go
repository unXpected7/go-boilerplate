package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/handler"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/service"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	// For debugging, create a super simple echo server
	router := echo.New()

	// Add a single test route
	router.GET("/", func(c echo.Context) error {
		return c.String(200, "SUPER SIMPLE TEST ROUTE - THIS SHOULD WORK!")
	})

	
	s.Logger.Info().Msg("SUPER SIMPLE ROUTER CREATED")
	s.Logger.Info().Msg("ROUTER CREATED")

	// Let's also create a standalone HTTP server for testing
	standaloneServer := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	go standaloneServer.ListenAndServe()
	s.Logger.Info().Msg("Standalone server started on :8081")

	return router
}

