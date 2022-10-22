package DAO

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_globaldb *gorm.DB = nil
	_db_addr  string   = "127.0.0.1:4000"
)

func InitConnect() error {
	user := os.Getenv("DATABASE_USER")
	passwd := os.Getenv("DATABASE_PASSWD")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/mysql?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		passwd,
		_db_addr)

	log.Printf("db dsn: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	_globaldb = db
	log.Print("db connected success.")
	return nil
}

func AddUser(usr string, passwd string) (err error) {
	sql := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s';", usr, passwd)
	log.Print("Exec: ", sql, "\n")

	db := getDB(context.Background())
	if db != nil {
		db.Exec(sql)
	}

	// privilege
	// db.Exec(fmt.Sprintf("GRANT ALL ON *.* TO '%s'@'%%'", usr))
	// db.Exec("flush privileges;")
	return nil
}

func UpdateUser(usr string, passwd string) (err error) {
	if !existUser(usr) {
		return AddUser(usr, passwd)
	}
	sql := fmt.Sprintf("ALTER USER '%s'@'%%' IDENTIFIED BY '%s';", usr, passwd)
	log.Print("Exec: ", sql)

	db := getDB(context.Background())
	if db != nil {
		db.Exec(sql)
	}
	return nil
}

func DeleteUser(usr string) (err error) {
	sql := fmt.Sprintf("DROP USER '%s'@'%%';", usr)
	log.Print("Exec: ", sql, "\n")

	db := getDB(context.Background())
	if db != nil {
		db.Exec(sql)
	}
	return nil
}

func Privilege_table() map[int]string {
	return map[int]string{
		1: "SELECT",
		2: "INSERT",
		4: "UPDATE",
		8: "DELETE",
	}
}

const _grant_template = "GRANT %s ON *.* TO '%s'@'%%';"

func GrantUser(usr string, code int) (string, error) {
	pri_str := ""
	for mask, priv := range Privilege_table() {
		if (mask & code) != 0 {
			pri_str += "," + priv
		}
	}
	sql := fmt.Sprintf(_grant_template, pri_str[1:], usr)
	log.Print(sql, "\n")

	db := getDB(context.Background())
	if db != nil {
		db.Exec(sql)
		db.Exec("flush privileges;")
	}
	return pri_str, nil
}

type temporary_user struct {
	gorm.Model
	User        string
	Select_priv string
	Insert_priv string
	Update_priv string
	Delete_priv string
}

func getDB(ctx context.Context) *gorm.DB {
	if _globaldb == nil {
		log.Print("Database not ready")
		return nil
	}
	return _globaldb.WithContext(ctx)
}

func existUser(usr string) bool {
	var u temporary_user
	db := getDB(context.Background())
	if db != nil {
		db.Table("user").Where("User = ?", usr).Scan(&u)
	}
	return u.User == usr
}
