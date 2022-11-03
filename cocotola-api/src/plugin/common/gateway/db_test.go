package gateway_test

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	"github.com/kujilabo/cocotola/cocotola-api/src/sqls"
	liberrors "github.com/kujilabo/cocotola/lib/errors"
)

func dbList() map[string]*gorm.DB {
	dbList := make(map[string]*gorm.DB)
	m, err := openMySQLForTest()
	if err != nil {
		panic(err)
	}

	dbList["mysql"] = m

	// s, err := openSQLiteForTest()
	// if err != nil {
	// 	panic(err)
	// }
	// dbList["sqlite3"] = s

	return dbList
}

func setupDB(db *gorm.DB, driverName string, withInstance func(sqlDB *sql.DB) (database.Driver, error)) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pos := strings.Index(wd, "src")
	// dir := wd[0:pos] + "sqls/" + driverName

	driver, err := withInstance(sqlDB)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to WithInstance. err: %w", err))
	}

	d, err := iofs.New(sqls.SQL, driverName)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("iofs", d, driverName, driver)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to NewWithDatabaseInstance. err: %w", err))
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(liberrors.Errorf("failed to Up. driver:%s, err: %w", driverName, err))
		}
	}
}
