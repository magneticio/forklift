package sql

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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

	createSchemaStatement := "CREATE SCHEMA IF NOT EXISTS \\?"

	mock.ExpectBegin()

	mock.ExpectPrepare(createSchemaStatement).
		ExpectExec().
		WithArgs("organization").
		WillReturnResult(sqlmock.NewResult(1, 1))

	createTableStatement := "CREATE TABLE \\? \\(`ID` int\\(11\\) NOT NULL AUTO_INCREMENT, `Record` mediumtext, PRIMARY KEY \\(`ID`\\) ENGINE=InnoDB DEFAULT CHARSET=utf8"

	mock.ExpectPrepare(createTableStatement).
		ExpectExec().
		WithArgs("organization").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	createError := client.SetupOrganization("organization")

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

	insertStatement := "INSERT INTO ENVIRONMENT VALUES\\( \\?, \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("test", "just a test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertError := client.Insert("test", "just a test")

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

	queryStatement := "SELECT DESCRIPTION FROM ENVIRONMENT WHERE NAME = \\?"

	rows := sqlmock.NewRows([]string{"DESCRIPTION"}).
		AddRow("just a test")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WithArgs("test").
		WillReturnRows(rows)

	result, queryError := client.Query("test")

	assert.Nil(t, queryError, fmt.Sprintf("Query resulted in error %v \n", queryError))
	assert.Equal(t, result, "just a test", fmt.Sprintf("Query returned wrong result %v \n", result))

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

	updateStatement := "UPDATE ENVIRONMENT SET DESCRIPTION = \\? WHERE NAME = \\?"

	mock.ExpectPrepare(updateStatement).
		ExpectExec().
		WithArgs("just a test2", "test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateError := client.Update("test", "just a test2")

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

	deleteStatement := "DELETE FROM ENVIRONMENT WHERE NAME = ?"

	mock.ExpectPrepare(deleteStatement).
		ExpectExec().
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteError := client.Delete("test")

	assert.Nil(t, deleteError, fmt.Sprintf("Delete resulted in error %v \n", deleteError))

}
