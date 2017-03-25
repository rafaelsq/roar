package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/rafaelsq/roar/hub"
)

var (
	upgrader = websocket.Upgrader{}
)

func Websocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Println(err)
		}
	}()

	channel := "all"
	client := hub.Add(uuid.NewV4().String(), channel, ws)

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			hub.Remove(channel, client.ID)
			break
		}
	}
}
