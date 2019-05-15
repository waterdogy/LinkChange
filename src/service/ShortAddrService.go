package service

import "dao"


//查询是否数据库中有短地址
func TranShortToLong(short string)(string, bool){
	if ok := dao.CheckShortAddr(short); !ok{//数据库暂时没有存短地址
		return "", false
	}
	return dao.FindLongAddr(short), true
}
