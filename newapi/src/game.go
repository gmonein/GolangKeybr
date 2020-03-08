package main

import (
	"fmt"
	"time"
)

var stateStarting = 1
var stateOnGoing = 2
var stateFinished = 3

var gameIDS Serial

type Game struct {
	Citation        []byte
	CitationLength  int
	State           int
	someoneFinished bool
	Winners         []*User
	ID              int
}

func (g *Game) IsOnGoing() bool {
	return g.State == stateOnGoing
}
func (g *Game) IsFinished() bool {
	return g.State == stateFinished
}
func (g *Game) IsStarging() bool {
	return g.State == stateStarting
}
func (g *Game) Finished() {
	if g.someoneFinished {
		return
	}
	g.someoneFinished = true
	go func(g *Game) {
		time.Sleep(3 * time.Second)
		fmt.Println("set finished")
		game.State = stateFinished
		go func(g *Game) {
			time.Sleep(3 * time.Second)
			fmt.Println("initialize new game")
			game.Initialize()
		}(g)
	}(g)
}

func (g *Game) Initialize() {
	g.Citation = findCitation()
	g.CitationLength = len(g.Citation)
	g.someoneFinished = false
	g.Winners = []*User{}
	g.ID = gameIDS.Next()
	g.State = stateStarting
	go func(g *Game) {
		time.Sleep(3 * time.Second)
		fmt.Println("GO !")
		g.State = stateOnGoing
	}(g)
}
