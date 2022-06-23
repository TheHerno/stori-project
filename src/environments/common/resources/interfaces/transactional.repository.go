package interfaces

/*
ITransactionalRepository handle transactions on repositories
*/
type ITransactionalRepository interface {
	Begin(initialTx interface{}) interface{}
	Commit() error
	Rollback() error
	SavePoint() error
	Clone() interface{}
}
