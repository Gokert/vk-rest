package middleware

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"vk-rest/pkg/models"
	httpResponse "vk-rest/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userId"

type Core interface {
	GetUserId(ctx context.Context, sid string) (uint64, error)
}

func AuthCheck(next http.Handler, core Core, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		userId, err := core.GetUserId(r.Context(), session.Value)
		if err != nil {
			lg.Error("auth check error", "err", err.Error())
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userId))
		if userId == 0 {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}
