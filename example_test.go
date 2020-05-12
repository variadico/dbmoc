package dbmoc

import (
	"database/sql"
	"log"
)

func ExampleMock_query() {
	mock := New()
	db := sql.OpenDB(mock)

	rows := []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}{
		{ID: 1, Name: "fizz"},
		{ID: 2, Name: "buzz"},
	}
	wantStmts := []Statement{
		{
			MatchQuery: `^select .*\s+from foo.bar`,
			MatchArgs:    []interface{}{3},

			Rows: NewRows(rows),
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
		{MatchQuery: "^select pg_advisory_lock"},
		{MatchQuery: "^create schema if not exists foo"},

		{MatchQuery: "^create table foo.bar", InTx: true},
		{MatchQuery: "^insert into foo.bar", InTx: true, SkipMatchArgs: true},

		{MatchQuery: "^select pg_advisory_unlock"},
	}
	if err := mock.SetStatements(wantStmts); err != nil {
		log.Fatalln(err)
	}

	db.Exec("select pg_advisory_lock(123)")
	// More database calls....
}
