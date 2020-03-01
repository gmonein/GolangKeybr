package main

import (
	"encoding/json"
	"fmt"
	"keybr/intraapi"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type context struct {
	CurrentToken *jwt.Claims
}

type handler func(http.ResponseWriter, *http.Request)

type handlerWithContext func(context, http.ResponseWriter, *http.Request)

func contextWrapper(hwc handlerWithContext) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		hwc(context{}, w, r)
	}
}

// Routes defined all keybr routes
func Routes() {
	fa := func(h handler) handler {
		return frontRequestWrapper(
			authWrapper(h))
	}
	http.HandleFunc("/citation", fa(citationHandler))
	http.HandleFunc("/data_ws", fa(dataWsHandler))
	http.HandleFunc("/type_ws", fa(typeWsHandler))
	http.HandleFunc("/oauth", frontRequestWrapper(oauthHandler))
	http.HandleFunc("/ssh", sshHandler)
	http.HandleFunc("/whoami", fa(whoamiHandler))
	http.HandleFunc("/logout", frontRequestWrapper(logoutHandler))
}

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	token := findTokenFromRequest(r)
	if token == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	parsedToken, err := parseJwtToken(*token)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if parsedToken["login"] == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	respMap := map[string]string{"login": parsedToken["login"].(string)}
	respJSON, err := json.Marshal(respMap)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = w.Write(respJSON)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(respJSON)
}

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

func findTokenFromRequest(r *http.Request) *string {
	headerToken := r.Header.Get("token")
	if headerToken != "" {
		return &headerToken
	}
	cookieToken, _ := r.Cookie("token")
	if cookieToken != nil {
		return &cookieToken.Value
	}
	return nil
}

func authWrapper(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := findTokenFromRequest(r)

		fmt.Println(token)
		if token != nil {
			parsedToken, err := parseJwtToken(*token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Println(err)
				return
			}
			fmt.Println(parsedToken)
			next(w, r)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func frontRequestWrapper(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8083")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintln(err)))
		return
	}
	resp := map[string]string{"token": tokenString}
	respJSON, _ := json.Marshal(resp)
	w.Write([]byte(respJSON))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	// http.Redirect(w, r, "http://localhost:8083/login", http.StatusAccepted)
}
