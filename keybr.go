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
	"sync"
	"time"
)

type User struct {
	Name           string
	Place          int
	Index          int `json:"index"`
	CurrentError   int
	LastErrorIndex int
	CurrentSpeed   float32
	getConn        *websocket.Conn
	getConnMutex   sync.Mutex
}

var users map[string]*User
var umux sync.Mutex
var citation []byte
var userFinished bool
var finish bool

func main() {
	// http.HandleFunc("/", handler)
	citation = findCitation()
	finish = false
	// http.HandleFunc("/", templateHandler)
	http.HandleFunc("/citation", citationHandler)
	http.HandleFunc("/data_ws", dataWsHandler)
	http.HandleFunc("/type_ws", typeWsHandler)
	http.HandleFunc("/oauth", OauthHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	if users == nil {
		users = make(map[string]*User, 10)
	}
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("omg", err)
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

func deleteUser(name string) {
	umux.Lock()
	delete(users, name)
	umux.Unlock()
}

func sendToEveryUser(message []byte) {
	for _, user := range users {
		user.send(message)
	}
}

func typeSubscribe(conn *websocket.Conn) *User {
	_, content, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return nil
	}

	var message struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(content, &message)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if users[message.Name] != nil {
		return nil
	}
	if len(message.Name) > 8 {
		fmt.Println("Nice try", message.Name)
		return nil
	}
	users[message.Name] = &User{Name: message.Name}
	return users[message.Name]
}

func getInput(conn *websocket.Conn) (byte, error) {
	var message struct {
		Input string `json:"input"`
	}
	_, content, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return '0', err
	}
	fmt.Println("here")
	err = json.Unmarshal(content, &message)
	if err != nil {
		fmt.Println(err)
		return '0', err
	}
	fmt.Println("heree")
	return message.Input[0], nil
}

func (user *User) input(input byte) bool {
	char := citation[user.Index]

	if input != '0' && (char == input || (char == '\n' && input == ' ')) {
		user.Index++
		if len(citation) == user.Index {
			fmt.Println(user.Name, "finished")
			user.finish()
		}
		return true
	} else {
		if user.LastErrorIndex != user.Index {
			user.LastErrorIndex = user.Index
			user.CurrentError++
			return true
		}
	}
	return false
}

func typeReader(conn *websocket.Conn) {
	user := typeSubscribe(conn)
	if user == nil {
		conn.Close()
		return
	}
	sendUsers()
	for {
		fmt.Println("wait")
		input, err := getInput(conn)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		user.input(input)
		fmt.Println("wait 1")
		sendUsers()
		fmt.Println("wait 2")
	}
}

func newGame() {
	time.Sleep(3 * time.Second)
	finish = true
	time.Sleep(3 * time.Second)
	citation = findCitation()
	for _, user := range users {
		user.reset()
		user.send([]byte("new game"))
	}
	sendUsers()
	sendToEveryUser([]byte("Wait"))
	time.Sleep(3 * time.Second)
	finish = false
	sendToEveryUser([]byte("Go"))
}

func (user *User) finish() {
	if !userFinished {
		user.Place = 1
		go newGame()
	}
}

func (user *User) reset() {
	user.Index = 0
	user.CurrentError = 0
	user.LastErrorIndex = 0
	user.Place = -1
}

func (user *User) send(message []byte) {
	if user.getConn == nil {
		return
	}
	user.getConnMutex.Lock()
	err := user.getConn.WriteMessage(1, message)
	user.getConnMutex.Unlock()
	if err != nil {
		fmt.Println(err)
		deleteUser(user.Name)
	}
}

func sendUsers() {
	m := usersJson()
	sendToEveryUser(m)
}

func dataWs(conn *websocket.Conn) {
	var message struct {
		Name string `json:"name"`
	}
	_, content, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data ws receiving connection")
	err = json.Unmarshal(content, &message)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := users[message.Name]
	if user == nil {
		fmt.Println(message.Name, ": data ws user doesn't exist")
		conn.Close()
		return
	}
	if user.getConn != nil {
		conn.Close()
		return
	}
	user.getConn = conn
	sendUsers()
	fmt.Println(user.Name, ": conn succeed sending users")
	for {
		_, _, err = user.getConn.ReadMessage()
		if err != nil {
			fmt.Println(user.Name, ": crashed")
			fmt.Println(err)
			user.getConn.Close()
			deleteUser(user.Name)
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
	_, err := w.Write(citation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func findCitation() []byte {
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))
	citation_nb := rand.Int() % 100
	dat, err := ioutil.ReadFile("./citations/" + strconv.Itoa(citation_nb))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return dat
}

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
