package pg

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

const c_testDns = "host=127.0.0.1 user=postgres dbname=usermanagement_testing sslmode=disable"

// setup test data
func testSetup() *sql.DB {
	// test password is 'Password123'
	var setupSql = `
		create sequence um_users_id_seq start 1;

		create table um_users (
			id bigint not null default nextval('um_users_id_seq'),
			email_addr varchar(128) not null,
			display_name varchar(64) not null default '',
			status integer not null default 0,
			hash varchar(128) not null default '',
			salt varchar(32) not null default '',
			created_on timestamp with time zone not null default 'now',
			last_login timestamp with time zone not null default 'epoch',

			constraint pk_um_users primary key(id),
			constraint ux_um_users_email_addr unique(email_addr)
		);

		create index ix_um_users_created_on on um_users(created_on);
		create index ix_um_users_last_login on um_users(last_login);

		insert into um_users(display_name, email_addr, status, salt, hash) values
			('Fixture User 1', 'fixtureuser1@example.com', 10, '950333c33bf84dbbceea95126a23974e', '2432612431302471466d6e7241755864322f7647727a4469574754362e702e3465697432472e6b54502e652e6174396167327732587a492e65732f32'),
			('Fixture User 2', 'fixtureuser2@example.com', 20, '950333c33bf84dbbceea95126a23974e', '2432612431302471466d6e7241755864322f7647727a4469574754362e702e3465697432472e6b54502e652e6174396167327732587a492e65732f32');`

	var err error
	var session *sql.DB

	session, err = sql.Open("postgres", c_testDns)
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
