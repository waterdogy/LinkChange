package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service"
	"strings"
)

//长链接转短链接处理器
func LongAddrHandler(c *gin.Context){
	raw := c.DefaultQuery("addr","")//得到查询长网址
	//判断地址是否以http://或者https://开头
	if ok1, ok2:=strings.HasPrefix(raw, "http://"),strings.HasPrefix(raw,"https://");!ok1&&!ok2{
		c.String(http.StatusBadRequest, "Url Format Error!")//返回400和错误信息
		return
	}
	addr, err := service.TranLongToShort(raw)//得到长网址对应的短网址
	if err!=nil{//如果出错
		c.String(http.StatusBadRequest, err.Error())//返回400和错误信息
		return
	}
	addr =  "http://localhost:8000/trans/"+ addr//将地址加上开头
	c.JSON(http.StatusOK, gin.H{"shortAddr":addr,"longAddr":raw})//返回长网址和短网址给客户端
}

//短链接转长链接处理器
func ShortAddrHandler(c *gin.Context){
	short := c.Param("addr")//得到短网址
	long, ok:= service.TranShortToLong(short)//查询数据库是否存有短地址
	if !ok{
		c.String(http.StatusNotFound, "addr %s doesn't exit!", short)//返回400和错误信息
	}else{
		//返回长网址和短网址给客户端
		c.JSON(http.StatusOK, gin.H{"shortAddr":"http://localhost:8000/trans/" + short,"longAddr":long})
	}
}

