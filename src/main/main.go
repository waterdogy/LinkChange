package main

import (
	"github.com/gin-gonic/gin"
	"handler"
)

func main(){
	r := gin.Default()
	//处理长链接转短链接
	r.GET("/trans", handler.LongAddrHandler)

	//处理短链接转长链接
	r.GET("/trans/:addr", handler.ShortAddrHandler)

	//监听8000端口启动服务
	r.Run(":8000")
}