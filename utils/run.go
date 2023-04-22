package utils

import (
	"context"
	"fmt"
	log "github.com/huoxue1/go-utils/base/log"
	"github.com/huoxue1/qinglong-go/service/config"
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
	environ := os.Environ()
	dir, _ := os.Getwd()
	port := config.ListenPort()
	cmd.Env = append(append(cmd.Env, environ...), "QL_DIR="+dir, fmt.Sprintf("QL_PORT=%d", port))
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

type RunOption struct {
	Ctx     context.Context
	Command string
	Env     map[string]string
	OnStart func(ctx context.Context)
	OnEnd   func(ctx context.Context)
	LogFile io.Writer
	CmdDir  string
}

func RunWithOption(ctx context.Context, option *RunOption) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("执行command出现异常")
			log.Errorln(err)
		}
	}()
	cmd := exec.Command(strings.Split(option.Command, " ")[0], strings.Split(option.Command, " ")[1:]...)
	for s, s2 := range option.Env {
		cmd.Env = append(cmd.Env, s+"="+s2)
	}
	environ := os.Environ()
	dir, _ := os.Getwd()
	port := config.ListenPort()
	cmd.Env = append(append(cmd.Env, environ...), "QL_DIR="+dir, fmt.Sprintf("QL_PORT=%d", port))
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()
	cmd.Dir = option.CmdDir
	option.OnStart(context.WithValue(ctx, "log", option.LogFile))
	ch := make(chan int, 1)
	go func() {
		err := cmd.Start()
		if err != nil {
			ch <- 1
			return
		}
		go io.Copy(option.LogFile, stderrPipe)
		go io.Copy(option.LogFile, stdoutPipe)
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
			option.OnEnd(context.WithValue(ctx, "log", option.LogFile))
		}
	case <-cancel:
		{
			_ = cmd.Process.Kill()
			option.OnEnd(context.WithValue(context.Background(), "log", option.LogFile))
		}

	}
}
