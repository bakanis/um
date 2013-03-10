package um

import "time"

type User struct {
	Id        uint64
	UserName  string
	EmailAddr string
	Status    int32
	hash      string
	salt      string
	CreatedOn time.Time
	LastLogin time.Time
}
