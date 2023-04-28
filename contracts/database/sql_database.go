package database

type (
	Iterator              func(ICursor) error
	Parameter             = interface{}
	Command               = string
	SqlTransactionHandler func(transaction ISqlTransaction) error

	ICursor interface {
		Scan(...Parameter) error
	}

	ISqlDatabase interface {
		Initialize() error
		GetName() string
		GetSchema() ISqlSchema
		RunScript(string, string) error
		Query(Iterator, Command, ...Parameter) error
		QuerySingle(Iterator, Command, ...Parameter) error
		Execute(Command, ...Parameter) (int64, error)
		ExecuteAtomic(ISqlTransaction, Command, ...Parameter) (int64, error)
		ExecuteBatch(Command, int64, ...Parameter) (int64, error)
		InsertSingle(Command, ...Parameter) error
		InsertSingleAtomic(ISqlTransaction, Command, ...Parameter) error
		UpdateSingle(Command, ...Parameter) error
		UpdateSingleAtomic(ISqlTransaction, Command, ...Parameter) error
		DeleteSingle(Command, ...Parameter) error
		DeleteSingleAtomic(ISqlTransaction, Command, ...Parameter) error
		InsertAll(Command, int64, ...Parameter) error
		Count(Command, ...Parameter) (int, error)
		WithTransaction(SqlTransactionHandler) error
	}

	ISqlTransaction interface {
		OnCommit(func())
	}

	ISqlSchema interface {
		GetTables() []string
		GetColumns(table string) []string
		GetTriggers() []string
		HasTable(table string) bool
		HasHistoryTable(table string) bool
		HasColumn(table, column string) bool
		HasTrigger(trigger string) bool
	}
)
