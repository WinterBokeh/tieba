package main

import (
	"Tieba/Controller"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	registerRouter(engine)

	engine.Run()
}

func registerRouter(engine *gin.Engine) {
	new(Controller.HelloController).Router(engine)
	new(Controller.UserController).Router(engine)
}