package pg

import (
	_ "github.com/bmizerany/pq"
	"github.com/golibs/um"
)

func init() {
	um.Register("postgres", &t_manager{})
}
