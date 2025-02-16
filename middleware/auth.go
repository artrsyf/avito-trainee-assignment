package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	JSONResponse "github.com/artrsyf/avito-trainee-assignment/pkg/json_response"
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
			JSONResponse.JSONResponse(
				w,
				http.StatusUnauthorized,
				map[string]string{"errors": "missing token"},
			)
			return
		}

		fieldParts := strings.Split(tokenString, " ")
		if len(fieldParts) != 2 || fieldParts[0] != "Bearer" {
			sendBadTokenError(w, logger, "Bad token format")
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
			sendBadTokenError(w, logger, "Token is invalid")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendBadTokenError(w, logger, "User claims type cast err")
			return
		}

		claimsUser, ok := claims["user"].(map[string]interface{})
		if !ok {
			sendBadTokenError(w, logger, "User claims are missing")
			return
		}

		userIDString, ok := claimsUser["id"].(string)
		if !ok {
			sendBadTokenError(w, logger, "Couldn't parse user id")
			return
		}

		userID64, err := strconv.ParseUint(userIDString, 10, 32)
		if err != nil || userID64 == 0 {
			sendBadTokenError(w, logger, "Encounter user id <= 0")
			return
		}

		userID := uint(userID64)
		username, ok := claimsUser["username"].(string)
		if !ok {
			sendBadTokenError(w, logger, "Couldn't cast username to string")
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

func sendBadTokenError(w http.ResponseWriter, logger *logrus.Logger, msg string) {
	logger.Warn(msg)
	JSONResponse.JSONResponse(w, http.StatusUnauthorized,
		map[string]string{"errors": "bad token"})
}
