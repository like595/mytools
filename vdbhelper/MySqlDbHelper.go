package vdbhelper

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/like595/mytools/vtools"
	"strings"
)

// MySql工具
type MySqlDBHelper struct {
	ip       string
	port     int
	username string
	password string
	dbname   string
	DB       *sql.DB
	//mylog MyLog.MyLog
}

// 初始化配置文件
func (this *MySqlDBHelper) init() {

	iniUtil := vtools.IniUtil{}
	iniUtil.Init("./conf/Config.ini")
	this.ip = iniUtil.GetVal("mysql", "ip")
	this.port = iniUtil.GetInt("mysql", "port")
	this.username = iniUtil.GetVal("mysql", "username")
	this.password = iniUtil.GetVal("mysql", "password")
	this.dbname = iniUtil.GetVal("mysql", "dbname")
}

// 初始化配置文件
func (this *MySqlDBHelper) SetParam(ip string, port int, username, password, dbname string) {
	this.ip = ip
	this.port = port
	this.username = username
	this.password = password
	this.dbname = dbname

}

func (this *MySqlDBHelper) Open() bool {
	if this.ip == "" {
		this.init()
	}
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", this.username, this.password, this.ip, this.port, this.dbname)
	vtools.SugarLogger.Info("DZXT.dzxt_1.", "打开数据库", dataSource)

	db, err := sql.Open("mysql", dataSource)
	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)
	if err == nil {
		this.DB = db
	} else {
		vtools.SugarLogger.Error("DZXT.dzxt_1.", "MySql数据库打开失败，dataSource=", dataSource, "err=", err)
		return false
	}
	return true
}

/*
根据传入的数据库地址打开数据库
*/
func (this *MySqlDBHelper) OpenByURL(dataSource string) bool {
	db, err := sql.Open("mysql", dataSource)
	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)

	vtools.SugarLogger.Info("DZXT.dzxt_1.", "打开数据库*OpenByURL", dataSource)

	if err == nil {
		this.DB = db
		vtools.SugarLogger.Info("DZXT.dzxt_1.", "门架数据库打开成", dataSource)
		return true
	} else {
		vtools.SugarLogger.Info("DZXT.dzxt_1.", "MySql数据库打开失败", err)
		return false
	}
}

func (this *MySqlDBHelper) Close() {
	if this.DB != nil {
		this.DB.Close()
	}
}

// 查询数据
func (this *MySqlDBHelper) Select(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		vtools.SugarLogger.Error("DZXT.dzxt_1.", "MySql语句执行失败,sql=", query, "参数=", args, "异常消息=", err)
	}
	//defer rows.Close()
	//defer this.db.Close()
	return rows, err
}

// 执行sql语句
func (this *MySqlDBHelper) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := this.DB.Exec(query, args...)
	if err != nil {
		if strings.Index(err.Error(), "PRIMARY") == -1 {
			vtools.SugarLogger.Error("DZXT.dzxt_1.", "MySql语句执行失败,sql=", query, "参数=", args, "异常消息=", err.Error())
		} else {
			//err = nil
		}
	}
	//defer this.db.Close()
	return res, err
}
