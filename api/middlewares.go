package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxUserKey struct{}

type AuthMiddleware struct {
	us UserService
}

func NewAuthMiddleware(us UserService) *AuthMiddleware {
	return &AuthMiddleware{
		us: us,
	}
}

func (m AuthMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sendJsonError(w, err, http.StatusNotFound)
			return
		}

		cookieUUID, err := uuid.Parse(cookie.Value)
		if err != nil {
			sendJsonError(w, err, http.StatusInternalServerError)
			return
		}

		//найти сессию в БД(по ID из cookie)
		session, err := m.us.FindSessionByID(cookieUUID)
		if err != nil {
			sendJsonError(w, err, http.StatusInternalServerError)
			return
		}
		// найти пользователя в базе по userid из сессии
		user, err := m.us.GetUser(session.UserID)
		if err != nil {
			sendJsonError(w, err, http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserKey{}, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

}
