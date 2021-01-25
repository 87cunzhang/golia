package golib

import (
	"log"
	"strconv"
	"time"

	"zkds/src/confParser"

	"github.com/xormplus/xorm"
)

var engines = make(map[string]*xorm.Engine)

func init() {
	conf := confParser.DefaultConf()
	data, err := conf.DIY("database")

	if err != nil {
		log.Fatal(err)
	}

	databaseConf := data.(map[string]interface{})

	for k := range databaseConf {
		driverName := conf.DefaultString("database::"+k+"::driverName", "mysql")
		host := conf.DefaultString("database::"+k+"::host", "127.0.0.1")
		port := conf.DefaultInt("database::"+k+"::port", 3306)
		user := conf.DefaultString("database::"+k+"::user", "root")
		password := conf.DefaultString("database::"+k+"::password", "")
		charset := conf.DefaultString("database::"+k+"::charset", "utf8")
		dbName := conf.String("database::" + k + "::dbName")
		maxIdleConns := conf.DefaultInt("database::"+k+"::maxIdleConns", 0)
		maxOpenConns := conf.DefaultInt("database::"+k+"::maxOpenConns", 0)
		connMaxLifetime := conf.DefaultInt("database::"+k+"::connMaxLifetime", 0)
		dsnString := user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbName + "?charset=" + charset
		engines[k], err = xorm.NewEngine(driverName, dsnString)

		if err != nil {
			log.Fatal(err)
		}

		if err := engines[k].Ping(); err != nil {
			log.Fatal(err)
		}

		engines[k].SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
		engines[k].SetMaxIdleConns(maxIdleConns)
		engines[k].SetMaxOpenConns(maxOpenConns)
	}
}

func defaultDB() *xorm.Engine {
	return engines["default"]
}

func DB(DBName string) *xorm.Engine {
	return engines[DBName]
}
