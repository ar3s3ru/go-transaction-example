package http

import (
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"db-transaction-example/internal/app"
)

func addUserHandler(service app.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		ageRaw := r.URL.Query().Get("age")

		if name == "" || ageRaw == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		age, err := strconv.Atoi(ageRaw)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		var appError app.Error
		users, err := service.AddUser(ctx, name, uint8(age))

		switch {
		case errors.As(err, &appError):
			writeJSON(ctx, w, appError.StatusCode(), WrapError(err), logger)
		case err != nil:
			writeJSON(ctx, w, http.StatusInternalServerError, WrapError(err), logger)
		default:
			writeJSON(ctx, w, http.StatusCreated, users, logger)
		}
	}
}
