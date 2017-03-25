package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/rafaelsq/roar/handler"
)

var port = flag.Int("port", 4000, "Port")

func main() {
	flag.Parse()

	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/api", handler.API)
	http.HandleFunc("/ws", handler.Websocket)
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	fmt.Printf("Listening :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
