package dbmoc

import (
	"database/sql/driver"
	"io"
)

type RowData [][]interface{}

type Rows struct {
	Cols []string
	Data RowData

	readIndex int
}

func (r *Rows) Columns() []string {
	return r.Cols
}

func (r *Rows) Next(dest []driver.Value) error {
	rs := r.Data
	if len(rs) == r.readIndex {
		return io.EOF
	}

	for i := range rs[r.readIndex] {
		dest[i] = rs[r.readIndex][i]
	}

	r.readIndex++
	return nil
}

func (r *Rows) Close() error {
	return nil
}
