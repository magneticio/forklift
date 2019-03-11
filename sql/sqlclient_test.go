package sql

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSetupEvironment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createTableStatement := "CREATE TABLE IF NOT EXISTS `environment` \\(ID int\\(11\\) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY \\(ID\\)\\)"

	mock.
		ExpectExec(createTableStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertStatement := "INSERT INTO `environment` \\( Record \\) VALUES\\( \\? \\)"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test1 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test2 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test3 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	createError := client.SetupEnvironment("testdb", "environment", []string{"test1 value", "test2 value", "test3 value"})

	assert.Nil(t, createError, fmt.Sprintf("Create resulted in error %v \n", createError))

}

func TestUpdateEvironment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteFromTableStatement := "DELETE FROM `environment`"

	mock.ExpectBegin()

	mock.
		ExpectExec(deleteFromTableStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertStatement := "INSERT INTO `environment` \\( Record \\) VALUES\\( \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test1 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test2 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test3 value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	createError := client.UpdateEnvironment("testdb", "environment", []string{"test1 value", "test2 value", "test3 value"})

	assert.Nil(t, createError, fmt.Sprintf("Create resulted in error %v \n", createError))

}

func TestSetupOrganization(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	createSchemaStatement := "CREATE SCHEMA IF NOT EXISTS `testdb`"

	mock.
		ExpectExec(createSchemaStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createTableStatement := "CREATE TABLE IF NOT EXISTS `organization` \\(ID int\\(11\\) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY \\(ID\\)\\)"

	mock.
		ExpectExec(createTableStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createError := client.SetupOrganization("testdb", "organization")

	assert.Nil(t, createError, fmt.Sprintf("Create resulted in error %v \n", createError))

}

func TestInsert(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertStatement := "INSERT INTO `organization` \\( Record \\) VALUES\\( \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("just a test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertError := client.Insert("testdb", "organization", "just a test")

	assert.Nil(t, insertError, fmt.Sprintf("Insert resulted in error %v \n", insertError))

}

func TestInsertOrReplace(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectBegin()

	deleteStatement := "DELETE FROM `organization` WHERE Record LIKE '%\"name\":\"name\"%' AND Record LIKE '%\"kind\":\"kind\"%'"

	mock.ExpectPrepare(deleteStatement).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertStatement := "INSERT INTO `organization` \\( Record \\) VALUES\\( \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("record").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	insertError := client.InsertOrReplace("testdb", "organization", "name", "kind", "record")

	assert.Nil(t, insertError, fmt.Sprintf("Insert resulted in error %v \n", insertError))

}

func TestFindByNameAndKind(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	queryStatement := "SELECT \\* FROM `testtable` WHERE Record LIKE '%\"name\":\"recordname\"%' AND Record LIKE '%\"kind\":\"recordkind\"%'"

	rows := sqlmock.NewRows([]string{"ID", "Record"}).
		AddRow(1, "just a test")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WillReturnRows(rows)

	result, queryError := client.FindByNameAndKind("testdb", "testtable", "recordname", "recordkind")

	expected := &Row{
		Id:     1,
		Record: "just a test",
	}

	assert.Nil(t, queryError, fmt.Sprintf("Query resulted in error %v \n", queryError))
	assert.Equal(t, result, expected, fmt.Sprintf("Query returned wrong result %v \n", result))

}

func TestFindById(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	queryStatement := "SELECT \\* FROM `organization` WHERE ID = \\?"

	rows := sqlmock.NewRows([]string{"ID", "Record"}).
		AddRow(1, "just a test")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(rows)

	result, queryError := client.FindById("testdb", "organization", 1)

	expected := &Row{
		Id:     1,
		Record: "just a test",
	}

	assert.Nil(t, queryError, fmt.Sprintf("Query resulted in error %v \n", queryError))
	assert.Equal(t, result, expected, fmt.Sprintf("Query returned wrong result %v \n", result))

}

func TestList(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	queryStatement := "SELECT \\* FROM `organization`"

	rows := sqlmock.NewRows([]string{"ID", "Record"}).
		AddRow(1, "just a test").
		AddRow(2, "just a test2").
		AddRow(3, "just a test3")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WillReturnRows(rows)

	result, queryError := client.List("testdb", "organization")

	expected := []Row{
		Row{
			Id:     1,
			Record: "just a test",
		},
		Row{
			Id:     2,
			Record: "just a test2",
		},
		Row{
			Id:     3,
			Record: "just a test3",
		},
	}

	assert.Nil(t, queryError, fmt.Sprintf("Query resulted in error %v \n", queryError))
	assert.Equal(t, result, expected, fmt.Sprintf("Query returned wrong result %v \n", result))

}

func TestUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateStatement := "UPDATE `organization` SET `Record` = \\? WHERE ID = \\?"

	mock.ExpectPrepare(updateStatement).
		ExpectExec().
		WithArgs("just a test2", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateError := client.Update("testdb", "organization", 1, "just a test2")

	assert.Nil(t, updateError, fmt.Sprintf("Update resulted in error %v \n", updateError))

}

func TestDelete(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteStatement := "DELETE FROM `organization` WHERE ID = \\?"

	mock.ExpectPrepare(deleteStatement).
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteError := client.Delete("testdb", "organization", 1)

	assert.Nil(t, deleteError, fmt.Sprintf("Delete resulted in error %v \n", deleteError))

}

func TestDeleteByNameAndKind(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := MySqlClient{
		User:     "",
		Password: "",
		Host:     "",
		DbName:   "",
		Db:       db,
	}

	useDbStatement := "USE `testdb`"

	mock.
		ExpectExec(useDbStatement).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteStatement := "DELETE FROM `organization` WHERE Record LIKE '%\"name\":\"name\"%' AND Record LIKE '%\"kind\":\"users\"%'"

	mock.ExpectPrepare(deleteStatement).
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteError := client.DeleteByNameAndKind("testdb", "organization", "name", "users")

	assert.Nil(t, deleteError, fmt.Sprintf("Delete resulted in error %v \n", deleteError))

}
