package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)
//62位字符表
var elements = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


const(
	flag = 0x3fffffff//30位截取
	index = 0x0000003D//0~61相与
)

//长短链接数据库
var ltos = make(map[string]string)
var stol = make(map[string]string)

var mu sync.RWMutex//操作数据库需要加锁


func main(){
	r := gin.Default()
	r.GET("/trans", func(c *gin.Context) {
		raw := c.Query("addr")//得到查询长网址
		var addr string
		addr, err := transform(raw)//得到长网址对应的短网址
		if err!=nil{//如果出错
			c.String(http.StatusBadRequest, err.Error())//返回400和错误信息
		}
		addr =  "http://localhost:8000/trans/"+ addr
		c.JSON(http.StatusOK, gin.H{"shortAddr":addr,"longAddr":raw})//返回长网址和短网址给客户端
	})
	r.GET("/trans/:addr", func(c *gin.Context) {
		short := c.Param("addr")
		long, ok:= stol[short]
		if !ok{
			c.String(http.StatusNotFound, "addr %s doesn't exit!", short)//返回400和错误信息
		}else{
			c.Redirect(http.StatusMovedPermanently, long)
		}

	})
	r.Run(":8000")
}

//长链接转短链接
func transform(raw string)(string,error){
	var addr string
	mu.RLock()//对数据库加读锁
	cache,ok:=ltos[raw]
	mu.RUnlock()
	if !ok{//检查数据库是否以存在
		for{//循环直到找到不碰撞的字符串为止
			tmp,err := handleAddr(raw)//不存在就使用MD5算法得到4个字符串
			if err !=nil{//如果有错误
				return "", err
			}
			addr, ok = checkValid(tmp)
			if ok{//找到可用字符串
				mu.Lock()//加锁写入
				//fmt.Println(addr)
				ltos[raw] = addr//结果入数据库
				stol[addr] = raw//结果入数据库
				mu.Unlock()
				break
			}else{
				raw = MD5(raw)//没找到就用MD5加密后的字符串再加密
			}
		}
	}else{
		addr = cache//存在返回数据值
	}
	return addr,nil
}


//检验字符串数组中是否有可用的字符串
func checkValid(s [4]string)(string, bool){
	mu.RLock()
	defer mu.RUnlock()
	for _,ele:=range s{
		if _,ok:= stol[ele];!ok{
			return ele,true
		}
	}
	return "",false
}


//生成32位MD5
func MD5(s string) string{
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}

//将长链接通过MD5和缩短处理得到4个字符串
func handleAddr(s string)([4]string,error){
	md := MD5(s)//得到长地址的32位MD5字符串
	var res [4]string//返回4个字符串
	var tempVal int64
	var tempRes []byte
	for i:=0;i<4;i++{//32位每8位分一段
		tempStr := md[8*i:8*i+8]//截取8位字符串
		rawVal,err := strconv.ParseInt(tempStr,16, 64)//将8位字符串转换成16进制的int64类型
		if err!=nil{
			return res, fmt.Errorf("transform %s fail", s)
		}
		tempVal = rawVal&flag//取后30位
		tempRes = []byte{}
		for j:=0;j<6;j++{//每5位取一个字符
			tempRes = append(tempRes, elements[tempVal&index])
			tempVal >>= 5//往右移5位
		}
		res[i] = string(tempRes)
	}
	return res,nil
}