package dbmoc

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

const (
	TxStateCommitted  = "committed"
	TxStateRolledback = "rolledback"
)

// Mock is a mock connection to a database.
type Mock struct {
	stmts []Statement
	tx    *Tx
}

// New returns a new mock connection.
func New() *Mock {
	return new(Mock)
}

// Prepare returns a prepared statement. This is called internally by
// database/sql.
func (m *Mock) Prepare(query string) (driver.Stmt, error) {
	if len(m.stmts) == 0 {
		return nil, fmt.Errorf("dbmoc: unexpected query: got %s, want <none>",
			query)
	}

	var nextStmt Statement
	nextStmt, m.stmts = m.stmts[0], m.stmts[1:]

	if !nextStmt.re.MatchString(query) {
		return nil, fmt.Errorf("dbmoc: unexpected query: got %s, want %s",
			query, nextStmt.MatchQuery)
	}

	if got := calledInTx(); got != nextStmt.InTx {
		return nil, fmt.Errorf("dbmoc: in tx mismatch: got %t, want %t",
			got, nextStmt.InTx)
	}

	nextStmt.query = query
	return &nextStmt, nil
}

// Begin stores a transaction object. Its state can be checked later with
// TxState(). This is called internally by database/sql.
func (m *Mock) Begin() (driver.Tx, error) {
	m.tx = &Tx{}
	return m.tx, nil
}

// Close satisfies driver.Conn interface. No-op.
func (m *Mock) Close() error {
	return nil
}

// Connect satisfies driver.Connector interface. No-op.
func (m *Mock) Connect(context.Context) (driver.Conn, error) {
	return m, nil
}

// Driver satisfies driver.Connector interface. No-op.
func (m *Mock) Driver() driver.Driver {
	return &mockDriver{}
}

// SetStatements takes in statements and checks if they have valid MatchQuery
// patterns. The order of the statements matter. If a SQL statement is seen out
// of order, you'll get an error. After a statement is seen, it's removed from
// the internal queue.
func (m *Mock) SetStatements(stmts []Statement) error {
	for i := 0; i < len(stmts); i++ {
		re, err := regexp.Compile(stmts[i].MatchQuery)
		if err != nil {
			return err
		}
		stmts[i].re = re

		if stmts[i].Rows == nil {
			// Avoids panic, Scan will return sql.ErrNoRows
			stmts[i].Rows = new(Rows)
		}
	}

	m.stmts = stmts

	return nil
}

// LenStatements returns the number of remaining statements.
func (m *Mock) LenStatements() int {
	return len(m.stmts)
}

// TxState returns the transaction's state. It can be "", TxStateCommitted, or
// TxStateRolledback. Note, dbmoc doesn't currently support nested
// transactions.
func (m *Mock) TxState() string {
	return m.tx.state
}

func calledInTx() bool {
	pc := make([]uintptr, 10)
	n := runtime.Callers(1, pc)
	if n == 0 {
		return false
	}

	fs := runtime.CallersFrames(pc[:n])
	for {
		f, more := fs.Next()

		if strings.Contains(f.Function, "*Tx") {
			return true
		}

		if !more {
			break
		}
	}

	return false
}
