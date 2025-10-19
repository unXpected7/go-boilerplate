package middleware

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type APILogger struct {
	next http.Handler
}

func NewAPLogger(next http.Handler) *APILogger {
	return &APILogger{next: next}
}

func (al *APILogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only log API endpoints, not static files
	if r.URL.Path == "/" {
		al.next.ServeHTTP(w, r)
		return
	}

	start := time.Now()

	// Capture response
	recorder := &responseRecorder{ResponseWriter: w}
	al.next.ServeHTTP(recorder, r)

	duration := time.Since(start)

	// TODO: Add proper logging here if needed
	_ = duration
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       interface{}
	err        string
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}

	// Try to parse JSON response
	var data interface{}
	if json.Unmarshal(b, &data) == nil {
		r.body = data
	} else {
		r.body = string(b)
	}

	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "cannot hijack response")
	}
	return hijacker.Hijack()
}

// Echo middleware version
func APIMonitor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Capture response
		recorder := &echoResponseRecorder{
			ResponseWriter: c.Response(),
			statusCode:     0,
			body:           nil,
			err:            "",
		}

		c.Response().Writer = recorder

		// Call next handler
		err := next(c)

		duration := time.Since(start)

		// TODO: Add proper logging here if needed
		responseData := ""
		errorMsg := ""

		if err != nil {
			errorMsg = err.Error()
		} else if recorder.body != nil {
			if bodyStr, ok := recorder.body.(string); ok {
				responseData = bodyStr
			}
		}

		_ = duration
		_ = responseData
		_ = errorMsg

		return err
	}
}

type echoResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       interface{}
	err        string
}

func (e *echoResponseRecorder) WriteHeader(statusCode int) {
	e.statusCode = statusCode
	e.ResponseWriter.WriteHeader(statusCode)
}

func (e *echoResponseRecorder) Write(b []byte) (int, error) {
	if e.statusCode == 0 {
		e.statusCode = http.StatusOK
	}

	// Try to parse JSON response
	var data interface{}
	if json.Unmarshal(b, &data) == nil {
		e.body = data
	} else {
		e.body = string(b)
	}

	return e.ResponseWriter.Write(b)
}