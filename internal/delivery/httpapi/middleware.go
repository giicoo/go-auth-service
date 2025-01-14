package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/giicoo/go-auth-service/internal/jwt"
	"github.com/giicoo/go-auth-service/pkg/apiError"
	"github.com/sirupsen/logrus"
)

type authKey struct{}

func (h *Handler) MiddlewareGetSessionHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.WithFields(
				logrus.Fields{
					"url": r.URL.String(),
				},
			).Errorf("token service: %s", apiError.ErrNotAuthService)
			httpError(w, fmt.Errorf("token service: %w", apiError.ErrNotAuthService))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.WithFields(
				logrus.Fields{
					"url": r.URL.String(),
				},
			).Errorf("token service: %s", apiError.ErrNotAuthService)
			httpError(w, fmt.Errorf("token service: %w", apiError.ErrNotAuthService))
			return
		}

		token := parts[1]
		jwtToken, err := jwt.GetJWTFromEnv(".env")
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"url": r.URL.String(),
				},
			).Errorf("get jwt token .env: %s", apiError.ErrNotAuthService)
			httpError(w, fmt.Errorf("get jwt token .env: %w", apiError.ErrNotAuthService))
		}

		if token != jwtToken {
			logrus.WithFields(
				logrus.Fields{
					"url": r.URL.String(),
				},
			).Errorf("jwt token invalid: %s", apiError.ErrNotAuthService)
			httpError(w, fmt.Errorf("jwt token invalid: %w", apiError.ErrNotAuthService))
		}
		next.ServeHTTP(w, r)
	})
}
