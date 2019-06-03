package service

import (
	"crypto/md5"
	"dao"
	"encoding/hex"
	"fmt"
	"models"
	"strconv"
)

//62位字符表
var elements = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const(
	flag = 0x3fffffff//30位截取
	index = 0x0000003D//0~61相与
)

//将长链接转成短链接
func TranLongToShort(long string)(string, error){
	for{//循环直到找到不碰撞的字符串为止
		//先去缓存中查找
		if ok := checkRedisLong(long); ok{
			return getRedisShort(long), nil
		}
		//缓存没有就去数据库查找
		if ok := dao.CheckLongAddr(long);ok{
			//找到后添加缓存
			short := dao.FindShortAddr(long)
			err := setCache(long, short)
			if err!=nil{
				return short, fmt.Errorf("setCache fail %v", err)
			}
			return short, nil
		}
		tmp, err := handleAddr(long)//不存在就使用MD5算法得到4个字符串
		if err !=nil{//如果有错误
			return "", err
		}
		short, ok := checkValid(tmp)
		if ok{//找到可用字符串
			ist := dao.InsertAddr(models.Addr{LongAddr: long, ShortAddr: short}) //将数据插入数据库
			if !ist{
				continue//插入失败重新处理
			}
			//找到后添加缓存
			err := setCache(long, short)
			if err!=nil{
				return short, fmt.Errorf("setCache fail %v", err)
			}
			return short, nil
		}else{
			long = getMD5(long)//没找到就用MD5加密后的字符串再加密
		}
	}
}

//将长链接通过MD5和缩短处理得到4个字符串
func handleAddr(s string)([4]string,error){
	md := getMD5(s)//得到长地址的32位MD5字符串
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

//生成32位MD5
func getMD5(s string) string{
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}

//检验字符串数组中是否有可用的字符串
func checkValid(s [4]string)(string, bool){
	for _,short:=range s{
		if ok:= dao.CheckShortAddr(short);!ok{
			return short,true
		}
	}
	return "",false
}