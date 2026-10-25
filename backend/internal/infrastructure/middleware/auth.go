package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"backend/internal/domain/apperror"
	"backend/internal/domain/repository/inputport"
	logger "backend/internal/log"

	"github.com/labstack/echo/v4"
)

// TODO: take the TimeProvider as a parameter, Auth(sessionRepository, timeProvider),
// and pass in the one the router already creates. Replace every time.Now() below with
// timeProvider.Now() so the middleware can be tested.
func Auth(sessionRepository inputport.SessionRepositoryInputPort) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("Middleware: auth")

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// TODO: return new(apperror.Unauthorized), like the other branches; let the error handler write the 401
				return c.JSON(http.StatusUnauthorized, apperror.InvalidCredentials)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				// TODO: return new(apperror.Unauthorized) instead of writing the JSON here
				return c.JSON(http.StatusUnauthorized, apperror.Unauthorized)
			}

			token := parts[1]
			// TODO: move these two lines into util.HashToken(token) and use it here and in
			// TokenGenerator.GenerateToken, so login and the middleware always hash the same way
			hash := sha256.Sum256([]byte(token))
			tokenHash := hex.EncodeToString(hash[:])

			session, err := sessionRepository.GetSessionByTokenHash(c.Request().Context(), tokenHash)
			if err != nil {
				return err
			}

			if session == nil {
				err = new(apperror.Unauthorized)
				return err
			}

			// TODO: expired session -> return new(apperror.Unauthorized)
			if session.ExpiresAt.Before(time.Now()) {
			}

			// TODO: rotation, not deletion. When expires_at < now + 29 days (and the token matched the
			// current token_hash, not the previous one): new token, previous_token_hash = old hash,
			// previous_valid_until = now + 5 min, expires_at = now + sessionLifetime, as one
			// UPDATE ... WHERE token_hash = <old hash>; only if 1 row changed, set the X-New-Token
			// response header. The -1 case below deletes a fresh session, which would log the user out.
			switch time.Now().Add(696 * time.Hour).Compare(session.ExpiresAt) {
			case -1:
			case 0:
			case 1:
			}

			// TODO: put user_id and session_id in the context (logout needs session_id)
			return next(c)
		}
	}
}
