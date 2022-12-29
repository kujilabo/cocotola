package gateway

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDBHost string
var testDBPort int

func openMySQLForTest() (*gorm.DB, error) {
	c := mysql.Config{
		DBName:               "testdb",
		User:                 "user",
		Passwd:               "password",
		Addr:                 fmt.Sprintf("%s:%d", testDBHost, testDBPort),
		Net:                  "tcp",
		ParseTime:            true,
		MultiStatements:      true,
		Params:               map[string]string{"charset": "utf8"},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
		Loc:                  time.UTC,
		// Loc:                  jst,
	}
	dsn := c.FormatDSN()
	return gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func setupMySQL(sqlFS embed.FS, db *gorm.DB) error {
	driverName := "mysql"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}

	return setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_mysql.WithInstance(sqlDB, &migrate_mysql.Config{})
	})
}

func InitMySQL(sqlFS embed.FS, dbHost string, dbPort int) (*gorm.DB, error) {
	testDBHost = dbHost
	testDBPort = dbPort
	db, err := openMySQLForTest()
	if err != nil {
		return nil, err
	}

	if err := setupMySQL(sqlFS, db); err != nil {
		return nil, err
	}

	return db, nil
}
