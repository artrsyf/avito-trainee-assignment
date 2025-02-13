package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"
const UsernameContextKey contextKey = "username"

func ValidateJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenKey := []byte(os.Getenv("TOKEN_KEY"))
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			fmt.Println("missing token")
			return
		}

		fieldParts := strings.Split(tokenString, " ")
		if len(fieldParts) != 2 || fieldParts[0] != "Bearer" {
			fmt.Println("bad token format")
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
			fmt.Println("bad token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("no payload")
			return
		}

		claimsUser := claims["user"].(map[string]interface{})
		userIDString := claimsUser["id"].(string)
		userID64, err := strconv.ParseUint(userIDString, 10, 32)
		if err != nil {
			fmt.Println("type cast error")
			return
		}

		userID := uint(userID64)
		username := claimsUser["username"].(string)

		// _, err = sessionRepo.Check(userID) /*TODO MB*/
		// if err != nil {
		// 	fmt.Println("no session")
		// }

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		ctx = context.WithValue(ctx, UsernameContextKey, username)
		fmt.Println("Username: ", username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
