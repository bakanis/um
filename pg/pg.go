package pg

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"github.com/golibs/um"
)

type Manager struct {
	session *sql.DB

	createUserStmt *sql.Stmt // prepared statement for creating user
}

func init() {
	um.Register("postgres", &Manager{})
}
