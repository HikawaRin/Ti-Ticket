package DAO

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"ti-ticket/utils"
)

var (
	_globaldb     *gorm.DB
	_password_len int = 16
    _db_addr      string = "127.0.0.1:4000"
)

func InitConnect() {
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
		panic(err)
	}
	_globaldb = db
	log.Print("db connected success.")
}

func getDB(ctx context.Context) *gorm.DB {
	return _globaldb.WithContext(ctx)
}

func AddUser(usr string) (passwd string, err error) {
	secret := utils.GetSecret(_password_len)
	sql := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s';", usr, secret)
    log.Print("Exec: ", sql)

	db := getDB(context.Background())
	db.Exec(sql)

	// privilege
	// db.Exec(fmt.Sprintf("GRANT ALL ON *.* TO '%s'@'%%'", usr))
	// db.Exec("flush privileges;")
	return secret, nil
}

func DeleteUser(usr string) (err error) { 
    sql := fmt.Sprintf("DROP USER '%s'@'%%';", usr)
    log.Print("Exec: ", sql)

    db := getDB(context.Background())
	db.Exec(sql)

    return nil
}
