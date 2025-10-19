package testing

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/config"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/database"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"
)

// CreateTestServer creates a server instance for testing
func CreateTestServer(logger *zerolog.Logger, db *TestDB) *server.Server {
	// Set up observability config with defaults if not present
	if db.Config.Observability == nil {
		db.Config.Observability = &config.ObservabilityConfig{
			ServiceName: "alfred-test",
			Environment: "test",
			Logging: config.LoggingConfig{
				Level:              "info",
				Format:             "json",
				SlowQueryThreshold: 100 * time.Millisecond,
			},
			NewRelic: config.NewRelicConfig{
				LicenseKey:                "",    // Empty for tests
				AppLogForwardingEnabled:   false, // Disabled for tests  
				DistributedTracingEnabled: false, // Disabled for tests
				DebugLogging:              false, // Disabled for tests
			},
			HealthChecks: config.HealthChecksConfig{
				Enabled: false,
			},
		}
	}

	testServer := &server.Server{
		Logger: logger,
		DB: &database.Database{
			Pool: db.Pool,
		},
		Config: db.Config,
	}

	return testServer
}