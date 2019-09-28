package main

import (
	"bytes"
	"fmt"
	"sync"
)

var gamesMutex = sync.Mutex{}
var games = make(map[string]*Game)

type Game struct {
	Name string

	turns map[int]*Turn
	m     sync.Mutex
}

func (g *Game) GetTurn(number int) *Turn {
	g.m.Lock()
	defer g.m.Unlock()

	turn, init := g.turns[number]
	if !init {
		turn = NewTurn()
		g.turns[number] = turn
	}
	return turn

}

func NewGame(name string) *Game {
	openbrowser(fmt.Sprintf("http://localhost:8888/visor/%s", name))
	return &Game{
		Name: name,

		turns: make(map[int]*Turn),
	}
}

func GetGame(name string) *Game {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	game, init := games[name]
	if !init {
		game = NewGame(name)
		games[name] = game
	}
	return game
}

// Turn stores the SVG data as a buffer ready to be sent in a web request
type Turn struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *Turn) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}

func (b *Turn) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *Turn) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Bytes()
}

func NewTurn() *Turn {
	return &Turn{
		b: bytes.Buffer{},
		m: sync.Mutex{},
	}
}

func (t *Turn) Svg() []byte {
	buf := &bytes.Buffer{}

	buf.Write(t.Bytes())

	return buf.Bytes()
}
