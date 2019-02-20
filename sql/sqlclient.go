package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlClient interface {
	CreateTable() error
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

	db, err := sql.Open("mysql", "zEXmohRjfa:zaqT1JkXil@tcp(remotemysql.com)/zEXmohRjfa")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	return &MySqlClient{
		User:     user,
		Password: password,
		Host:     host,
		DbName:   dbName,
		Db:       db,
	}, nil
}

func (client MySqlClient) Close() error {

	err := client.Db.Close()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client MySqlClient) CreateTable() error {

	stmtIns, err := client.Db.Prepare("CREATE TABLE IF NOT EXISTS ENVIRONMENT (NAME TEXT NULL, DESCRIPTION TEXT NULL)") // ? = placeholder
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client MySqlClient) Insert(name string, description string) error {

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

func (client MySqlClient) Query(name string) (string, error) {

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

func (client MySqlClient) Update(name string, description string) error {

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

func (client MySqlClient) Delete(name string) error {

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
