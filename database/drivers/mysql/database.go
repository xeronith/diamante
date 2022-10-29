package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/xeronith/diamante/contracts/database"
	. "github.com/xeronith/diamante/contracts/logging"
	. "github.com/xeronith/diamante/contracts/settings"
	"github.com/xeronith/diamante/contracts/system"
	. "github.com/xeronith/diamante/utility/collections"
)

var user, token string

type sqlDatabase struct {
	// TODO: Tune MySQL configuration (innodb_buffer_pool_size, ...)
	name             string
	connectionString string
}

func NewDatabase(configuration IConfiguration, logger ILogger, name string) ISqlDatabase {
	if !configuration.IsDockerized() {
		if configuration.IsTestEnvironment() {
			name = fmt.Sprintf("%s_test", name)
		} else if configuration.IsDevelopmentEnvironment() {
			name = fmt.Sprintf("%s_dev", name)
		} else if configuration.IsStagingEnvironment() {
			name = fmt.Sprintf("%s_staging", name)
		}
	}

	user = configuration.GetMySQLConfiguration().GetUsername()
	token = configuration.GetMySQLConfiguration().GetPassword()

	if configuration.GetMySQLConfiguration().IsPasswordSkipped() {
		token = ""
	}

	address := configuration.GetMySQLConfiguration().GetAddress()
	logger.SysComp(fmt.Sprintf("┄ Using MySQL(%s@%s/%s)", user, address, name))

	return &sqlDatabase{
		name: name,
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		connectionString: fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, token, address, name),
	}
}

func (database *sqlDatabase) GetName() string {
	return database.name
}

func (database *sqlDatabase) Initialize() error {
	command := "CREATE TABLE IF NOT EXISTS `__system__`(`id` BIGINT NOT NULL AUTO_INCREMENT, `script` VARCHAR(10240) NOT NULL, `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`)) ENGINE = InnoDB DEFAULT CHARSET = `utf8mb4` COLLATE = `utf8mb4_unicode_ci`;"
	db, err := sql.Open("mysql", database.connectionString)
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
	historyTables := make(map[string][]string)
	triggers := make([]string, 0)

	dbName := database.name
	historyDbName := fmt.Sprintf("%s_history", dbName)

	if err := database.Query(func(cursor ICursor) error {
		var table, column string
		if err := cursor.Scan(&table, &column); err != nil {
			return err
		}

		if _, exists := tables[table]; exists {
			tables[table] = append(tables[table], column)
		} else {
			tables[table] = []string{column}
		}

		return nil
	}, "SELECT `x`.`TABLE_NAME`, `y`.`COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`TABLES` AS `x` INNER JOIN `INFORMATION_SCHEMA`.`COLUMNS` AS `y` ON `x`.`TABLE_NAME` = `y`.`TABLE_NAME` WHERE `x`.`TABLE_SCHEMA` = ? AND `y`.`TABLE_SCHEMA` = ?;", dbName, dbName); err != nil {
		panic(err)
	}

	if err := database.Query(func(cursor ICursor) error {
		var table, column string
		if err := cursor.Scan(&table, &column); err != nil {
			return err
		}

		if _, exists := historyTables[table]; exists {
			historyTables[table] = append(historyTables[table], column)
		} else {
			historyTables[table] = []string{column}
		}

		return nil
	}, "SELECT `x`.`TABLE_NAME`, `y`.`COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`TABLES` AS `x` INNER JOIN `INFORMATION_SCHEMA`.`COLUMNS` AS `y` ON `x`.`TABLE_NAME` = `y`.`TABLE_NAME` WHERE `x`.`TABLE_SCHEMA` = ? AND `y`.`TABLE_SCHEMA` = ?;", historyDbName, historyDbName); err != nil {
		panic(err)
	}

	if err := database.Query(func(cursor ICursor) error {
		var trigger string
		if err := cursor.Scan(&trigger); err != nil {
			return err
		}

		triggers = append(triggers, trigger)
		return nil
	}, "SELECT `TRIGGER_NAME` FROM `INFORMATION_SCHEMA`.`TRIGGERS` WHERE `TRIGGER_SCHEMA` = ?;", dbName); err != nil {
		panic(err)
	}

	return &sqlSchema{
		tables:        tables,
		historyTables: historyTables,
		triggers:      triggers,
	}
}

func (database *sqlDatabase) RunScript(script string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s&multiStatements=true", database.connectionString))
	if err != nil {
		return err
	}

	defer func() { _ = db.Close() }()

	result, err := db.Query(script)
	if err != nil {
		return err
	}

	defer func() { _ = result.Close() }()

	return nil
}

func (database *sqlDatabase) Query(iterator Iterator, command Command, parameters ...Parameter) error {
	db, err := sql.Open("mysql", database.connectionString)
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
	db, err := sql.Open("mysql", database.connectionString)
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
	db, err := sql.Open("mysql", database.connectionString)
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

	db, err := sql.Open("mysql", database.connectionString)
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
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) InsertSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows != 1 {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) UpdateSingle(command Command, parameters ...Parameter) error {
	if affectedRows, err := database.Execute(command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) UpdateSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) DeleteSingle(command Command, parameters ...Parameter) error {
	if affectedRows, err := database.Execute(command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) DeleteSingleAtomic(transaction ISqlTransaction, command Command, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteAtomic(transaction, command, parameters...); err != nil {
		return err
	} else if affectedRows > 1 {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d}", command, affectedRows))
	}

	return nil
}

func (database *sqlDatabase) InsertAll(command Command, count int64, parameters ...Parameter) error {
	if affectedRows, err := database.ExecuteBatch(command, count, parameters...); err != nil {
		return err
	} else if affectedRows != count {
		return errors.New(fmt.Sprintf("affected_rows_inconsistency: '%s' {%d, %d}", command, affectedRows, count))
	}

	return nil
}

func (database *sqlDatabase) Count(command Command, parameters ...Parameter) (int, error) {
	db, err := sql.Open("mysql", database.connectionString)
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
	db, err := sql.Open("mysql", database.connectionString)
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
	tables        map[string][]string
	historyTables map[string][]string
	triggers      []string
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

func (schema *sqlSchema) GetHistoryTables() []string {
	tables := make([]string, 0)
	for table := range schema.historyTables {
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
		for _, column := range schema.tables[table] {
			columns = append(columns, column)
		}
	}

	return columns
}

func (schema *sqlSchema) GetTriggers() []string {
	triggers := make([]string, 0)
	for _, trigger := range schema.triggers {
		triggers = append(triggers, trigger)
	}

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
	for _, _table := range schema.GetHistoryTables() {
		if _table == table {
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
