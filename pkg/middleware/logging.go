package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		
		// Buffer the request body
		var requestBody []byte
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			requestBody, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
		}
		
		// Create a response writer that captures the response
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer
		
		// Process request
		c.Next()
		
		// Calculate request duration
		duration := time.Since(start)
		
		// Log request details
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()
		
		if query != "" {
			path = path + "?" + query
		}
		
		// Format the log entry
		logEntry := fmt.Sprintf("[%s] %s | %3d | %13v | %15s | %s | %s",
			method,
			path,
			statusCode,
			duration,
			clientIP,
			userAgent,
			requestBody,
		)
		
		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Printf("[ERROR] %s\n", logEntry)
		case statusCode >= 400:
			logger.Printf("[WARN] %s\n", logEntry)
		default:
			logger.Printf("[INFO] %s\n", logEntry)
		}
	}
}

// responseBodyWriter is a custom ResponseWriter that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body while writing it to the client
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString writes a string to the response and captures it
func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
