package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

type WSClient struct {
	ID      string
	Channel string
	Conn    *websocket.Conn
}

var (
	pool map[string]map[string]*WSClient
)

func init() {
	pool = map[string]map[string]*WSClient{}
}

func Add(ID, channel string, conn *websocket.Conn) *WSClient {
	mutex.Lock()
	defer mutex.Unlock()

	c := &WSClient{ID: ID, Channel: channel, Conn: conn}

	if pool[c.Channel] == nil {
		pool[c.Channel] = map[string]*WSClient{}
	}
	pool[c.Channel][c.ID] = c

	return c
}

func Remove(channel, ID string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(pool[channel], ID)
}

func Send(channel string, message *Message) (err error) {
	for ID := range pool[channel] {
		err = pool[channel][ID].Conn.WriteJSON(message)
	}

	return
}
