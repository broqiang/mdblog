package main

import (
	"github.com/broqiang/mdblog/app"
	"github.com/broqiang/mdblog/app/mylog"
)

func main() {
	mylog.LogInfo.Println("Hello")
	engine := app.Init()

	app.Run(engine)
}
