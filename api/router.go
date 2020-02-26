package main

import (
	"fmt"
	"keybr/intraapi"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type handlerWraper func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
type handler func(http.ResponseWriter, *http.Request)

// Routes defined all keybr routes
func Routes() {
	fa := func(h handler) handler {
		return frontRequestWrapper(
			authWrapper(h))
	}
	http.HandleFunc("/citation", fa(citationHandler))
	http.HandleFunc("/data_ws", fa(dataWsHandler))
	http.HandleFunc("/type_ws", fa(typeWsHandler))
	http.HandleFunc("/oauth", oauthHandler)
	http.HandleFunc("/ssh", sshHandler)
}

func parseJwtToken(tokenString string) (jwt.Claims, error) {
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

func authWrapper(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("token")
		cookieToken, _ := r.Cookie("token")

		var parsedToken jwt.Claims
		var err error

		if headerToken != "" {
			parsedToken, err = parseJwtToken(headerToken)
		} else if cookieToken != nil && cookieToken.Value != "" {
			parsedToken, err = parseJwtToken(cookieToken.Value)
		}
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}
		if parsedToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			fmt.Println(parsedToken)
			next(w, r)
		}
	}
}

func frontRequestWrapper(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "localhost:8083")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next(w, r)
	}
}

func citationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	_, err := w.Write(citation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func sshHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDDac87uqRmn8Y+KnvTbmwTdVPhdZkmIwnxJXzksCplgOaQb86m25KBPRvlt8jDMv7OYeuAVgvH8a0I+hyeAmBZWHAuzBxH1UPRR2F4CTsMMBuAqyL8/zw9SlPrG6lacyfSXh6pjtL6kwrAR3ZKDJzT97q4xgWEeM9ZOz7aCmaxqTVhzoD1uuJ0CBNb+q98ZfxhqEg3g6C+5zNldH3ksGYXO/eguajQXRx1zpKBgjfZRTSaJNj1nNUk6Wx4YqYvbs/guEj7vzMr9fof6w2g580yt/dqVWiYqQ5xYpqvCwoONlz41r0ZQ4U8wZ/6v5xf1D6Y55X/nG6yQdfd4K1meGclYjsmLoaWKy+5TeTLc2bAiEHOjVf/vEIfUw/bsu493ZQ6zvCujZdDKM+X/C/gPpd3myxbHp7CpZLVZLrG4U6k3VhqnRULUrJa8FPdithNyRwW6EqUqfeVmy6Xyok32k3lK8AwfH2oBCGGc4GT5L9jPGdd8tS/R2CfY7ipvvTsNzHVDSvowvoXFoDYJaL0FNi/eY+/0H36/k8u24uc11qJ2Diac17TiVxkE0gTgljO5ZSTtzD4pfrZhtJxTxKYlWfcBrDk4W4IYPpx+u7tQHCFe8RHUnY+gJtbqQLsD89lBddR1OaeOCVVDw0uMkjkQfeplcWDikAnHtfCk73yKCHH4w== gmonein@student.42.fr"))
}

func typeWsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	typeReader(ws)
}

func dataWsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	dataWs(ws)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	user, token, err := intraapi.GetUserFromCode(r.URL.Query().Get("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":       user.Login,
		"accessToken": token.AccessToken,
	})
	tokenString, err := jwtToken.SignedString([]byte("toto"))

	fmt.Println(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	cookie := &http.Cookie{Name: "token", Value: tokenString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:8083/login", http.StatusTemporaryRedirect)
}
