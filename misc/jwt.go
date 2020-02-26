package main

// import (
// 	"github.com/dgrijalva/jwt-go"
// 	"time"
// )

// var jwtKey = []byte("my_secret_key")

// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.StandardClaims
// }

// type JwtToken struct {
// }

// // Create the Signin handler
// func generateJWT(login string) map[string]string {
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	claims := &Claims{
// 		Username: login,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Value:   tokenString,
// 		Expires: expirationTime,
// 	})
// }
