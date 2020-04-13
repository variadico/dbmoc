package dbmoc

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
)

// Statement holds information for expected database statements. They're used
// to validate expected and actual queries.
type Statement struct {
	// Pattern is a Go regexp pattern. This will be used to match against a
	// query, so you know it showed up.
	MatchPattern string
	// Args are used to verify that the correct arguments got passed to your
	// query. These are checked for equality with the actual arguments that
	// your function passes in.
	MatchArgs []interface{}

	// InTx will check to see if a statement is being called within a
	// transaction.
	InTx bool

	// SkipMatchArgs allows the arg matching to be skipped.
	SkipMatchArgs bool

	// Rows returns the rows from a mocked query.
	Rows *Rows
	// Results returns the results from a mocked executed statement.
	Result *Result

	re    *regexp.Regexp
	query string
}

// Exec handles calls to db.Exec.
func (s *Statement) Exec(args []driver.Value) (driver.Result, error) {
	if !s.SkipMatchArgs {
		iargs := toInterfaceSlice(args)
		if err := sliceEqual(s.MatchArgs, iargs); err != nil {
			return nil, fmt.Errorf("dbmoc: unexpected exec args: %w", err)
		}
	}

	return s.Result, nil
}

// Query handles calls to db.Query and db.QueryRow.
func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
	if !s.SkipMatchArgs {
		iargs := toInterfaceSlice(args)
		if err := sliceEqual(s.MatchArgs, iargs); err != nil {
			return nil, fmt.Errorf("dbmoc: unexpected query args: %w", err)
		}
	}

	return s.Rows, nil
}

// Close satisfies the driver.Stmt interface. No-op.
func (s *Statement) Close() error {
	return nil
}

// NumInput satisfies the driver.Stmt interface. It always returns -1. No-op.
func (s *Statement) NumInput() int {
	return -1
}

func sliceEqual(a, b []interface{}) error {
	if len(a) != len(b) {
		return fmt.Errorf(
			"len args not equal: %v(%d) vs %v(%d)",
			a, len(a),
			b, len(b),
		)
	}

	for i := 0; i < len(a); i++ {
		if !reflect.DeepEqual(a[i], b[i]) {
			return fmt.Errorf(
				"element not equal: %[1]T(%[1]v) vs %[2]T(%[2]v)",
				a[i], b[i],
			)
		}
	}

	return nil
}

// toInterfaceSlice converts driver.Value to interface{}. For primitive
// types, the underlying value is one of these types.
//
//   int64
//   float64
//   bool
//   []byte
//   string
//   time.Time
//
// This is important if you try to use reflect.DeepEqual.
func toInterfaceSlice(vs []driver.Value) []interface{} {
	if len(vs) == 0 {
		return nil
	}

	is := make([]interface{}, 0, len(vs))
	for _, v := range vs {
		is = append(is, interface{}(v))
	}
	return is
}
