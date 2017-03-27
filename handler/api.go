package handler

import (
	"fmt"
	"net/http"

	"github.com/rafaelsq/roar/async"
	"github.com/rafaelsq/roar/cmd"
	"github.com/rafaelsq/roar/hub"
)

func API(w http.ResponseWriter, r *http.Request) {
	channel := "all"
	cmds, ok := r.URL.Query()["cmd"]
	if ok {
		var fs []async.TypeFunc

		output := make(chan string)
		done := make(chan struct{}, 1)

		go func() {
			for {
				select {
				case m := <-output:
					hub.Send(channel, &hub.Message{Payload: m})
				case <-done:
					return
				}
			}
		}()

		for _, c := range cmds {
			command := c
			fs = append(fs, func(cancel chan bool) error {
				return cmd.Run(command, cancel, output)
			})
		}

		err := async.Go(fs...)
		done <- struct{}{}
		if err != nil {
			hub.Send(channel, &hub.Message{Type: hub.MessageTypeError, Payload: err.Error()})
			fmt.Fprintf(w, "err; %v\n", err)
		} else {
			hub.Send(channel, &hub.Message{Type: hub.MessageTypeSuccess, Payload: "done without error"})
			fmt.Fprintf(w, "")
		}
	} else {
		fmt.Fprintf(w, "\nno cmd found.\n")
	}
}
