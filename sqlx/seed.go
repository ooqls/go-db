package sqlx

import (
	"go.uber.org/zap"
)

func SeedSQLX(tableStmts []string, indexStmts []string) {
	// db.Get().Exec(tables.GetDropTableStmt())

	for _, stmt := range tableStmts {
		l.Info("executing table statement", zap.String("stmt", stmt))
		_, err := GetSQLX().Exec(stmt) // Changed to GetSQLX() for SQLX usage
		if err != nil {
			l.Error("failed to execute table statement", zap.String("stmt", stmt), zap.Error(err))
			panic(err)
		}
	}

	for _, stmt := range indexStmts {
		l.Info("executing index statement", zap.String("stmt", stmt))
		_, err := GetSQLX().Exec(stmt) // Changed to GetSQLX() for SQLX usage
		if err != nil {
			l.Error("failed to execute index statement", zap.String("stmt", stmt), zap.Error(err))
			panic(err)
		}
	}
}
