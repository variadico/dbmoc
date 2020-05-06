package dbmoc

import (
	"database/sql"
	"log"
)

func ExampleMock_query() {
	mock := New()
	db := sql.OpenDB(mock)

	wantStmts := []Statement{
		{
			MatchPattern: `^select .*\s+from foo.bar`,
			MatchArgs:    []interface{}{3},

			Rows: &Rows{
				Cols: []string{"id", "name"},
				Data: RowData{
					{1, "fizz"},
					{2, "buzz"},
				},
			},
		},
	}
	if err := mock.SetStatements(wantStmts); err != nil {
		log.Fatalln(err)
	}

	db.Query("select id, name from foo.bar where id < $1", 3)
	// More database calls....
}

func ExampleMock_exec() {
	mock := New()
	db := sql.OpenDB(mock)

	wantStmts := []Statement{
		{MatchPattern: "^select pg_advisory_lock"},
		{MatchPattern: "^create schema if not exists foo"},

		{MatchPattern: "^create table foo.bar", InTx: true},
		{MatchPattern: "^insert into foo.bar", InTx: true, SkipMatchArgs: true},

		{MatchPattern: "^select pg_advisory_unlock"},
	}
	if err := mock.SetStatements(wantStmts); err != nil {
		log.Fatalln(err)
	}

	db.Exec("select pg_advisory_lock(123)")
	// More database calls....
}
