package v1

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type claims struct {
	*jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func userHasAnyRoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			rolesFromToken, err := getRolesFromToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			rolesMap := map[string]struct{}{}
			for _, role := range rolesFromToken {
				rolesMap[role] = struct{}{}
			}

			for _, role := range roles {
				if _, ok := rolesMap[role]; ok {
					next.ServeHTTP(w, r)
				}
			}

			w.WriteHeader(http.StatusForbidden)
		})
	}
}

func getRolesFromToken(tokenString string) ([]string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("SECRET"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims.Roles, nil
}
