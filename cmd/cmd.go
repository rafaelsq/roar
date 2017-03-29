package cmd

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/rafaelsq/roar/async"
)

func Run(ctx context.Context, command string, out chan Msg) error {
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
		<-ctx.Done()
		_ = cmd.Process.Kill()
	}()

	err = async.Run(ctx, func(ctx context.Context) error {
		bi := bufio.NewScanner(stdout)
		for {
			if !bi.Scan() {
				break
			}

			out <- Msg{Text: bi.Text(), Command: command}
		}

		if err := bi.Err(); err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context) error {
		bi := bufio.NewScanner(stderr)
		for {
			if !bi.Scan() {
				break
			}

			out <- Msg{Text: bi.Text(), Type: Error, Command: command}
		}

		if err := bi.Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
