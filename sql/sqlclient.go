package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlClient interface {
	SetupOrganisation() error
	Insert(name string, description string) error
	Query(name string) (string, error)
	Update(name string, description string) error
	Delete(name string) error
	Close() error
}

type MySqlClient struct {
	User     string
	Password string
	Host     string
	DbName   string
	Db       *sql.DB
}

func NewMySqlClient(user string, password string, host string, dbName string) (*MySqlClient, error) {

	//"zEXmohRjfa:zaqT1JkXil@tcp(remotemysql.com)/zEXmohRjfa"

	db, connectionErr := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", host, ")/", dbName))
	if connectionErr != nil {
		fmt.Printf("Error: %v\n", connectionErr.Error())
		return nil, connectionErr
	}

	c := context.Background()

	if c == nil {
		return nil, errors.New("SQL: Context background can not be initialized.")
	}

	return &MySqlClient{
		User:     user,
		Password: password,
		Host:     host,
		DbName:   dbName,
		Db:       db,
	}, nil
}

func (client *MySqlClient) Close() error {

	err := client.Db.Close()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) SetupOrganization(name string) error {

	tx, txErr := client.Db.Begin()
	if txErr != nil {
		fmt.Printf("Error: %v\n", txErr.Error())
		return txErr
	}

	stmtSchema, schemaErr := client.Db.Prepare("CREATE SCHEMA IF NOT EXISTS ?")
	if schemaErr != nil {
		_ = tx.Rollback()
		fmt.Printf("Error: %v\n", schemaErr.Error())
		return schemaErr
	}

	defer stmtSchema.Close()

	_, stmtSchemaErr := stmtSchema.Exec(name)
	if stmtSchemaErr != nil {
		_ = tx.Rollback()
		fmt.Printf("Error: %v\n", stmtSchemaErr.Error())
		return stmtSchemaErr
	}

	stmtIns, tableErr := client.Db.Prepare("CREATE TABLE ? (`ID` int(11) NOT NULL AUTO_INCREMENT, `Record` mediumtext, PRIMARY KEY (`ID`) ENGINE=InnoDB DEFAULT CHARSET=utf8")
	if tableErr != nil {
		_ = tx.Rollback()
		fmt.Printf("Error: %v\n", tableErr.Error())
		return tableErr
	}

	defer stmtIns.Close()

	_, stmtInsErr := stmtIns.Exec(name)
	if stmtInsErr != nil {
		fmt.Printf("Error: %v\n", stmtInsErr.Error())
		return stmtInsErr
	}

	if commitError := tx.Commit(); commitError != nil {
		fmt.Printf("Error: %v\n", commitError.Error())
		return commitError
	}

	return nil
}

func (client *MySqlClient) Insert(name string, description string) error {

	stmtIns, err := client.Db.Prepare("INSERT INTO ENVIRONMENT VALUES( ?, ? )")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(name, description)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) Query(name string) (string, error) {

	stmtOut, err := client.Db.Prepare("SELECT DESCRIPTION FROM ENVIRONMENT WHERE NAME = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return "", err
	}

	defer stmtOut.Close()

	var result string

	err = stmtOut.QueryRow(name).Scan(&result)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return "", err
	}

	return result, nil
}

func (client *MySqlClient) Update(name string, description string) error {

	stmtIns, err := client.Db.Prepare("UPDATE ENVIRONMENT SET DESCRIPTION = ? WHERE NAME = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(description, name)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) Delete(name string) error {

	stmtIns, err := client.Db.Prepare("DELETE FROM ENVIRONMENT WHERE NAME = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}
