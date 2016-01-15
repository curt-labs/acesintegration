package database

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
)

var (
	EmptyDb      = flag.String("clean", "", "bind empty database with structure defined")
	DB           *sql.DB
	VCDB         *sql.DB
	MongoSession *mgo.Session
	MongoDB      string
)

type Scanner interface {
	Scan(...interface{}) error
}

func ConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")

		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}
	return "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true&loc=America%2FChicago"
}

func VcdbConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("VCDB_NAME")

		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}
	return "root:@tcp(127.0.0.1:3306)/vcdb?parseTime=true&loc=America%2FChicago"
}

func Init() error {
	var err error
	if DB == nil {
		DB, err = sql.Open("mysql", ConnectionString())
		if err != nil {
			return err
		}
	}
	if VCDB == nil {
		VCDB, err = sql.Open("mysql", VcdbConnectionString())
		if err != nil {
			return err
		}
	}
	return nil
}

func InitMongo() error {
	var err error
	if MongoSession == nil {
		MongoSession, err = mgo.DialWithInfo(mongoConnectionString())
		if err != nil {
			return err
		}
	}
	MongoDB = mongoConnectionString().Database
	return nil
}

func mongoConnectionString() *mgo.DialInfo {
	var info mgo.DialInfo
	addr := os.Getenv("MONGO_URL")
	if addr == "" {
		addr = "127.0.0.1"
	}
	addrs := strings.Split(addr, ",")
	info.Addrs = append(info.Addrs, addrs...)

	info.Username = os.Getenv("MONGO_USERNAME")
	info.Password = os.Getenv("MONGO_PASSWORD")
	info.Database = os.Getenv("MONGO_DATABASE")
	info.Timeout = time.Second * 2
	info.FailFast = true
	if info.Database == "" {
		info.Database = "aries"
	}
	info.Source = "admin"

	return &info
}

func TestMongoConnection() error {
	err := InitMongo()
	if err != nil {
		return err
	}

	return MongoSession.Ping()

}
