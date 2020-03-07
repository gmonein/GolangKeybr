package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keybr/intraapi"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var eventQueue EventQueue
var userIDS Serial
var game Game

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1024,
}

func main() {
	eventQueue.Initialize()
	game.Initialize()

	router := gin.Default()
	router.GET("/oauth/intrav2", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{"error": "Empty parameter 'code'"})
			return
		}
		user, _, err := intraapi.GetUserFromCode(code)
		if err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"login": user.Login})
		tokenString, err := jwtToken.SignedString([]byte("toto"))
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	router.GET("/type_ws", func(c *gin.Context) {
		token := c.Query("token")
		user := User{Name: "Guest"}
		if token != "" {
			parsedToken, err := parseJwtToken(token)
			if err != nil {
				c.JSON(400, gin.H{"error": err})
				return
			}
			user.Name = parsedToken["login"].(string)
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		user.ID = userIDS.Next()
		conn.WriteMessage(1, append([]byte("3 - Hello "), []byte(user.Name)...))
		conn.WriteMessage(1, append([]byte("4"), game.Citation...))

		for {
			_, content, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			if game.Citation[user.Index] == content[0] {
				user.Index++
				conn.WriteMessage(1, []byte("1ok"))
				if user.Index == game.CitationLength {
					eventQueue.Push(&Event{
						UserID:    user.ID,
						EventType: Finish,
						NextIndex: user.Index})
				} else {
					eventQueue.Push(&Event{
						UserID:    user.ID,
						EventType: TypingValid,
						NextIndex: user.Index})
				}
			} else {
				eventQueue.Push(&Event{
					UserID:    user.ID,
					EventType: TypingError,
					NextIndex: user.Index})
			}
		}
	})

	router.GET("/data_ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
		s := QueueSubscription{queue: &eventQueue}
		for {
			next := s.Next()
			if next == nil {
				continue
			}
			resp, err := json.Marshal(next)
			if err != nil {
				return
			}
			conn.WriteMessage(1, resp)
		}
	})
}

func findCitation() []byte {
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))
	citationNb := rand.Int() % 100
	dat, err := ioutil.ReadFile(os.Getenv("RESOURCES_PATH") + "/citations/" + strconv.Itoa(citationNb))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return dat
}
