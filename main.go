package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/rafaelsq/roar/async"
)

func do(command string, cc chan bool) error {
	cmd := exec.Command(command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		if <-cc {
			_ = cmd.Process.Kill()
		}
	}()

	err = async.Go(func(cancel chan bool) error {
		if err := cmd.Wait(); err != nil {
			return err
		}
		return nil
	}, func(cancel chan bool) error {
		bi := bufio.NewScanner(stdout)
		for {
			if !bi.Scan() {
				break
			}

			fmt.Println(" -", bi.Text())
		}

		if err := bi.Err(); err != nil {
			return err
		}
		return nil
	}, func(cancel chan bool) error {
		bi := bufio.NewScanner(stderr)
		for {
			if !bi.Scan() {
				break
			}

			fmt.Println(" *", bi.Text())
		}

		if err := bi.Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func api(w http.ResponseWriter, r *http.Request) {
	cmds, ok := r.URL.Query()["cmd"]
	if ok {
		var fs []async.TypeFunc

		for _, cmd := range cmds {
			fs = append(fs, func(cancel chan bool) error {
				return do(cmd, cancel)
			})
		}

		err := async.Go(fs...)
		if err != nil {
			fmt.Fprintf(w, "\nerr; %v\n", err)
		} else {
			fmt.Fprintf(w, "\ndone.\n")
		}
	} else {
		fmt.Fprintf(w, "\nno cmd found.\n")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "\nusage: curl http://localhost:8080/api?cmd=./do_build.sh&cmd=./do_too.sh'\n")
}

func main() {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/api", api)
	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
