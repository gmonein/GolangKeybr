package main

type Game struct {
	Citation        []byte
	CitationLength  int
	someoneFinished bool
	Winners         []*User
}

func (g *Game) Initialize() {
	g.Citation = findCitation()
	g.CitationLength = len(g.Citation)
	g.someoneFinished = false
	g.Winners = []*User{}
}
