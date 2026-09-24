package handler

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"pos-backend/internal/bootstrap"
	"pos-backend/internal/config"
)

var (
	appInstance *bootstrap.App
	initOnce    sync.Once
	initErr     error
)

// Handler adalah entrypoint utama untuk Vercel Serverless Function (@vercel/go).
// Semua request (/health, /api/v1/*) akan diteruskan ke Gin Router.
func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			initErr = err
			return
		}
		appInstance, initErr = bootstrap.New(context.Background(), cfg)
	})

	if initErr != nil {
		http.Error(w, "Backend initialization failed: "+initErr.Error(), http.StatusInternalServerError)
		return
	}

	// Pulihkan path asli dari rewrite vercel.json (?__path=/...)
	if origPath := r.URL.Query().Get("__path"); origPath != "" {
		if !strings.HasPrefix(origPath, "/") {
			origPath = "/" + origPath
		}
		r.URL.Path = origPath
		r.URL.RawPath = origPath
	}

	appInstance.Router.ServeHTTP(w, r)
}
