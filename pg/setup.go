package pg

import "database/sql"

// Setup prepares the manager's database connections
func (this *Manager) Setup(dns string) error {
	var err error
	this.session, err = sql.Open("postgres", dns)

	if err == nil {
		query := "insert into um_users(user_name, email_addr, status, created_on, last_login) values($1, $2, $3, $4, $5) returning id;"
		this.createUserStmt, err = this.session.Prepare(query)
	}

	return err
}

// Close cleans up the database connection of the current manager
func (this *Manager) Close() error {
	if this.session != nil {
		return this.session.Close()
	}
	return nil
}
