package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"models"
	"utils"
)

var confFilename = "D:/LinkChange/src/conf/app.ini"
var tp, user, pw, tableName string

//初始化建表
func init(){
	fmt.Println("开始初始化建表")
	//建立数据库连接
	//读取ini配置文件
	iniParser := utils.IniParser{}
	if err:=iniParser.Load(confFilename); err!=nil{
		fmt.Println("loading ini file fail",err)
		return
	}
	//数据库类型
	tp = iniParser.GetString("database","Type")
	//数据库用户名
	user = iniParser.GetString("database","User")
	//数据库密码
	pw = iniParser.GetString("database","Password")
	//数据库表名
	tableName = iniParser.GetString("database","Table")
	//数据库驱动
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	if err!=nil{
		panic("连接数据库失败")
	}
	defer db.Close()
	//如果没有表就建表
	if !db.HasTable(&models.Addr{}){
		if err:= db.Set("gorm:table_options","ENGINE=InnoDB DEFAULT CHARSET=utf8").
			CreateTable(&models.Addr{}).Error; err!=nil{
			panic(err)
		}
	}
}

//查询短地址是否存在
func CheckShortAddr(short string) bool{
	var count int
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	defer db.Close()
	if err!=nil{
		panic("连接数据库失败")
	}
	db.Model(&models.Addr{}).Where("short_addr=?", short).Count(&count)
	return count==1
}

//根据短地址返回长地址
func FindLongAddr(short string) string{
	var addr models.Addr
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	defer db.Close()
	if err!=nil{
		panic("连接数据库失败")
	}
	db.Model(&models.Addr{}).Where("short_addr=?", short).First(&addr)
	return addr.LongAddr
}

//查询长地址是否存在
func CheckLongAddr(long string) bool{
	//缓存中没有去数据库查找
	var count int
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	defer db.Close()
	if err!=nil{
		panic("连接数据库失败")
	}
	db.Model(&models.Addr{}).Where("long_addr=?", long).Count(&count)
	return count==1
}

//根据长地址返回短地址
func FindShortAddr(long string) string{
	var addr models.Addr
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	defer db.Close()
	if err!=nil{
		panic("连接数据库失败")
	}
	db.Model(&models.Addr{}).Where("long_addr=?", long).First(&addr)
	return addr.ShortAddr
}

//将新的地址映射关系存入数据库
func InsertAddr(a models.Addr) bool{
	dr := user+":"+pw+"@/"+tableName+"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(tp, dr)
	defer db.Close()
	if err!=nil{
		panic("连接数据库失败")
	}
	err = db.Create(&a).Error
	if err!=nil{
		return false
	}
	return true
}