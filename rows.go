package dbmoc

import (
	"database/sql/driver"
	"io"
	"reflect"
)

// Rows is the data that gets returned from a database.
type Rows struct {
	// Cols is the name of the columns.
	Cols []string
	// Data is the data for each column. The number of items must match the
	// number of cols.
	Data [][]interface{}

	readIndex int
}

// Columns returns the row columns. This is called internally by database/sql.
func (r *Rows) Columns() []string {
	return r.Cols
}

// Next scans the next driver value into a row. This is called internally by
// database/sql.
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

// Close is a no-op.
func (r *Rows) Close() error {
	return nil
}

// NewRows converts a slice of structs into database rows. The columns names
// are taken from fields that have the `db:"col"` struct tag.
func NewRows(structSlice interface{}) *Rows {
	ss := reflect.ValueOf(structSlice)
	if !isStructSlice(ss) {
		return nil
	}

	return &Rows{
		Cols: dbTags(ss),
		Data: sliceValues(ss),
	}
}

func isStructSlice(v reflect.Value) bool {
	k := v.Kind()
	if k != reflect.Slice {
		return false
	}
	return isStruct(v.Type().Elem())
}

func isStruct(v reflect.Type) bool {
	k := v.Kind()
	return k == reflect.Struct
}

func dbTags(slice reflect.Value) []string {
	sv := slice.Index(0)
	st := makeValue(slice.Type().Elem()).Type()

	var tags []string
	for i := 0; i < sv.NumField(); i++ {
		f := st.Field(i)
		if v, ok := f.Tag.Lookup("db"); ok {
			tags = append(tags, v)
		}
	}

	return tags
}

func makeValue(t reflect.Type) reflect.Value {
	return reflect.New(t).Elem()
}

func sliceValues(slice reflect.Value) [][]interface{} {
	var vals [][]interface{}

	for i := 0; i < slice.Len(); i++ {
		vals = append(vals, structValues(slice.Index(i)))
	}

	return vals
}

func structValues(sv reflect.Value) []interface{} {
	var fields []interface{}

	for i := 0; i < sv.NumField(); i++ {
		fields = append(fields, sv.Field(i).Interface())
	}

	return fields
}
