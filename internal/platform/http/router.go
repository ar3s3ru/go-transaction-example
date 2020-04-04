package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"go-transaction-example/internal/app"
)

func Router(srv app.Service, logger *zap.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(1 * time.Second))

	router.Route("/users", func(r chi.Router) {
		r.Post("/", addUserHandler(srv, logger))
		r.Get("/history", listHistory(srv, logger))
	})

	return router
}

func writeJSON(ctx context.Context, w http.ResponseWriter, code int, target interface{}, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(target); err != nil {
		writeError(ctx, w, http.StatusInternalServerError, err, logger)
	}
}

func writeError(ctx context.Context, w http.ResponseWriter, code int, err error, logger *zap.Logger) {
	w.WriteHeader(code)

	if _, err = io.Copy(w, bytes.NewBufferString(err.Error())); err != nil {
		id := middleware.GetReqID(ctx)

		logger.Error("Failed to write error to HTTP response",
			zap.String("request_id", id),
			zap.Error(err),
		)
	}
}
