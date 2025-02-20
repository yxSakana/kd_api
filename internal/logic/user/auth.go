package user

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	StudentId uint32
	jwt.RegisteredClaims
}

// var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var jwtSecret = []byte("test")

func GenerateToken(studentId uint32) (tokenStr string, err error) {
	claims := Claims{
		StudentId: studentId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", gerror.New(err.Error())
	}
	return
}

func ParseToken(tokenStr string) (claims *Claims, err error) {
	claims = &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	if !token.Valid {
		return nil, gerror.New("invalid token")
	}
	return
}
