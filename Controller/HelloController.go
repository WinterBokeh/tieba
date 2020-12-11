package Controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HelloController struct {

}

func (h *HelloController) Router(engine *gin.Engine)  {
	engine.GET("/hello", h.hello)
}

//测试
func (h *HelloController) hello(ctx *gin.Context) {
	fmt.Println("HelloTieba")
	ctx.String(200, "HelloTieba")
}
