package cmd

import (
	"bufio"
	"os/exec"

	"github.com/rafaelsq/roar/async"
)

func Run(command string, cc chan bool, out chan string) error {
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

			out <- bi.Text()
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

			out <- bi.Text()
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
