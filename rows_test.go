package dbmoc

import (
	"reflect"
	"testing"
)

func TestNewRows(t *testing.T) {
	rowData := []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}{
		{ID: 1, Name: "foo"},
		{ID: 2, Name: "bar"},
		{ID: 3, Name: "fizz"},
	}

	rows := NewRows(rowData)

	wantCols := []string{
		"id",
		"name",
	}
	wantData := [][]interface{}{
		{int64(1), "foo"},
		{int64(2), "bar"},
		{int64(3), "fizz"},
	}

	if len(rows.Cols) != len(wantCols) {
		t.Error("unexpected columns")
		t.Fatalf("got=%v; want=%v", rows.Cols, wantCols)
	}
	for i := range rows.Cols {
		if !reflect.DeepEqual(rows.Cols[i], wantCols[i]) {
			t.Error("unexpected column")
			t.Fatalf("got=%v; want=%v", rows.Cols[i], wantCols[i])
		}
	}

	if len(rows.Data) != len(wantData) {
		t.Error("unexpected data")
		t.Fatalf("got=%v; want=%v", rows.Data, wantData)
	}
	for i := range rows.Data {
		if !reflect.DeepEqual(rows.Data[i], wantData[i]) {
			t.Error("unexpected data")
			t.Fatalf("got=%v; want=%v", rows.Data[i], wantData[i])
		}
	}
}
