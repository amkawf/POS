package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/pkg/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestID(t *testing.T) {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.GET("/ping", func(c *gin.Context) {
		reqID := middleware.GetRequestID(c)
		c.String(http.StatusOK, reqID)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	headerID := w.Header().Get(middleware.HeaderXRequestID)
	if headerID == "" {
		t.Fatalf("expected X-Request-ID header to be present")
	}
	if w.Body.String() != headerID {
		t.Fatalf("expected context request ID to match header: %s vs %s", w.Body.String(), headerID)
	}
}

func TestCORS(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS())
	r.GET("/api/v1/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/api/v1/ping", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("expected origin header to be reflected")
	}
}

func TestRecovery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.GET("/panic", func(c *gin.Context) {
		panic("something went critically wrong!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", w.Code)
	}
}
