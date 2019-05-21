package dao

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
	"utils"
)

var redisClient *redis.Pool
func init() {
	fmt.Println("初始化redis连接池")
	iniParser := utils.IniParser{}
	if err:=iniParser.Load(confFilename); err!=nil{
		fmt.Println("loading ini file fail",err)
		return
	}
	//最大空闲数
	mi := iniParser.GetInt64("redis","MaxIdle")
	//最大活跃数
	ma := iniParser.GetInt64("redis","MaxActive")
	//闲置时间
	it := iniParser.GetInt64("redis","IdleTimeout")
	//主机名
	host := iniParser.GetString("redis","Host")
	redisClient = &redis.Pool{
		MaxActive: int(ma),
		MaxIdle: int(mi),
		IdleTimeout: time.Second * time.Duration(it),
		Wait: true,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", host)
		},
	}
}
//在redis缓存中判断是否存在短地址
func checkRedisShort(short string) bool{
	rc := redisClient.Get()
	defer rc.Close()
	ok, err:= redis.Bool(rc.Do("EXISTS",short))
	if err!=nil{
		fmt.Println("redis exists error", err)
		return false
	}
	return ok
}
//在redis缓存中判断是否存在长地址
func checkRedisLong(long string) bool{
	rc := redisClient.Get()
	defer rc.Close()
	ok, err:= redis.Bool(rc.Do("EXISTS",long))
	if err!=nil{
		fmt.Println("redis exists error", err)
		return false
	}
	return ok
}
//在redis缓存中查找长地址
func getRedisLong(short string) string{
	rc := redisClient.Get()
	defer rc.Close()
	longArr, err:= redis.String(rc.Do("GET",short))
	if err!=nil{
		return ""
	}else{
		return longArr
	}
}

//在redis缓存中查找短地址
func getRedisShort(long string) string{
	rc := redisClient.Get()
	defer rc.Close()
	shortArr, err:= redis.String(rc.Do("GET",long))
	if err!=nil{
		return ""
	}else{
		return shortArr
	}
}
//在redis缓存中添加数据
func setCache(long, short string) error{
	rc := redisClient.Get()
	defer rc.Close()
	_, err := rc.Do("SET",short, long, "EX", "600")
	if err!=nil{
		fmt.Println("redis set error", err)
		return err
	}
	_, err = rc.Do("SET",long, short, "EX", "600")
	if err!=nil{
		fmt.Println("redis set error", err)
		return err
	}
	return nil
}