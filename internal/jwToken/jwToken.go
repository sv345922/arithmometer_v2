package jwToken

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// expireTime - время жизни токена в часах
func CreateJWT(userName string, expireTime int) (string, error) {
	hmacSampleSecret := configs.SecretString
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userName,
		"nbf":  now.Unix(),
		"exp":  now.Add(time.Duration(expireTime) * time.Hour).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}
	fmt.Println("token string:", tokenString) // todo delete

	return tokenString, nil
}
func CheckJWT(tokenString string) (string, bool, error) {
	hmacSampleSecret := configs.SecretString

	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
	}
	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		userName := claims["name"].(string)
		return userName, true, nil
	}
	return "", false, err
}
