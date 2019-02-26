package sql

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Should be used with a local database
// func TestSetupOrganizationOnLocal(t *testing.T) {
//
// 	client, _ := NewMySqlClient("root", "root", "localhost", "organization")
//
// 	createError := client.SetupOrganization("testdb", "organization")
// 	assert.Nil(t, createError, fmt.Sprintf("Create resulted in error %v \n", createError))
//
// 	insertError := client.Insert("testdb", "organization", 1, "ciao")
// 	assert.Nil(t, insertError, fmt.Sprintf("Insert resulted in error %v \n", createError))
//
// 	insertError2 := client.Insert("testdb", "organization", 2, "ciao2")
// 	assert.Nil(t, insertError2, fmt.Sprintf("Insert resulted in error %v \n", createError))
//
// 	_, findError := client.FindById("testdb", "organization", 1)
// 	assert.Nil(t, findError, fmt.Sprintf("Find resulted in error %v \n", findError))
//
// 	_, listError := client.List("testdb", "organization")
// 	assert.Nil(t, listError, fmt.Sprintf("List resulted in error %v \n", createError))
//
// 	deleteError := client.Delete("testdb", "organization", 1)
// 	assert.Nil(t, deleteError, fmt.Sprintf("Delete resulted in error %v \n", createError))
//
// 	deleteError2 := client.Delete("testdb", "organization", 2)
// 	assert.Nil(t, deleteError2, fmt.Sprintf("Delete resulted in error %v \n", deleteError2))
//
// }

// func TestConnection(t *testing.T) {
//
// 	_, connectionErr := sql.Open("mysql", "root:secret@tcp(api.dev.vamp.merapar.net:32401)/")
//
// 	assert.Nil(t, connectionErr, fmt.Sprintf("Could not connect due to %v \n", connectionErr))
//
// }

// func TestRemoteInsert(t *testing.T) {
//
// 	client, _ := NewMySqlClient("root", "secret", "api.dev.vamp.merapar.net:32401", "")
//
// 	insertError := client.Insert("vamp-neworg", "neworg", 2, "ciao-test")
//
// 	assert.Nil(t, insertError, fmt.Sprintf("Insert resulted in error %v \n", insertError))
//
// }

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

	insertStatement := "INSERT INTO `organization` VALUES\\( \\?, \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs(1, "just a test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertError := client.Insert("testdb", "organization", 1, "just a test")

	assert.Nil(t, insertError, fmt.Sprintf("Insert resulted in error %v \n", insertError))

}

func TestQuery(t *testing.T) {

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

	expected := &Organization{
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

	expected := []Organization{
		Organization{
			Id:     1,
			Record: "just a test",
		},
		Organization{
			Id:     2,
			Record: "just a test2",
		},
		Organization{
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
