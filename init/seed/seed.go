package seed

import (
	"log"

	db "github.com/braumsmilk/go-db"
)

func SeedPostgresDatabase(tableStmts []string, indexStmts []string) {
	// db.Get().Exec(tables.GetDropTableStmt())

	for _, stmt := range tableStmts {
		log.Printf("%s", stmt)
		_, err := db.Get().Exec(stmt)
		if err != nil {
			panic(err)
		}
	}

	for _, stmt := range indexStmts {
		log.Printf("%s", stmt)
		_, err := db.Get().Exec(stmt)
		if err != nil {
			panic(err)
		}
	}
}
