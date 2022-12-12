package gateway

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(username, password, host string, port int, database string) (*gorm.DB, error) {
	c := mysql.Config{
		DBName:          database,
		User:            username,
		Passwd:          password,
		Addr:            fmt.Sprintf("%s:%d", host, port),
		Net:             "tcp",
		ParseTime:       true,
		MultiStatements: true,
		Params:          map[string]string{"charset": "utf8"},
		Collation:       "utf8mb4_unicode_ci",
		Loc:             jst,
	}
	dsn := c.FormatDSN()
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&multiStatements=true", username, password, host, port, database)
	return gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func MigrateMySQLDB(db *gorm.DB, sqlFS embed.FS) error {
	driverName := "mysql"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}

	return migrateDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_mysql.WithInstance(sqlDB, &migrate_mysql.Config{})
	})
}
