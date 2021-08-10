package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

const cfgFile string = "config.toml"

type cfgInfo struct {
	DBType string
	DBInfo dbInfo
}

type dbInfo struct {
	Usr  string `toml:"DBUser"`
	Pwd  string `toml:"DBPwd"`
	IP   string `toml:"DBIP"`
	Port uint   `toml:"DBPort"`
	Name string `toml:"DBName"`
}

type tableAccount struct {
	accountID   int
	accountName string
}

type accountID int32

//function to get config file info into struct variable
func getTomlInfo(file string, cfgInfoStruct *cfgInfo) {
	if _, err := toml.DecodeFile(cfgFile, cfgInfoStruct); err != nil {
		panic(err)
	}
}

//crud之查
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

//crud之增
func (tbAct tableAccount) insert(db *sql.DB) {
	sql := "insert into account(acct_id, account_name) values (?, ?)"

	// var result string
	result, err := db.Exec(sql, tbAct.accountID, tbAct.accountName)
	if err != nil {
		log.Println("meet error when inserting data", err)
		return
	}
	insertedID, _ := result.LastInsertId()
	affectedRows, _ := result.RowsAffected()
	log.Printf("ID of inserted data : %d, rows of inserted data : %d \n", insertedID, affectedRows)
}

//crud之删
func (actID accountID) delete(db *sql.DB) {
	sql := "delete from account where acct_id = ?"

	// var result string
	result, err := db.Exec(sql, actID)
	if err != nil {
		log.Println("meet error when deleting data", err)
		return
	}
	insertedID, _ := result.LastInsertId()
	affectedRows, _ := result.RowsAffected()
	log.Printf("ID of deleted data : %d, rows of deleted data : %d \n", insertedID, affectedRows)
}

//crud之改
func (tbAct tableAccount) update(db *sql.DB) {
	sql := "update account set account_name= ? where acct_id= ?"

	// var result string
	result, err := db.Exec(sql, tbAct.accountName, tbAct.accountID)
	if err != nil {
		log.Println("meet error when updating data", err)
		return
	}
	insertedID, _ := result.LastInsertId()
	affectedRows, _ := result.RowsAffected()
	log.Printf("ID of updated data : %d, rows of updated data : %d \n", insertedID, affectedRows)
}

// func handleUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	switch r.URL.Path {
// 	case "/account":
// 		switch r.Method {
// 		case "GET":
// 			w.WriteHeader(http.StatusOK)
// 			vars := r.URL.Query()
// 			actID := vars.Get("accountID")
// 			actIDInt, err := strconv.Atoi(actID)
// 			if err != nil {
// 				log.Println("illegal number :", actID)
// 				return
// 			}

// 			var accountID accountID = accountID(actIDInt)
// 			fmt.Fprintf(w, "the result of accountID=", actIDInt, " is : ", accountID.query(db))
// 		default:
// 			fmt.Fprintf(w, "404 no such method : %s\n", r.Method)
// 		}
// 	default:
// 		fmt.Fprintf(w, "404 no such page : %s\n", r.URL)
// 	}
// }

// type database struct {
// 	Database *sql.DB
// }
type dsn string

func (con dsn) handleUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", string(con))
	if err != nil {
		fmt.Println("open mysql database failed while fontend GET, dsn : ", con)
	}
	defer db.Close()

	switch r.URL.Path {
	case "/account":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			getVars := r.URL.Query()
			getActID := getVars.Get("accountID")
			getActIDInt, err := strconv.Atoi(getActID)
			if err != nil {
				log.Println("illegal number : ", getActID)
				return
			}

			var getAccountID accountID = accountID(getActIDInt)
			fmt.Fprintf(w, "query result of accountID = %d is :\n %v\n", getActIDInt, getAccountID.query(db))
		case "DELETE":
			w.WriteHeader(http.StatusOK)
			deleteVars := r.URL.Query()
			deleteActID := deleteVars.Get("accountID")
			deleteActIDInt, err := strconv.Atoi(deleteActID)
			if err != nil {
				log.Println("illegal number : ", deleteActID)
				return
			}

			var deleteAccountID accountID = accountID(deleteActIDInt)
			deleteBeforeResult := deleteAccountID.query(db)
			deleteAccountID.delete(db)
			deleteAfterResult := deleteAccountID.query(db)
			fmt.Fprintf(w, "delete successfully, accountID is %d\n before query result is :\n%v\n after query result is :\n%v", deleteAccountID, deleteBeforeResult, deleteAfterResult)
		case "PUT":
			w.WriteHeader(http.StatusOK)
			putVars := r.URL.Query()
			putActID := putVars.Get("accountID")
			putActIDInt, err := strconv.Atoi(putActID)
			if err != nil {
				log.Println("illegal number : ", putActID)
				return
			}
			putActName := putVars.Get("accountName")

			var putAccountID = accountID(putActIDInt)
			var putTbAct = tableAccount{putActIDInt, putActName}
			putBeforeResult := putAccountID.query(db)
			putTbAct.update(db)
			putAfterResult := putAccountID.query(db)
			fmt.Fprintf(w, "upate successfully, accountID is %d\n before query result is :\n%v\n after query result is :\n%v", putAccountID, putBeforeResult, putAfterResult)
		case "POST":
			w.WriteHeader(http.StatusOK)
			postVars := r.URL.Query()
			postActID := postVars.Get("accountID")
			postActIDInt, err := strconv.Atoi(postActID)
			if err != nil {
				log.Println("illegal number : ", postActID)
				return
			}
			postActName := postVars.Get("accountName")

			var postAccountID = accountID(postActIDInt)
			var postTbAct = tableAccount{postActIDInt, postActName}
			postBeforeResult := postAccountID.query(db)
			postTbAct.insert(db)
			postAfterResult := postAccountID.query(db)
			fmt.Fprintf(w, "insert successfully, accountID is %d\n before query result is :\n%v\n after query result is :\n%v", postAccountID, postBeforeResult, postAfterResult)
		default:
			fmt.Fprintf(w, "404 no such method : %s\n", r.Method)
		}
	default:
		fmt.Fprintf(w, "404 no such page : %s\n", r.URL)
	}
}

func main() {
	log.Println("starting")

	//get database connection info
	var cfg cfgInfo
	getTomlInfo(cfgFile, &cfg)
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.DBInfo.Usr, cfg.DBInfo.Pwd, cfg.DBInfo.IP, cfg.DBInfo.Port, cfg.DBInfo.Name)
	fmt.Println("connecting mysql database...")
	fmt.Println("connection is : " + dbConn)

	//create connection with mysql database
	//*notice, need to consider do a loop to try reconnect database if database is down
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("connect mysql database failed.")
		panic(err)
	}
	defer db.Close()

	log.Println("connect mysql database succeed.")

	//start server to get request from frontend
	dbCoon2 := dsn(dbConn)
	http.HandleFunc("/account", dbCoon2.handleUsers)
	http.ListenAndServe(":12610", nil)

	// ///test method of data
	// //测试查
	// var a accountID = 0
	// qResult := a.query(db)
	// log.Println("the result of accountID=", a, " is : ", qResult)

	// //测试增(查)删(查)
	// var insertData = tableAccount{1, "zhuyu2"}
	// var b accountID = 1
	// insertData.insert(db)
	// qResult = b.query(db)
	// log.Println("the result of accountID=", b, " is : ", qResult)
	// b.delete(db)
	// qResult = b.query(db)
	// log.Println("the result of accountID=", b, " is : ", qResult)

	// //测试更(查)更（查）
	// var c accountID = 2
	// var updateData1 = tableAccount{2, "zhuyu3"}
	// updateData1.update(db)
	// log.Println("the result of accountID=", c, " is : ", c.query(db))
	// var updateData2 = tableAccount{2, "zhuyu2"}
	// updateData2.update(db)
	// log.Println("the result of accountID=", c, " is : ", c.query(db))
}
