package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rafaelsq/roar/async"
	"github.com/rafaelsq/roar/cmd"
	"github.com/rafaelsq/roar/hub"
)

func API(w http.ResponseWriter, r *http.Request) {
	channel := "all"
	rawId := r.URL.Query()["id"]
	var Id int
	if len(rawId) > 0 {
		Id, _ = strconv.Atoi(rawId[0])
	}

	cmds, ok := r.URL.Query()["cmd"]
	if ok {
		go hub.Send(channel, &hub.Message{
			Type: hub.MessageTypeNewChannel,
			Payload: map[string]interface{}{
				"Id":       Id,
				"Commands": cmds,
			},
		})

		var fs []async.TypeFunc

		output := make(chan cmd.Msg)
		done := make(chan struct{}, 1)

		go func() {
			for {
				select {
				case m := <-output:
					hub.Send(channel, &hub.Message{Payload: map[string]interface{}{
						"Id":      Id,
						"Text":    m.Text,
						"Type":    m.Type,
						"Command": m.Command,
					}})
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
			hub.Send(channel, &hub.Message{
				Type: hub.MessageTypeError,
				Payload: map[string]interface{}{
					"Id":   Id,
					"Text": err.Error(),
				},
			})
			fmt.Fprintf(w, "err; %v\n", err)
		} else {
			hub.Send(channel, &hub.Message{
				Type: hub.MessageTypeSuccess,
				Payload: map[string]interface{}{
					"Id":   Id,
					"Text": "done without error",
				},
			})
			fmt.Fprintf(w, "")
		}
	} else {
		fmt.Fprintf(w, "\nno cmd found.\n")
	}
}
