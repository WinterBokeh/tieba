package Tool

import (
	"github.com/gin-gonic/gin"
)

//检查登录状态，无cookie返回""，否则返回用户名
func CheckLogin(ctx *gin.Context) string {
	redisConn := GetRedisConn()
	cmd := redisConn.Get("isLogin")
	if cmd.Err() != nil {
		return ""
	}

	return cmd.String()
}
