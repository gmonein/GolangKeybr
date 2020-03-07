package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func parseJwtToken(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, nil
	}
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("toto"), nil
		})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	fmt.Println(err)
	return nil, err
}
