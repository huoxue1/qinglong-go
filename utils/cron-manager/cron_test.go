package cron_manager

import (
	"fmt"
	"testing"
	"time"
)

func TestAddCron(t *testing.T) {
	_ = AddCron("test1", "*/5 * * * * ?", func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})
	time.Sleep(time.Minute)
	DeleteCron("test1")
	fmt.Println("已停止")
	time.Sleep(10 * time.Second)
}
