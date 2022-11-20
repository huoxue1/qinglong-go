package utils

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Context struct {
	process *os.Process
}

func RunTask(ctx context.Context, command string, env map[string]string, onStart func(ctx context.Context), onEnd func(ctx context.Context), logFile io.Writer) {
	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	for s, s2 := range env {
		cmd.Env = append(cmd.Env, s+"="+s2)
	}
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	cmd.Dir = "./data/scripts/"
	onStart(context.WithValue(ctx, "log", logFile))
	ch := make(chan int, 1)
	go func() {
		err := cmd.Start()
		if err != nil {
			ch <- 1
			return
		}
		go io.Copy(logFile, stderrPipe)
		go io.Copy(logFile, stdoutPipe)
		err = cmd.Wait()
		if err != nil {
			ch <- 1
			return
		}
		ch <- 1
	}()
	cancel := ctx.Value("cancel").(chan int)
	select {
	case <-ch:
		{
			onEnd(context.WithValue(ctx, "log", logFile))
		}
	case <-cancel:
		{
			_ = cmd.Process.Kill()
			onEnd(context.WithValue(context.Background(), "log", logFile))
		}

	}
}
