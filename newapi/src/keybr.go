package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"keybr/intraapi"
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
		commandTracker := Tracker{}
		commandTracker.Add([]byte("restart"), func() {
			game.Finished()
		})
		for {
			conn.WriteMessage(1, append([]byte("4"), game.Citation...))

			currentGameID := game.ID
			user.Index = 0
			eventQueue.Push(&Event{
				UserID:    user.ID,
				EventType: TypingValid,
				NextIndex: user.Index})
			for game.IsStarging() {
			}
			conn.WriteMessage(1, append([]byte("5"), []byte("GO")...))
			go func(gameID int) {
				for game.IsOnGoing() {
					_, content, err := conn.ReadMessage()
					if err != nil {
						return
					}
					commandTracker.Push(content[0])
					if !game.IsOnGoing() || gameID != game.ID {
						break
					}
					if content[0] == '2' {
						continue
					}
					if game.CitationLength == 0 ||
						game.Citation[user.Index] == content[0] ||
						(game.Citation[user.Index] == '\n' && content[0] == ' ') {
						user.Index++
						if user.Index == game.CitationLength {
							eventQueue.Push(&Event{
								UserID:    user.ID,
								EventType: Finish,
								NextIndex: user.Index})
							game.Finished()
							break
						} else {
							eventQueue.Push(&Event{
								UserID:    user.ID,
								EventType: TypingValid,
								NextIndex: user.Index})
						}
						conn.WriteMessage(1, []byte("1ok"))

					} else {
						eventQueue.Push(&Event{
							UserID:    user.ID,
							EventType: TypingError,
							NextIndex: user.Index})
						conn.WriteMessage(1, []byte("2nop"))
					}
				}
			}(currentGameID)
			for game.IsOnGoing() {
				time.Sleep(1)
			}
			for game.IsFinished() {
				time.Sleep(1)
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
				time.Sleep(30 * time.Millisecond)
				continue
			}
			resp, err := json.Marshal(next)
			if err != nil {
				return
			}
			if err = conn.WriteMessage(1, resp); err != nil {
				return
			}
		}
	})
	router.Run(":8084")
}

func findCitation() []byte {
	rand.Seed(time.Now().UnixNano() / int64(time.Millisecond))
	citationNb := rand.Int() % 100
	fmt.Println(citationNb)
	dat, err := ioutil.ReadFile(os.Getenv("RESOURCES_PATH") + "/citations/" + strconv.Itoa(citationNb))
	if err != nil {
		return []byte(err.Error())
	}
	return dat
}
