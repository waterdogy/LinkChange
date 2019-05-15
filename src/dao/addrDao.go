package dao

import (
	"domain"
	"sync"
)

//长短链接数据库
var ltos = make(map[string]string)
var stol = make(map[string]string)

var mu sync.RWMutex//操作数据库需要加锁

//查询短地址是否存在
func CheckShortAddr(short string) bool{
	mu.RLock()
	defer mu.RUnlock()
	_,ok := stol[short]
	return ok
}

//根据短地址返回长地址
func FindLongAddr(short string) string{
	mu.RLock()
	defer mu.RUnlock()
	return stol[short]
}

//查询长地址是否存在
func CheckLongAddr(long string) bool{
	mu.RLock()
	defer mu.RUnlock()
	_,ok := ltos[long]
	return ok
}

//根据长地址返回短地址
func FindShortAddr(long string) string{
	mu.RLock()
	defer mu.RUnlock()
	return ltos[long]
}

//将新的地址映射关系存入数据库
func InsertAddr(a domain.Addr){
	mu.Lock()
	defer mu.Unlock()
	ltos[a.LongAddr] = a.ShortAddr
	stol[a.ShortAddr] = a.LongAddr
}