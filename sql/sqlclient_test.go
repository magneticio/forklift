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

	createTableStatement := "CREATE TABLE \\?.\\? \\(`ID` int\\(11\\) NOT NULL AUTO_INCREMENT, `Record` mediumtext, PRIMARY KEY \\(`ID`\\) ENGINE=InnoDB DEFAULT CHARSET=utf8"

	mock.ExpectPrepare(createTableStatement).
		ExpectExec().
		WithArgs("organization", "organization").
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

	insertStatement := "INSERT INTO \\?.\\? VALUES\\( \\?, \\? \\)"

	mock.ExpectPrepare(insertStatement).
		ExpectExec().
		WithArgs("organization", "organization", 1, "just a test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	insertError := client.Insert("organization", 1, "just a test")

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

	queryStatement := "SELECT \\* FROM \\?.\\? WHERE ID = \\?"

	rows := sqlmock.NewRows([]string{"ID", "Record"}).
		AddRow(1, "just a test")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WithArgs("organization", "organization", 1).
		WillReturnRows(rows)

	result, queryError := client.FindById("organization", 1)

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

	queryStatement := "SELECT \\* FROM \\?.\\?"

	rows := sqlmock.NewRows([]string{"ID", "Record"}).
		AddRow(1, "just a test").
		AddRow(2, "just a test2").
		AddRow(3, "just a test3")

	mock.ExpectPrepare(queryStatement).
		ExpectQuery().
		WithArgs("organization", "organization").
		WillReturnRows(rows)

	result, queryError := client.List("organization")

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

	updateStatement := "UPDATE \\?.\\? SET `Record` = \\? WHERE ID = \\?"

	mock.ExpectPrepare(updateStatement).
		ExpectExec().
		WithArgs("organization", "organization", "just a test2", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updateError := client.Update("organization", 1, "just a test2")

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

	deleteStatement := "DELETE FROM \\?.\\? WHERE ID = \\?"

	mock.ExpectPrepare(deleteStatement).
		ExpectExec().
		WithArgs("organization", "organization", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	deleteError := client.Delete("organization", 1)

	assert.Nil(t, deleteError, fmt.Sprintf("Delete resulted in error %v \n", deleteError))

}
