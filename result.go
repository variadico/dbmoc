package dbmoc

type Result struct {
	LastID   int64
	Affected int64
}

func (r *Result) LastInsertId() (int64, error) {
	return r.LastID, nil
}

func (r *Result) RowsAffected() (int64, error) {
	return r.Affected, nil
}
