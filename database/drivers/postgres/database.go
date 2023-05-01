package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	. "github.com/xeronith/diamante/contracts/database"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/settings"
	"github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

var user, password string

type sqlDatabase struct {
	name             string
	connectionString string
}

func NewDatabase(configuration IConfiguration, logger ILogger, dbname string) ISqlDatabase {
	if configuration.IsTestEnvironment() {
		dbname = fmt.Sprintf("%s_test", dbname)
	} else if configuration.IsDevelopmentEnvironment() {
		dbname = fmt.Sprintf("%s_dev", dbname)
	} else if configuration.IsStagingEnvironment() {
		dbname = fmt.Sprintf("%s_staging", dbname)
	}

	host := configuration.GetPostgreSQLConfiguration().GetHost()
	port := configuration.GetPostgreSQLConfiguration().GetPort()
	user = configuration.GetPostgreSQLConfiguration().GetUsername()
	password = configuration.GetPostgreSQLConfiguration().GetPassword()

	logger.SysComp(fmt.Sprintf("┄ Using PostgreSQL(%s@%s:%s/%s)", user, host, port, dbname))

	return &sqlDatabase{
		name:             dbname,
		connectionString: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname),
	}
}

func (database *sqlDatabase) GetName() string {
	return database.name
}

func (database *sqlDatabase) Initialize() error {
	command := `CREATE TABLE IF NOT EXISTS "__system__"("id" BIGSERIAL NOT NULL, "script" VARCHAR(10240) NOT NULL, "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));`
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	result, err := db.Query(command)
	if err != nil {
		return err
	}

	_ = result
	return nil
}

func (database *sqlDatabase) GetSchema() ISqlSchema {
	tables := make(map[string][]string)
	triggers := make([]string, 0)

	dbName := database.name

	if err := database.Query(func(cursor ICursor) error {
		var table, column string
		if err := cursor.Scan(&table, &column); err != nil {
			return err
		}

		tables[table] = append(tables[table], column)

		return nil
	}, `SELECT "x"."table_name", "y"."column_name" FROM "information_schema"."tables" AS "x" INNER JOIN "information_schema"."columns" AS "y" ON "x"."table_name" = "y"."table_name" WHERE "x"."table_catalog" = $1 AND "y"."table_catalog" = $2 AND "x"."table_schema" = 'public' AND "y"."table_schema" = 'public';`, dbName, dbName); err != nil {
		panic(err)
	}

	if err := database.Query(func(cursor ICursor) error {
		var trigger string
		if err := cursor.Scan(&trigger); err != nil {
			return err
		}

		triggers = append(triggers, trigger)
		return nil
	}, `SELECT "trigger_name" FROM "information_schema"."triggers" WHERE "trigger_catalog" = $1 AND "trigger_schema" = 'public';`, dbName); err != nil {
		panic(err)
	}

	return &sqlSchema{
		tables:   tables,
		triggers: triggers,
	}
}

func (database *sqlDatabase) RunScript(script string, separator string) error {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, statement := range strings.Split(script, separator) {
		if strings.TrimSpace(statement) != "" {
			if _, err := tx.Exec(statement); err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}

				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (database *sqlDatabase) Query(iterator Iterator, command Command, parameters ...Parameter) error {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	result, err := db.Query(command, parameters...)
	if err != nil {
		return err
	}

	defer func() { _ = result.Close() }()

	if iterator != nil {
		for result.Next() {
			if err := iterator(result); err != nil {
				return err
			}
		}
	}

	return nil
}

func (database *sqlDatabase) QuerySingle(iterator Iterator, command Command, parameters ...Parameter) error {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	result, err := db.Query(command, parameters...)
	if err != nil {
		return err
	}

	defer func() { _ = result.Close() }()

	if iterator != nil {
		if result.Next() {
			if err := iterator(result); err != nil {
				return err
			}
		} else {
			return errors.New("not_found")
		}
	}

	return nil
}

func (database *sqlDatabase) Execute(command Command, parameters ...Parameter) (int64, error) {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return 0, err
	}

	defer func() { _ = db.Close() }()

	result, err := db.Exec(command, parameters...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (database *sqlDatabase) ExecuteAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) (int64, error) {
	if sqlTransaction, ok := transaction.(*sqlTransaction); ok {
		result, err := sqlTransaction.databaseTransaction.Exec(command, parameters...)
		if err != nil {
			return 0, err
		}

		return result.RowsAffected()
	}

	return 0, errors.New("transaction_not_valid")
}

func (database *sqlDatabase) ExecuteBatch(command Command, count int64, parameters ...Parameter) (int64, error) {
	if count == 0 {
		return 0, nil
	}

	parametersCount := int64(len(parameters)) / count
	if int64(len(parameters))%count > 0 {
		parametersCount++
	}

	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return 0, err
	}

	defer func() { _ = db.Close() }()

	transaction, err := db.Begin()
	if err != nil {
		return 0, err
	}

	statement, err := transaction.Prepare(command)
	if err != nil {
		return 0, err
	}

	defer func() { _ = statement.Close() }()

	var lastError error
	total := int64(0)
	for i := int64(0); i < count; i++ {
		offset := parametersCount * i
		result, err := statement.Exec(parameters[offset : offset+parametersCount]...)
		if err != nil {
			lastError = err
			break
		}

		affectedRowsCount, err := result.RowsAffected()
		if err != nil {
			lastError = err
			break
		}

		total += affectedRowsCount
	}

	if lastError != nil {
		if err := transaction.Rollback(); err != nil {
			// TODO: Where is your god now?
			_ = err
		}

		return 0, lastError
	}

	if err := transaction.Commit(); err != nil {
		return 0, nil
	}

	return total, nil
}

func (database *sqlDatabase) InsertSingle(command Command, parameters ...Parameter) error {
	if affectedRows, err := database.Execute(command, parameters...); err != nil {
		return err
	} else if affectedRows != 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) InsertSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows != 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) UpdateSingle(command Command, parameters ...Parameter) error {
	if affectedRows, err := database.Execute(command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) UpdateSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) DeleteSingle(command Command, parameters ...Parameter) error {
	if affectedRows, err := database.Execute(command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) DeleteSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows)
	}

	return nil
}

func (database *sqlDatabase) InsertAll(command Command, count int64, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteBatch(command, count, parameters...); err != nil {
		return err
	} else if affectedRows != count {
		return fmt.Errorf("affected_rows_inconsistency: '%s' {%d, %d}", command, affectedRows, count)
	}

	return nil
}

func (database *sqlDatabase) Count(command Command, parameters ...Parameter) (int, error) {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return 0, err
	}

	defer func() { _ = db.Close() }()

	count := 0
	if err := db.QueryRow(command, parameters...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (database *sqlDatabase) WithTransaction(handler SqlTransactionHandler) (err error) {
	db, err := sql.Open("postgres", database.connectionString)
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlTransaction := NewTransaction(tx).(*sqlTransaction)

	defer func() {
		if reason := recover(); reason != nil {
			sqlTransaction.Rollback()
			panic(reason)
		} else if err != nil {
			sqlTransaction.Rollback()
		} else {
			err = sqlTransaction.Commit()
		}
	}()

	err = handler(sqlTransaction)
	return err
}

type sqlTransaction struct {
	databaseTransaction *sql.Tx
	callbacks           ISlice
}

func NewTransaction(databaseTransaction *sql.Tx) ISqlTransaction {
	return &sqlTransaction{
		databaseTransaction: databaseTransaction,
		callbacks:           NewConcurrentSlice(),
	}
}

func (transaction *sqlTransaction) OnCommit(callback func()) {
	if callback != nil {
		transaction.callbacks.Append(callback)
	}
}

func (transaction *sqlTransaction) Commit() error {
	if err := transaction.databaseTransaction.Commit(); err != nil {
		return err
	}

	transaction.callbacks.ForEach(func(_ int, object system.ISystemObject) {
		if object != nil {
			object.(func())()
		}
	})

	return nil
}

func (transaction *sqlTransaction) Rollback() {
	_ = transaction.databaseTransaction.Rollback()
}

type sqlSchema struct {
	tables   map[string][]string
	triggers []string
}

func (schema *sqlSchema) GetTables() []string {
	tables := make([]string, 0)
	for table := range schema.tables {
		tables = append(tables, table)
	}

	sort.Slice(tables, func(x, y int) bool {
		return tables[x] < tables[y]
	})

	return tables
}

func (schema *sqlSchema) GetColumns(table string) []string {
	columns := make([]string, 0)
	if _, exists := schema.tables[table]; exists {
		columns = append(columns, schema.tables[table]...)
	}

	return columns
}

func (schema *sqlSchema) GetTriggers() []string {
	triggers := make([]string, 0)
	triggers = append(triggers, schema.triggers...)

	sort.Slice(triggers, func(x, y int) bool {
		return triggers[x] < triggers[y]
	})

	return triggers
}

func (schema *sqlSchema) HasTable(table string) bool {
	for _, _table := range schema.GetTables() {
		if _table == table {
			return true
		}
	}

	return false
}

func (schema *sqlSchema) HasHistoryTable(table string) bool {
	historyTable := fmt.Sprintf("%s_history", table)
	for _, _table := range schema.GetTables() {
		if _table == historyTable {
			return true
		}
	}

	return false
}

func (schema *sqlSchema) HasColumn(table, column string) bool {
	for _, _column := range schema.GetColumns(table) {
		if _column == column {
			return true
		}
	}

	return false
}

func (schema *sqlSchema) HasTrigger(trigger string) bool {
	for _, _trigger := range schema.GetTriggers() {
		if _trigger == trigger {
			return true
		}
	}

	return false
}

func (schema *sqlSchema) String() string {
	result := ""
	for _, table := range schema.GetTables() {
		result += "- "
		result += table
		result += ":\n"
		for _, column := range schema.GetColumns(table) {
			result += "\t- "
			result += column
			result += "\n"
		}

		result += "\n"
	}

	for _, trigger := range schema.GetTriggers() {
		result += "- "
		result += trigger
		result += "\n"
	}

	return result
}
