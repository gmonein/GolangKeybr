package main

import (
	"fmt"
	"log"
	"net/http"
)

func typeWsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("connection upgraded")
	_, content, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		fmt.Println(err)
		return
	}
	log.Println("message read", content)

	parsedToken, err := parseJwtToken(string(content))
	if err != nil {
		conn.WriteMessage(1, []byte("2Auth Error"))
		log.Println(err)
		return
	}
	fmt.Println("token parsed", parsedToken)
	conn.WriteMessage(1, []byte("3Auth success"))
	conn.WriteMessage(1, append([]byte("4"), findCitation()...))

	for {
		fmt.Println("wait")
		_, content, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(content[0]), string(content[1]))
		conn.WriteMessage(1, []byte("1ok"))
	}
}
