package middleware

import (
	"backend/internal/service/user"
	"backend/pkg/chi/response"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string, userRepo user.UserRepository, logger zerolog.Logger) AuthMiddleware {
	return AuthMiddleware{userRepo, jwtauth.New("HS256", []byte(verificationKey), nil), logger}
}

type AuthMiddleware struct {
	userRepo user.UserRepository
	Ja       *jwtauth.JWTAuth
	logger   zerolog.Logger
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func (a AuthMiddleware) HandleToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			tokenString := ""
			bearer := r.Header.Get("Authorization")
			if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
				tokenString = bearer[7:]
			}

			if tokenString == "" {
				err := fmt.Errorf("token empty")
				a.logger.Error().Err(err).Msg("")
				a.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			//do decode token here
			token, err := a.Ja.Decode(tokenString)
			if err != nil {
				a.logger.Error().Err(err).Msg("")
				a.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, _ := token.AsMap(ctx)
			identity := Identity{
				UserUID: claims["user_uid"].(string),
			}

			userEntity, err := a.userRepo.GetByUID(r.Context(), identity.UserUID)
			if err != nil || userEntity == nil {
				a.logger.Error().Err(err).Msg("")
				a.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			identity.Email = userEntity.Email.String
			identity.Role = userEntity.Role.String

			ctx = WithUser(
				r.Context(),
				identity,
			)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func (a AuthMiddleware) writeErrorResponse(w http.ResponseWriter, status int, message string) {
	errResp := response.HTTPError{
		HTTPStatus: status,
		Message:    message,
	}

	jsonBody, err := json.Marshal(errResp)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to marshal error response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBody)
}

type contextKey int

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, userKey, identity)
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) *Identity {
	if userEntity, ok := ctx.Value(userKey).(Identity); ok {
		return &userEntity
	}
	return nil
}
