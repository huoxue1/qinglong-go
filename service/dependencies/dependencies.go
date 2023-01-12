package dependencies

import (
	"bytes"
	"context"
	"fmt"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils"
	"io"
	"strings"
	"time"
)

func AddDep(dep *models.Dependences) {
	if dep.Type == models.NODE {
		addNodeDep(dep)
	} else if dep.Type == models.PYTHON {
		addPythonDep(dep)
	} else {
		addLinuxDep(dep)
	}
}

func addNodeDep(dep *models.Dependences) {
	log := ""
	buffer := bytes.NewBufferString(log)
	ctx := context.WithValue(context.Background(), "cancel", make(chan int, 1))
	now := time.Now()
	utils.RunTask(ctx, fmt.Sprintf("yarn add %s", dep.Name), map[string]string{}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", now.Format("2006-01-02 15:04:05"))))
	}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("\n##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(now).Seconds())))
		dep.Status = 1
		var logs []string
		for _, i2 := range strings.Split(buffer.String(), "\n") {
			logs = append(logs, i2+"\n\n")
		}
		dep.Log = logs
		models.AddDependences(dep)
	}, buffer)
}

func addPythonDep(dep *models.Dependences) {
	log := ""
	buffer := bytes.NewBufferString(log)
	ctx := context.WithValue(context.Background(), "cancel", make(chan int, 1))
	now := time.Now()
	utils.RunTask(ctx, fmt.Sprintf("pip install %s", dep.Name), map[string]string{}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", now.Format("2006-01-02 15:04:05"))))
	}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("\n##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(now).Seconds())))
		dep.Status = 1
		var logs []string
		for _, i2 := range strings.Split(buffer.String(), "\n") {
			logs = append(logs, i2+"\n\n")
		}
		dep.Log = logs
		models.AddDependences(dep)
	}, buffer)
}

func addLinuxDep(dep *models.Dependences) {
	log := ""
	buffer := bytes.NewBufferString(log)
	ctx := context.WithValue(context.Background(), "cancel", make(chan int, 1))
	now := time.Now()
	utils.RunTask(ctx, fmt.Sprintf("apk add %s", dep.Name), map[string]string{}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", now.Format("2006-01-02 15:04:05"))))
	}, func(ctx context.Context) {
		writer := ctx.Value("log").(io.Writer)
		writer.Write([]byte(fmt.Sprintf("\n##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(now).Seconds())))
		dep.Status = 1
		var logs []string
		for _, i2 := range strings.Split(buffer.String(), "\n") {
			logs = append(logs, i2+"\n\n")
		}
		dep.Log = logs
		models.AddDependences(dep)
	}, buffer)
}
