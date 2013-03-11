package pg

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

const testDns = "user=postgres dbname=usermanagement_testing sslmode=disable"

// setup test data
func testSetup() *sql.DB {
	var setupSql = `
		create sequence um_users_id_seq start 1;

		create table um_users (
			id bigint not null default nextval('um_users_id_seq'),
			user_name varchar(64) not null,
			email_addr varchar(128) not null default '',
			status integer not null default 0,
			hash varchar(128) not null default '',
			salt varchar(32) not null default '',
			created_on timestamp with time zone not null default 'now',
			last_login timestamp with time zone not null default 'epoch',

			constraint pk_um_users primary key(id),
			constraint ux_um_users_username unique(user_name)
		);

		create index ix_um_users_email_addr on um_users(email_addr);
		create index ix_um_users_created_on on um_users(created_on);
		create index ix_um_users_last_login on um_users(last_login);

		insert into um_users(user_name, email_addr, status) values
			('fixtureuser1', 'fixtureuser1@example.com', 10),
			('fixtureuser2', 'fixtureuser2@example.com', 20);`

	var err error
	var session *sql.DB

	session, err = sql.Open("postgres", testDns)
	if err != nil {
		panic(err)
	}
	_, err = session.Exec(setupSql)
	if err != nil {
		panic(err)
	}

	return session
}

func testTearDown(session *sql.DB) {
	var teardownSql = `
		drop table um_users;
		drop sequence um_users_id_seq;`
	if session != nil {
		var err error
		_, err = session.Exec(teardownSql)
		if err != nil {
			panic(err)
		}
		session.Close()
	}
}

// assertRecord returns true if the table contains a record with the specified properties
func assertRecord(session *sql.DB, table string, props map[string]interface{}) bool {
	if session == nil {
		panic("database session is not established")
	}

	var buff bytes.Buffer
	var vals = make([]interface{}, len(props))
	var i = 0
	var count int

	buff.WriteString("select count(*) from ")
	buff.WriteString(table)
	buff.WriteString(" where ")

	for k, v := range props {
		if i != 0 {
			buff.WriteString(" and ")
		}
		vals[i] = v
		i++
		buff.WriteString(fmt.Sprintf("%s = $%d", k, i))
	}

	stmt, err := session.Prepare(buff.String())
	if err != nil {
		panic(err)
	}
	row := stmt.QueryRow(vals...)
	err = row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return (count != 0)
}
