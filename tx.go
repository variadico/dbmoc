package dbmoc

type Tx struct {
	state string
}

func (t *Tx) Commit() error {
	t.state = TxStateCommitted
	return nil
}

func (t *Tx) Rollback() error {
	t.state = TxStateRolledback
	return nil
}
