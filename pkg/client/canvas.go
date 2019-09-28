package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/metalblueberry/sicase/internal/actions"
	log "github.com/sirupsen/logrus"
)

type Canvas struct {
	Server  string
	GameID  string
	client  *http.Client
	actions []Action
}

func NewCanvasServer(server, gameID string) *Canvas {
	// initialize http client
	return &Canvas{
		Server:  server,
		GameID:  gameID,
		client:  &http.Client{},
		actions: make([]Action, 0),
	}
}

var defaultCanvas *Canvas

func InitializeDefaultCanvas(server, gameID string) {
	defaultCanvas = NewCanvasServer(server, gameID)
}

// Send connect with the server to draw the current set of instructions in the given turn.
func Send(turn int) { defaultCanvas.Send(turn) }
func (c *Canvas) Send(turn int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	c.SendWithContext(ctx, turn)
}

func (c *Canvas) SendWithContext(ctx context.Context, turn int) error {
	// marshal User to json
	bytes, err := json.Marshal(c.actions)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	req, err := c.buildRequest(ctx, turn, bytes)
	if err != nil {
		log.Printf("error building request for debug server %v", err)
		return err
	}

	err = c.do(req)
	if err != nil {
		log.Printf("error building request for debug server %v", err)
		return err
	}
	c.actions = c.actions[:0]
	return nil
}

func Circle(p Circler, class ...string) { defaultCanvas.Circle(p, class...) }
func (c *Canvas) Circle(p Circler, class ...string) {
	c.actions = append(c.actions, circleAction(p, class...))
}

func circleAction(p Circler, class ...string) Action {
	x, y, r := p.Circle()
	return Action{
		"Method": "Circle",
		"Class":  class,
		"X":      x,
		"Y":      y,
		"R":      r,
	}
}

func Line(p Liner, class ...string) { defaultCanvas.Line(p, class...) }
func (c *Canvas) Line(p Liner, class ...string) {
	c.actions = append(c.actions, lineAction(p, class...))
}

func lineAction(p Liner, class ...string) Action {
	x1, y1, x2, y2 := p.Line()
	return Action{
		"Method": "Line",
		"Class":  class,
		"X1":     x1,
		"Y1":     y1,
		"X2":     x2,
		"Y2":     y2,
	}
}

func (c *Canvas) buildRequest(ctx context.Context, turn int, data []byte) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/img/%s/%d", c.Server, c.GameID, turn), bytes.NewBuffer(data))
}

func (c *Canvas) do(req *http.Request) error {
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
