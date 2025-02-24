package pgx

func SeedPGX(tableStmts []string, indexStmts []string) {
	// db.Get().Exec(tables.GetDropTableStmt())

	for _, stmt := range tableStmts {
		_, err := GetDBX().Exec(stmt) // Changed to GetSQLX() for SQLX usage
		if err != nil {
			panic(err)
		}
	}

	for _, stmt := range indexStmts {
		_, err := GetDBX().Exec(stmt) // Changed to GetSQLX() for SQLX usage
		if err != nil {
			panic(err)
		}
	}
}
