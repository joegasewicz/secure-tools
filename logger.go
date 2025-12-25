package secure_tools

import (
	"fmt"
	"net/http"
	"time"
)

const (
	RESET = "\033[0m"
	RED   = "\033[31m"
	GREEN = "\033[32m"
	BLUE  = "\033[34m"
)

func PrintWithColor(msg string, color string) string {
	return fmt.Sprintf("%s%s%s", color, msg, RESET)
}

// Logging adds logging for each request
func Logging(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var out string
		start := time.Now()
		// Log response
		duration := time.Duration(time.Now().Sub(start)) * time.Nanosecond

		// Set status
		sw := status_writer.New(w)
		next.ServeHTTP(sw, r)
		statusCode := sw.Status
		msg := fmt.Sprintf("[INFO] %s %s %ds Status: %d\n", r.Method, r.RequestURI, duration, statusCode)

		if statusCode < 400 {
			out = PrintWithColor(msg, BLUE)
		} else {
			out = PrintWithColor(msg, RED)
		}
		fmt.Printf(out)
	})
}
