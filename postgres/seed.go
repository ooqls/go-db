package postgres

import (
	"log"
)

func Seed(tableStmts []string, indexStmts []string) {
	// db.Get().Exec(tables.GetDropTableStmt())

	for _, stmt := range tableStmts {
		log.Printf("%s", stmt)
		_, err := Get().Exec(stmt)
		if err != nil {
			panic(err)
		}
	}

	for _, stmt := range indexStmts {
		log.Printf("%s", stmt)
		_, err := Get().Exec(stmt)
		if err != nil {
			panic(err)
		}
	}
}
