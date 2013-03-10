package pg

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"github.com/golibs/um"
)

type Manager struct {
	dns     string
	session *sql.DB
}

func init() {
	um.Register("postgres", &Manager{})
}
