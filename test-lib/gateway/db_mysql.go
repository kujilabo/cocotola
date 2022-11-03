package gateway

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDBHost string
var testDBPort int
var testDBURL string

func openMySQLForTest() (*gorm.DB, error) {
	return gorm.Open(gormMySQL.Open(testDBURL), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func InitMySQL(sqlFS embed.FS, dbHost string, dbPort int) {
	testDBHost = dbHost
	testDBPort = dbPort
	testDBURL = fmt.Sprintf("user:password@tcp(%s:%d)/testdb?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", testDBHost, testDBPort)

	setupMySQL(sqlFS)
}

func setupMySQL(sqlFS embed.FS) {
	driverName := "mysql"
	db, err := openMySQLForTest()
	if err != nil {
		log.Fatal(err)
	}
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		log.Fatal(err)
	}
	setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return mysql.WithInstance(sqlDB, &mysql.Config{})
	})
}
