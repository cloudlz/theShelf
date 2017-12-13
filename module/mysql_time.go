//it is list example to introduce timestamp attribute
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"time"
)

//the `test` database is use to test you can rename it
var (
	db, _     = sql.Open("mysql", "root:root@/test?charset=utf8&loc=Local")
	tableName = "`i`" //the table name use to test will be delete after use
)

//according to the parameter to create table
func createDifferentTable(diffCreate string, diffUpdate string) {
	sql := "CREATE TABLE IF NOT EXISTS " + tableName + " (`id` BIGINT UNSIGNED NOT NULL  AUTO_INCREMENT PRIMARY KEY," +
		" `name` varchar(256) NOT NULL," +
		" `datetime1` datetime DEFAULT NULL," +
		" `datetime2` datetime DEFAULT NULL," +
		" `created` TIMESTAMP NOT NULL " + diffCreate + " , " +
		" `timestamp` TIMESTAMP NOT NULL," +
		"`updated`  TIMESTAMP NOT NULL " + diffUpdate + " )" +
		" ENGINE = InnoDB " +
		"DEFAULT CHARSET = utf8 " +
		"COLLATE = utf8_bin;"
	//fmt.Println(sql)
	_, err := db.Exec(sql)
	checkErr(err)

	fmt.Println("success to create table")

}

func deleteTable() {
	_, err := db.Exec("DROP TABLE " + tableName)
	checkErr(err)
	fmt.Println("success to delete table")
}

//get database ddl info
func getDBInfo(parameterName string) {
	stmt, err := db.Query("DESC " + tableName + " " + parameterName)
	checkErr(err)
	for stmt.Next() {
		var Field string
		var Type string
		var Null string
		var Key string
		var Default string
		var Extra string

		err = stmt.Scan(&Field, &Type, &Null, &Key, &Default, &Extra)
		checkErr(err)
		fmt.Println("\tfield:"+Field+"\ttype:"+Type, "\tdefault:"+Default+"\textra:"+Extra)
	}
}

func insertValues() int64 {
	stmt, err := db.Prepare("INSERT " + tableName + " SET name=?,datetime1=?,datetime2=?,updated=?")
	checkErr(err)
	res, err := stmt.Exec("test", time.Now(), time.Now(), time.Now())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println("success insert values(name,datetime1,datetime2,updated) that id =", id)
	return id
}

func updateValues(id int64) {
	stmt, err := db.Prepare("update " + tableName + " set name=?,updated=? where id=?")
	checkErr(err)
	res, err := stmt.Exec("testUpdate", time.Now(), id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("success update value(name,updated) that id =", affect)
}

//get table all value
func getDBValues() {
	stmt, err := db.Query("SELECT * FROM " + tableName)
	checkErr(err)
	fmt.Println("\tid\tname  		datetime1  			datetime2  			created   				timestamp				updated")
	for stmt.Next() {
		var id int
		var name string
		var datetime1 []uint8
		var datetime2 []uint8
		var created []uint8
		var timestamp []uint8
		var updated []uint8

		err = stmt.Scan(&id, &name, &datetime1, &datetime2, &created, &timestamp, &updated)
		checkErr(err)
		fmt.Printf("\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", id, name, string(datetime1), string(datetime2), string(created), string(timestamp), string(updated))
	}
}

//including create table ,show info,insert values ,update values,delete table
func runTable(createParameter string, updateParameter string) {
	createDifferentTable(createParameter, updateParameter)
	fmt.Println("if create table use \n" +
		"\a`created`      TIMESTAMP    NOT NULL   " + createParameter + ",\n" +
		"\a`timestamp`      TIMESTAMP    NOT NULL   ,\n" +
		"\a`updated`      TIMESTAMP    NOT NULL   " + updateParameter + ")")
	fmt.Println("get db info")
	getDBInfo("created")
	getDBInfo("timestamp")
	getDBInfo("updated")
	id := insertValues()
	getDBValues()
	time.Sleep(2000000000)
	updateValues(id)
	getDBValues()
	deleteTable()
}

//the `test` database is use to test you can rename it
func main() {
	fmt.Println("table : \n" +

		"\a`id`      	 	int   ,\n" +
		"\a`name`      	varchar(256),\n" +
		"\a`datetime1`    datetime,\n" +
		"\a`datetime2`    datetime,\n" +
		"\a`created`    	TIMESTAMP,\n" +
		"\a`timestamp`    TIMESTAMP,\n" +
		"\a`updated`      TIMESTAMP)")
	runTable("", "")
	fmt.Println("MySQL datetime is not affected, no matter multiple or one, can not change automatically can only be assigned to change, otherwise it is empty")
	//mysql会默认为表中的第一个timestamp字段（且设置了NOT NULL）隐式设置DEFAULAT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	fmt.Println("if the first timestamp don't set up default,it will default set up DEFAULAT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	fmt.Println("the first timestamp will change to time now when any attribute changed")
	fmt.Println("other timestamp will default 0")
	fmt.Println("--------------------------------------------")

	runTable("DEFAULT CURRENT_TIMESTAMP", "")
	fmt.Println("MySQL datetime is not affected, no matter multiple or one, can not change automatically can only be assigned to change, otherwise it is empty")
	fmt.Println("the first timestamp will not change when other attribute changed")
	fmt.Println("other timestamp will default 0")
	fmt.Println("--------------------------------------------")

	runTable("DEFAULT 0", "DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
	fmt.Println("MySQL datetime is not affected, no matter multiple or one, can not change automatically can only be assigned to change, otherwise it is empty")
	fmt.Println("the first timestamp must be default 0 if table has other timestamp and other timestamp has default value")
	fmt.Println("other timestamp will default 0")
}

//在创建新记录和修改现有记录的时候都对这个数据列刷新：
//TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//在创建新记录的时候把这个字段设置为当前时间，但以后修改时，不再刷新它：
//TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//在创建新记录的时候把这个字段设置为0，以后修改时刷新它：
//TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//在创建新记录的时候把这个字段设置为给定值，以后修改时刷新它：
//TIMESTAMP DEFAULT ‘yyyy-mm-dd hh:mm:ss' ON UPDATE CURRENT_TIMESTAMP

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
