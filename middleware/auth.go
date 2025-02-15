package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"
const UsernameContextKey contextKey = "username"

func ValidateJWTToken(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Validate JWT for request")

		tokenKey := []byte(os.Getenv("TOKEN_KEY"))
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			logger.Warn("Missing token")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "missing token"})
			return
		}

		fieldParts := strings.Split(tokenString, " ")
		if len(fieldParts) != 2 || fieldParts[0] != "Bearer" {
			logger.Warn("Bad token format")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}
		pureToken := fieldParts[1]

		token, err := jwt.Parse(pureToken, func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method.Alg() != "HS256" {
				return nil, errors.New("bad sign method")
			}
			return tokenKey, nil
		})
		if err != nil || !token.Valid {
			logger.Warn("Invalid")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("No user claims")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		claimsUser, ok := claims["user"].(map[string]interface{})
		if !ok {
			logger.Warn("User claim missing")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		userIDString, ok := claimsUser["id"].(string)
		if !ok {
			logger.Warn("Couldn't parse user id")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		userID64, err := strconv.ParseUint(userIDString, 10, 32)
		if err != nil || userID64 == 0 {
			logger.Warn("Encounter user id <= 0")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		userID := uint(userID64)
		username := claimsUser["username"].(string)
		if !ok || len(username) < 3 || len(username) > 50 {
			logger.Warn("Couldn't parse username")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"errors": "bad token"})
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id":  userID,
			"username": username,
		}).Info("User authenticated")

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		ctx = context.WithValue(ctx, UsernameContextKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
