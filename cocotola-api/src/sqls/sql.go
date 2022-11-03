package sqls

import (
	"embed"
	_ "embed"
)

//go:embed mysql/*.sql
//go:embed sqlite3/*.sql
var SQL embed.FS

func init() {
	if _, err := SQL.ReadFile("mysql/2020050101_create_organization.up.sql"); err != nil {
		panic(err)
	}
}
