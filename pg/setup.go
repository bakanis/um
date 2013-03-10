package pg

import "database/sql"

// Setup prepares the manager's database connections
func (this *Manager) Setup(dns string) error {
	var err error
	this.session, err = sql.Open("postgres", dns)
	return err
}

// Close cleans up the database connection of the current manager
func (this *Manager) Close() error {
	if this.session != nil {
		return this.session.Close()
	}
	return nil
}
