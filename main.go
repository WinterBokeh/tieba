package main

import (
	"Tieba/Controller"
	"Tieba/Tool"
	"github.com/gin-gonic/gin"
)

func main() {
	err := Tool.ParseCfg("./Config/app.json")
	if err != nil {
		panic(err)
	}

	engine := gin.Default()
	registeRouter(engine)

	engine.Run()
}

func registeRouter(engine *gin.Engine) {
	new(Controller.HelloController).Router(engine)
}