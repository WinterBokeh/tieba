package Controller

import "github.com/gin-gonic/gin"

type UserController struct {

}

func (u *UserController) Router(engine *gin.Engine) {
	engine.POST("/register", u.register)
}

func (u *UserController) register(ctx *gin.Context) {

}
