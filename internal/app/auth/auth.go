package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/utils"
	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

type ContextKey string

const (
	tokenExp             = time.Hour * 3
	secretKey            = "supersecretkey"
	UserIDKey ContextKey = "userID"
)

func WithAuth(h http.Handler) http.HandlerFunc {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		needAuthString := false
		var userID int
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			logger.Log().Debug("No Cookies")
			needAuthString = true
		} else {
			userID = GetUserID(cookie.Value)
			if userID < 0 {
				logger.Log().Debug("Not correct UserId")
				needAuthString = true
			}
		}
		cookie = &http.Cookie{
			Name: "Authorization",
		}

		if needAuthString {
			userID = utils.RandInt()
			token, err := BuildJWTString(userID)
			if err == nil {
				cookie.Value = token
				http.SetCookie(w, cookie)
			}
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		h.ServeHTTP(w, r.WithContext(ctx))

	}

	return authFunc
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(userID int) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}

func CheckAuth(r *http.Request) int {
	token, err := r.Cookie("Authorization")
	var userID int
	if err != nil {
		return -1
	} else {
		userID = GetUserID(token.Value)
		if userID < 0 {
			return -1
		}
	}
	return userID
}
