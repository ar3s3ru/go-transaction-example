package http

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"db-transaction-example/internal/app"
)

func listHistory(srv app.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
		if err != nil {
			writeJSON(ctx, w, http.StatusBadRequest, WrapError(err), logger)
			return
		}

		to := time.Now()
		if ts := r.URL.Query().Get("to"); ts != "" {
			if to, err = time.Parse(time.RFC822, ts); err != nil {
				writeJSON(ctx, w, http.StatusBadRequest, WrapError(err), logger)
				return
			}
		}

		var appError app.Error
		history, err := srv.ListHistory(ctx, from, to)

		switch {
		case errors.As(err, &appError):
			writeJSON(ctx, w, appError.StatusCode(), WrapError(err), logger)
		case err != nil:
			writeJSON(ctx, w, http.StatusInternalServerError, WrapError(err), logger)
		default:
			writeJSON(ctx, w, http.StatusCreated, history, logger)
		}
	}
}
