package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

const cfgFile string = "config.toml"

type cfgInfo struct {
	DBType string
	DB     DBInfo `toml:"DBInfo"`
}

type DBInfo struct {
	DBUser string
	DBPwd  string
	DBIP   string
	DBPort uint
	DBName string
}

type tableAccount struct {
	accountID   int32
	accountName string
}

type accountID int32

//function to get config file info into struct variable
func getTomlInfo(file string, cfgInfoStruct *cfgInfo) {
	if _, err := toml.DecodeFile(cfgFile, cfgInfoStruct); err != nil {
		panic(err)
	}
}

func (actID accountID) query(db *sql.DB) []tableAccount {
	sql := "select acct_id, account_name from account where acct_id = ?"

	//定义结果集结构体数组用于接受查询结果
	var accountResult []tableAccount
	//get query result from database
	//从mysql数据库中使用接收器的值查询出结果存放于结果字段对应结构体切片中
	rows, err := db.Query(sql, actID)
	if err != nil {
		log.Println("query data base failed with accountID : ", actID, err)
		return accountResult
	}
	defer rows.Close()
	for rows.Next() {
		var act tableAccount
		rows.Scan(&act.accountID, &act.accountName)
		accountResult = append(accountResult, act)
	}

	return accountResult
}

func main() {
	log.Println("starting")

	//get database connection info
	var cfg cfgInfo
	getTomlInfo(cfgFile, &cfg)
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DB.DBUser, cfg.DB.DBPwd, cfg.DB.DBIP, cfg.DB.DBPort, cfg.DB.DBName)
	fmt.Println("connecting mysql database...")
	fmt.Println("connection is : " + dbConn)

	//create connection with mysql database
	//*notice, need to consider do a loop to try connect database if database is down
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("connect mysql database failed.")
		panic(err)
	}
	defer db.Close()

	log.Println("connect mysql database succeed.")

	var a accountID = 0
	qResult := a.query(db)
	log.Println("the result of accountID=", a, " is : ", qResult)

}
