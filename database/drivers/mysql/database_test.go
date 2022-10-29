package mysql_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/xeronith/diamante/contracts/database"
	"github.com/xeronith/diamante/database/drivers/mysql"
	"github.com/xeronith/diamante/logging"
	"github.com/xeronith/diamante/settings"
)

var (
	configuration = settings.NewTestConfiguration()
	logger        = logging.GetDefaultLogger()
)

func TestDatabase_InsertAll(test *testing.T) {
	database := mysql.NewDatabase(configuration, logger, "sandbox")

	args := make([]interface{}, 0)

	total := int64(10)
	for i := int64(0); i < total; i++ {
		args = append(args, time.Now().UnixNano())
	}

	err := database.InsertAll("INSERT INTO `items` (`title`) VALUES (?);", total, args...)
	if err != nil {
		test.Fatal(err)
	}
}

func TestDatabase_GetSchema(test *testing.T) {
	database := mysql.NewDatabase(configuration, logger, "suppline")

	var schema ISqlSchema
	if schema = database.GetSchema(); schema == nil {
		test.FailNow()
	}

	if schema.HasTable("invalid_table") {
		test.FailNow()
	}

	if !schema.HasTable("identities") {
		test.FailNow()
	}

	if !schema.HasTrigger("identities_after_update") {
		test.FailNow()
	}

	fmt.Println(schema)
}
