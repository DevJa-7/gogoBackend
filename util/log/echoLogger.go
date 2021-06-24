package log

import (
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"../../config"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// GetContextLogInfo gets context log information
func GetContextLogInfo(c echo.Context) (string, int, string, string) {
	method := c.Request().Method
	statusCode := c.Response().Status
	urlPath := c.Request().URL.Path
	clientIP := c.Request().RemoteAddr
	return method, statusCode, urlPath, clientIP
}

// AccessLogger write access log into a file.
func AccessLogger() echo.MiddlewareFunc {
	out := LumberJackLogger(config.AccessLogFilePath+config.AccessLogFileExtension, config.AccessLogMaxSize, config.AccessLogMaxBackups, config.AccessLogMaxAge)
	stdlogger := log.New(out, "", 0)

	config := middleware.DefaultLoggerConfig

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			// Start timer
			start := time.Now()
			// Process request
			next(c)
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			method, statusCode, urlPath, clientIP := GetContextLogInfo(c)

			stdlogger.Printf("[ECHO] %v |%3d| %12v |%s %-7s | %s",
				end.Format("2006/01/02 - 15:04:05"),
				statusCode,
				latency,
				method,
				urlPath,
				clientIP,
			)
			return
		}
	}
}

func colorForStatus(status int) string {
	switch {
	case status >= 200 && status <= 299:
		return green
	case status >= 300 && status <= 399:
		return white
	case status >= 400 && status <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
