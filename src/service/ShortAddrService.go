package service

import "dao"


//查询是否数据库中有短地址
func TranShortToLong(short string)(string, bool){
	//先去缓存中查找
	if ok := checkRedisShort(short);ok{
		return getRedisLong(short), true
	}
	//再去数据库中查找
	if ok := dao.CheckShortAddr(short); ok{//数据库暂时没有存短地址
		long := dao.FindLongAddr(short)
		//添加缓存
		_ = setCache(long, short)
		return long, true
	}
	return "", false
}
