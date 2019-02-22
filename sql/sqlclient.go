package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlClient interface {
	SetupOrganization() error
	Insert(organization string, id int, record string) error
	FindById(name string) (string, error)
	List(organization string) ([]Organization, error)
	Update(organization string, id int, record string) error
	Delete(organization string, id int) error
	Close() error
}

type MySqlClient struct {
	User     string
	Password string
	Host     string
	DbName   string
	Db       *sql.DB
}

type Organization struct {
	Id     int
	Record string
}

func NewMySqlClient(user string, password string, host string, dbName string) (*MySqlClient, error) {

	//"zEXmohRjfa:zaqT1JkXil@tcp(remotemysql.com)/zEXmohRjfa"

	db, connectionErr := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", host, ")/", dbName))
	if connectionErr != nil {
		fmt.Printf("Error: %v\n", connectionErr.Error())
		return nil, connectionErr
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

	stmtIns, tableErr := client.Db.Prepare("CREATE TABLE ?.? (`ID` int(11) NOT NULL AUTO_INCREMENT, `Record` mediumtext, PRIMARY KEY (`ID`) ENGINE=InnoDB DEFAULT CHARSET=utf8")
	if tableErr != nil {
		_ = tx.Rollback()
		fmt.Printf("Error: %v\n", tableErr.Error())
		return tableErr
	}

	defer stmtIns.Close()

	_, stmtInsErr := stmtIns.Exec(name, name)
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

func (client *MySqlClient) Insert(organization string, id int, record string) error {

	stmtIns, err := client.Db.Prepare("INSERT INTO ?.? VALUES( ?, ? )")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(organization, organization, id, record)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) FindById(organization string, id int) (*Organization, error) {

	stmtOut, err := client.Db.Prepare("SELECT * FROM ?.? WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	err = stmtOut.QueryRow(organization, organization, id).Scan(&resultId, &resultRecord)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	return &Organization{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) List(organization string) ([]Organization, error) {

	stmtOut, listErr := client.Db.Prepare("SELECT * FROM ?.?")
	if listErr != nil {
		fmt.Printf("Error: %v\n", listErr.Error())
		return []Organization{}, listErr
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	rows, err := stmtOut.Query(organization, organization)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return []Organization{}, err
	}
	defer rows.Close()

	var result []Organization

	for rows.Next() {
		err := rows.Scan(&resultId, &resultRecord)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return []Organization{}, err
		}

		result = append(result, Organization{
			Id:     resultId,
			Record: resultRecord,
		})

	}

	return result, nil
}

func (client *MySqlClient) Update(organization string, id int, record string) error {

	stmtIns, err := client.Db.Prepare("UPDATE ?.? SET `Record` = ? WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(organization, organization, record, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) Delete(organization string, id int) error {

	stmtIns, err := client.Db.Prepare("DELETE FROM ?.? WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(organization, organization, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}
