package main

import (
	"github.com/broqiang/mdblog/app"
	"github.com/broqiang/mdblog/app/config/mylog"
)

func main() {
	mylog.LogInfo.Println("Hello")
	engine := app.Init()

	engine.Run(":8080")
}
