package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	jwt.RegisteredClaims
	UserID int
}

const (
	tokenExp  = time.Hour * 3
	secretKey = "supersecretkey"
	UserID    = 17
)

func WithAuth(h http.Handler) http.HandlerFunc {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		needAuthString := false
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			logger.Log().Debug("No Cookies")
			needAuthString = true
		} else {
			if id := GetUserID(cookie.Value); id < 0 {
				logger.Log().Debug("Not correct UserId")
				needAuthString = true
			}
		}
		cookie = &http.Cookie{
			Name: "Authorization",
		}
		if needAuthString {
			token, err := BuildJWTString(UserID)
			if err == nil {
				cookie.Value = token
				http.SetCookie(w, cookie)
			}
		}

		h.ServeHTTP(w, r)

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
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token os valid")
	return claims.UserID
}

func CheckAuth(r *http.Request) bool {
	token, err := r.Cookie("Authorization")
	if err != nil {
		return false
	} else {
		if GetUserID(token.Value) < 0 {
			return false
		}
	}
	return true
}
