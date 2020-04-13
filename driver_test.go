package dbmoc

import (
	"database/sql"
	"testing"
)

func TestDriverOpenClose(t *testing.T) {
	t.Parallel()

	var driverFound bool
	for _, d := range sql.Drivers() {
		if d == "dbmoc" {
			driverFound = true
			break
		}
	}
	if !driverFound {
		sql.Register("dbmoc", &mockDriver{})
	}

	db, err := sql.Open("dbmoc", "")
	if err != nil {
		t.Error(err)
	}

	if err := db.Ping(); err != nil {
		t.Error(err)
	}

	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestDriverQuery(t *testing.T) {
	t.Parallel()

	t.Run("query row", func(t *testing.T) {
		mock := New()
		db := sql.OpenDB(mock)
		defer db.Close()

		err := mock.SetStatements([]Statement{
			{MatchPattern: `select 1`,
				Rows: &Rows{
					Cols: []string{"?column?"},
					Data: [][]interface{}{
						{1},
					},
				}},
		})
		if err != nil {
			t.Fatal(err)
		}

		row := db.QueryRow("select 1")

		var got int
		if err := row.Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != 1 {
			t.Error("unexpected scan")
			t.Fatalf("got=%d; want=%d", got, 1)
		}

		if got := mock.LenStatements(); got != 0 {
			t.Error("unexpected len statements")
			t.Fatalf("got=%d; want=%d", got, 0)
		}
	})

	t.Run("query rows", func(t *testing.T) {
		mock := New()
		db := sql.OpenDB(mock)
		defer db.Close()

		err := mock.SetStatements([]Statement{
			{MatchPattern: "select id, name from user",
				Rows: &Rows{
					Cols: []string{"id", "name"},
					Data: [][]interface{}{
						{1, "foo"},
						{2, "bar"},
						{3, "fizz"},
					},
				}},
		})
		if err != nil {
			t.Fatal(err)
		}

		type user struct {
			ID   int
			Name string
		}

		rows, err := db.Query("select id, name from user")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		var users []user
		for rows.Next() {
			var u user
			err := rows.Scan(&u.ID, &u.Name)
			if err != nil {
				t.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}

		if len(users) != 3 {
			t.Error("unexpected users")
			t.Fatalf("got=%d; want=%d", len(users), 3)
		}

		if got := mock.LenStatements(); got != 0 {
			t.Error("unexpected len statements")
			t.Fatalf("got=%d; want=%d", got, 0)
		}
	})
	t.Run("query rows with args", func(t *testing.T) {
		mock := New()
		db := sql.OpenDB(mock)
		defer db.Close()

		err := mock.SetStatements([]Statement{
			{MatchPattern: "select .*? from user", MatchArgs: []interface{}{int64(2)},
				Rows: &Rows{
					Cols: []string{"id"},
					Data: [][]interface{}{
						{1},
						{2},
						{3},
					},
				}},
		})
		if err != nil {
			t.Fatal(err)
		}

		rows, err := db.Query("select id from user where id > $1", 2)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var n int
			if err := rows.Scan(&n); err != nil {
				t.Fatal(err)
			}
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}

		if got := mock.LenStatements(); got != 0 {
			t.Error("unexpected len statements")
			t.Fatalf("got=%d; want=%d", got, 0)
		}
	})
}

func TestDriverExec(t *testing.T) {
	t.Parallel()

	t.Run("exec with args", func(t *testing.T) {
		mock := New()
		db := sql.OpenDB(mock)
		defer db.Close()

		err := mock.SetStatements([]Statement{
			{MatchPattern: "insert into [a-z]+ values",
				MatchArgs: []interface{}{int64(1), "foo"},
				Result: &Result{
					LastID:   3,
					Affected: 1,
				}},
		})
		if err != nil {
			t.Fatal(err)
		}

		res, err := db.Exec("insert into user values ($1, $2)", 1, "foo")
		if err != nil {
			t.Fatal(err)
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			t.Fatal(err)
		}
		if lastID != 3 {
			t.Error("unexpected last id")
			t.Fatalf("got=%d; want=%d", lastID, 3)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}
		if affected != 1 {
			t.Error("unexpected rows affected")
			t.Fatalf("got=%d; want=%d", affected, 1)
		}

		if got := mock.LenStatements(); got != 0 {
			t.Error("unexpected len statements")
			t.Fatalf("got=%d; want=%d", got, 0)
		}
	})
}
