// Package dbmoc creates a dummy database driver. The driver checks for SQL
// statements and arguments. It also returns arbitrary rows of data.
package dbmoc

import "database/sql/driver"

// mockDriver is a mock DB mockDriver.
type mockDriver struct{}

// Open return
func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &Mock{}, nil
}
