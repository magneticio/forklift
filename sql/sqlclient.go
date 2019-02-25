package sql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/util"
)

type SqlClient interface {
	SetupOrganization(dbName string, tableName string) error
	Insert(dbName string, tableName string, id int, record string) error
	FindById(dbName string, tableName string, id int) (*Organization, error)
	List(dbName string, tableName string) ([]Organization, error)
	Update(dbName string, tableName string, id int, record string) error
	Delete(dbName string, tableName string, id int) error
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

func NewSqlClient(config models.Database) (SqlClient, error) {
	if config.Type == "mysql" {
		// TODO: add params
		fmt.Printf("host: %v\n", config.Sql.Url)
		host, hostError := util.GetHostFromUrl(config.Sql.Url)
		if hostError != nil {
			return nil, hostError
		}

		sqlClient, mySqlError := NewMySqlClient(config.Sql.User, config.Sql.Password, host, config.Sql.Database)
		if mySqlError != nil {
			return nil, mySqlError
		}

		return sqlClient, nil
	}
	return nil, errors.New("Unsupported Key Value Store Client: " + config.Type)
}

func NewMySqlClient(user string, password string, host string, dbName string) (*MySqlClient, error) {

	fmt.Printf("%v\n", fmt.Sprint(user, ":", password, "@tcp(", host, ")/"))

	db, connectionErr := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", host, ")/"))
	if connectionErr != nil {
		fmt.Printf("Error in mysql client creation: %v\n", connectionErr.Error())
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
		fmt.Printf("Error in client close: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) SetupOrganization(dbName string, tableName string) error {
	fmt.Printf("dbName %v tableName %v\n", dbName, tableName)
	_, createSchemaErr := client.Db.Exec("CREATE SCHEMA IF NOT EXISTS `" + dbName + "`")
	if createSchemaErr != nil {
		fmt.Printf("Error in Setup Organization: %v\n", createSchemaErr.Error())
		return createSchemaErr
	}

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		_, dropError := client.Db.Exec("DROP SCHEMA `" + dbName + "`")
		if dropError != nil {
			fmt.Printf("Error in table drop: %v\n", dropError.Error())
			return dropError
		}
		return useDbErr
	}

	_, insertErr := client.Db.Exec("CREATE TABLE IF NOT EXISTS `" + tableName + "` (ID int(11) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY (ID))")
	if insertErr != nil {
		fmt.Printf("Error during create: %v\n", insertErr.Error())
		_, dropError := client.Db.Exec("DROP SCHEMA `" + dbName + "`")
		if dropError != nil {
			fmt.Printf("Error: %v\n", dropError.Error())
			return dropError
		}
		return insertErr
	}

	return nil
}

func (client *MySqlClient) Insert(dbName string, tableName string, id int, record string) error {

	_, useDbErr := client.Db.Exec("USE " + dbName)
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("INSERT INTO " + tableName + " VALUES( ?, ? )")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(id, record)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) FindById(dbName string, tableName string, id int) (*Organization, error) {

	_, useDbErr := client.Db.Exec("USE " + dbName)
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return nil, useDbErr
	}

	stmtOut, err := client.Db.Prepare("SELECT * FROM " + tableName + " WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	err = stmtOut.QueryRow(id).Scan(&resultId, &resultRecord)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	return &Organization{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) List(dbName string, tableName string) ([]Organization, error) {

	_, useDbErr := client.Db.Exec("USE " + dbName)
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return []Organization{}, useDbErr
	}

	stmtOut, listErr := client.Db.Prepare("SELECT * FROM " + tableName)
	if listErr != nil {
		fmt.Printf("Error: %v\n", listErr.Error())
		return []Organization{}, listErr
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	rows, err := stmtOut.Query()
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

func (client *MySqlClient) Update(dbName string, tableName string, id int, record string) error {

	_, useDbErr := client.Db.Exec("USE " + dbName)
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("UPDATE " + tableName + " SET `Record` = ? WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(record, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) Delete(dbName string, tableName string, id int) error {

	_, useDbErr := client.Db.Exec("USE " + dbName)
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("DELETE FROM " + tableName + " WHERE ID = ?")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}
