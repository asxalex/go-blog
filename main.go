package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"myblog/models"
	_ "myblog/routers"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func createTable() {
	name := "default"                          //数据库别名
	force := false                             //不强制建数据库
	verbose := true                            //打印建表过程
	err := orm.RunSyncdb(name, force, verbose) //建表
	if err != nil {
		beego.Error(err)
	}
}

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	//orm.RegisterDataBase("default", "mysql", "root:root@/myblogs?charset=utf8")
	orm.RegisterModel(new(models.Post), new(models.Tag))
	createTable()
	models.File2DB()
}

func main() {
	beego.Run()
}
