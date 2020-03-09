package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

// User is pretty cool
type User struct {
	Index int
}

var users map[int]*User
var usersCount int

var addr = flag.String("addr", "localhost:8084", "http service address")
var citation []byte
var citationLen int
var sanitizedCitation []byte
var index = 0
var onGoing bool
var debugg string

func refreshCitation() {
}

func grey(s string) string {
	return "\033[0m" + s + "\033[0m"
}
func underline(s string) string {
	return "\033[4m" + s + "\033[0m"
}
func bold(s string) string {
	return "\033[1;14m" + s + "\033[0m"
}
func red(s string) string {
	return "\033[1;34m" + s + "\033[0m"
}
func underlineBold(s string) string {
	return "\033[1;4m" + s + "\033[0m"
}
func redUnderlineBold(s string) string {
	return "\033[1m\033[31m\033[4m" + s + "\033[0m"
}
func highlight(s string) string {
	return "\033[7m;" + s + "\033[0m"
}

func sanitizeLineReturn(s []byte) []byte {
	return bytes.Replace(s, []byte("\n"), []byte(" "), 4)
}

var rand int = 0
var nextRefresh time.Time
var noRefresh = true

func refresh() {
	if noRefresh {
		return
	}
	if nextRefresh.After(time.Now()) {
		return
	}
	nextRefresh = time.Now().Add(40 * time.Millisecond)
	clearString := fmt.Sprintf("\033[1J")

	var keys []int
	for k := range users {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	usersString := ""
	// To perform the opertion you want
	for _, k := range keys {
		user := users[k]
		if user != nil {
			usersString += fmt.Sprintf("%d\t|", k)
			if user.Index < citationLen {
				usersString += fmt.Sprintf(grey("%s")+"|"+"\n", sanitizedCitation[0:user.Index])
			} else if user.Index == citationLen {
				usersString += fmt.Sprintf("%s"+"\n", sanitizedCitation)
			}
		}
	}
	usersString += fmt.Sprintf("\n---\n\n")

	citationString := ""
	if index < citationLen-1 {
		if onGoing {
			citationString += fmt.Sprintf(bold("%s")+underlineBold("%c")+grey("%s")+"\n", citation[0:index], citation[index], citation[index+1:])
		} else {
			citationString += fmt.Sprintf(bold("%s")+"\n", citation)
		}
	} else if index == citationLen {
		citationString += fmt.Sprintf(bold("%s")+"\n", citation[0:index])
	} else if index != 0 {
		debugg = string(index) + "-" + string(citationLen)
		citationString += fmt.Sprintf(bold("%s")+underlineBold("%c")+"\n", citation[0:index], citation[index])
	}
	fmt.Printf("%s%s%s\n", clearString, usersString, citationString)
}

func main() {
	users = make(map[int]*User, 100)
	usersCount = 0
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/type_ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if message[0] == '1' {
				index++
				if index < len(citation) {
					refresh()
				}
			}
			if message[0] == '2' {
				refresh()
			}
			if message[0] == '4' {
				citation = message[1:]
				citationLen = len(citation)
				sanitizedCitation = sanitizeLineReturn(citation)
				index = 0
				refresh()
				onGoing = false
			}
			if message[0] == '5' {
				index = 0
				onGoing = true
				refresh()
			}
			debugg = ""
		}
	}()

	uData := url.URL{Scheme: "ws", Host: *addr, Path: "/data_ws"}
	log.Printf("connecting to %s", uData.String())

	cData, _, err := websocket.DefaultDialer.Dial(uData.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer cData.Close()

	doneData := make(chan struct{})

	go func() {
		defer close(doneData)
		for {
			var JSONmessage struct {
				UserID    int    `json:"user_id"`
				Event     string `json:"event"`
				NextIndex int    `json:"next_index"`
			}
			_, message, err := cData.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			cData.WriteMessage(websocket.TextMessage, []byte(" "))
			json.Unmarshal(message, &JSONmessage)
			id := JSONmessage.UserID
			if users[id] == nil {
				users[id] = &User{}
			}
			users[id].Index = JSONmessage.NextIndex
			refresh()
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		c.WriteMessage(websocket.TextMessage, []byte("2"))
	}()
	for {
		if onGoing {
			for _, ch := range citation {
				err := c.WriteMessage(websocket.TextMessage, []byte(string(ch)))
				time.Sleep(10 * time.Millisecond)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}
}
