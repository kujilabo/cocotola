package gateway

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gormSQLite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDBFile string

func openSQLiteForTest() (*gorm.DB, error) {
	return gorm.Open(gormSQLite.Open(testDBFile), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func InitSQLite(sqlFS embed.FS) {
	testDBFile = "./test.db"
	os.Remove(testDBFile)
	setupSQLite(sqlFS)
}

func setupSQLite(sqlFS embed.FS) {
	driverName := "sqlite3"
	db, err := openSQLiteForTest()
	if err != nil {
		log.Fatal(err)
	}
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		log.Fatal(err)
	}
	setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	})
}
