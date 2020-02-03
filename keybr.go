package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Name           string
	CurrentIndex   int
	CurrentError   int
	LastErrorIndex int
	CurrentSpeed   float32
	getConn        *websocket.Conn
}

var users map[string]*User
var currentCitation []byte
var finish bool

func main() {
	// http.HandleFunc("/", handler)
	currentCitation = findCitation()
	finish = false
	http.HandleFunc("/index", templateHandler)
	http.HandleFunc("/citation", citationHandler)
	http.HandleFunc("/data_ws", dataWsHandler)
	http.HandleFunc("/type_ws", typeWsHandler)
	if users == nil {
		users = make(map[string]*User, 10)
	}
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("omg")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func usersJson() []byte {
	jsonString, _ := json.Marshal(users)
	return jsonString
}

func getUser(name string) *User {
	if users[name] == nil {
		users[name] = &User{
			Name:         name,
			CurrentIndex: 0}
	}
	return users[name]
}

func typeReader(conn *websocket.Conn) {
	var message struct {
		Name  string `json:"name"`
		Input string `json:"input"`
	}
	var currentUserName string = ""
	var currentUser *User
	for {
		_, content, err := conn.ReadMessage()
		if err != nil {
			delete(users, currentUserName)
			log.Println(err)
			return
		}
		err = json.Unmarshal(content, &message)
		if err != nil {
			delete(users, currentUserName)
			fmt.Println(err)
			return
		}

		currentUser = getUser(message.Name)
		sendUsers()

		if !finish {
			currentChar := currentCitation[currentUser.CurrentIndex]
			fmt.Println(currentCitation[currentUser.CurrentIndex])
			if message.Input != "" && (currentChar == message.Input[0] || (currentChar == '\n' && message.Input[0] == ' ')) {
				currentUser.CurrentIndex++
			} else {
				if currentUser.LastErrorIndex != currentUser.CurrentIndex {
					currentUser.LastErrorIndex = currentUser.CurrentIndex
					currentUser.CurrentError++
				}
			}
		}
		sendUsers()
		if currentUser.CurrentIndex == len(currentCitation) && !finish {
			go func() {
				time.Sleep(1 * time.Second)
				finish = true
				time.Sleep(1 * time.Second)
				currentCitation = findCitation()
				for _, user := range users {
					user.CurrentIndex = 0
					user.CurrentError = 0
					user.LastErrorIndex = 0
					_ = user.getConn.WriteMessage(1, []byte("new game"))
				}
				sendUsers()
				finish = false
				for _, user := range users {
					_ = user.getConn.WriteMessage(1, []byte("top"))
				}
			}()
		}
	}
}

func sendUser(to *User) {
	if to.getConn != nil {
		err := to.getConn.WriteMessage(1, usersJson())
		if err != nil {
			delete(users, to.Name)
			return
		}
	}
}

func sendUsers() {
	for _, user := range users {
		sendUser(user)
	}
}

func dataWs(conn *websocket.Conn) {
	var message struct {
		Name string `json:"name"`
	}
	_, content, err := conn.ReadMessage()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &message)
	if err != nil {
		return
	}
	currentUser := getUser(message.Name)
	if currentUser.getConn != nil {
		currentUser.getConn.Close()
	}
	currentUser.getConn = conn
	sendUsers()
	for {
		_, _, err = currentUser.getConn.ReadMessage()
		if err != nil {
			currentUser.getConn.Close()
			delete(users, currentUser.Name)
			return
		}
	}
}

func dataWsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	dataWs(ws)
}

func typeWsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	typeReader(ws)
}

func citationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	_, err := w.Write(currentCitation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func findCitation() []byte {
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))
	citation_nb := rand.Int() % 338
	dat, err := ioutil.ReadFile("./citations/" + strconv.Itoa(citation_nb))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return dat
}

//templateHandler renders a template and returns as http response.
func templateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
