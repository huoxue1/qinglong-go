package main

import (
	"github.com/huoxue1/qinglong-go/controller"
)

func main() {
	engine := controller.Router()
	_ = engine.Run(":8080")
}
