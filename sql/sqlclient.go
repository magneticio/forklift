package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/util"
)

type SqlClient interface {
	SetupOrganization(dbName string, tableName string) error
	SetupEnvironment(dbName string, tableName string, elements []string) error
	UpdateEnvironment(dbName string, tableName string, elements []string) error
	Insert(dbName string, tableName string, record string) error
	InsertOrReplace(dbName string, tableName string, name string, kind string, record string) error
	FindById(dbName string, tableName string, id int) (*Row, error)
	FindByNameAndKind(dbName string, tableName string, name string, kind string) (*Row, error)
	List(dbName string, tableName string, kind string) ([]Row, error)
	Update(dbName string, tableName string, id int, record string) error
	Delete(dbName string, tableName string, id int) error
	DeleteByNameAndKind(dbName string, tableName string, name string, kind string) error
	Close() error
}

type MySqlClient struct {
	User     string
	Password string
	Host     string
	DbName   string
	Db       *sql.DB
}

type Row struct {
	Id     int
	Record string
}

func NewSqlClient(config models.Database) (SqlClient, error) {
	if config.Type == "mysql" {
		// TODO: add params
		host, hostError := util.GetHostFromUrl(strings.TrimPrefix(config.Sql.Url, "jdbc:"))
		if hostError != nil {
			return nil, hostError
		}

		logging.Info.Log("Creating new sql client. User: %v - Host: %v - Database: %v\n", config.Sql.User, host, config.Sql.Database)

		sqlClient, mySqlError := NewMySqlClient(config.Sql.User, config.Sql.Password, host, config.Sql.Database)
		if mySqlError != nil {
			return nil, mySqlError
		}

		return sqlClient, nil
	}
	return nil, errors.New("Unsupported Sql Client: " + config.Type)
}

func NewMySqlClient(user string, password string, host string, dbName string) (*MySqlClient, error) {

	// fmt.Printf("%v\n", fmt.Sprint(user, ":", password, "@tcp(", host, ")/"))

	db, connectionErr := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", host, ")/"))
	if connectionErr != nil {
		logging.Error.Log("Error in mysql client creation: %v\n", connectionErr.Error())
		return nil, connectionErr
	}

	logging.Info.Log("Created new mysql client")

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
		logging.Error.Log("Error in client close: %v\n", err.Error())
		return err
	}

	logging.Info.Log("MySql Client was closed")

	return nil
}

func (client *MySqlClient) SetupOrganization(dbName string, tableName string) error {

	logging.Info.Log("Creating organization database %v\n", dbName)

	_, createSchemaErr := client.Db.Exec("CREATE SCHEMA IF NOT EXISTS `" + dbName + "`")
	if createSchemaErr != nil {
		logging.Error.Log("Error while creating organization database %v - %v\n", dbName, createSchemaErr.Error())
		return createSchemaErr
	}

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		_, dropError := client.Db.Exec("DROP SCHEMA `" + dbName + "`")
		if dropError != nil {
			logging.Error.Log("Error during rollback: %v\n", dropError.Error())
			return dropError
		}
		return useDbErr
	}

	logging.Info.Log("Creating organization table %v\n", tableName)

	_, insertErr := client.Db.Exec("CREATE TABLE IF NOT EXISTS `" + tableName + "` (ID int(11) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY (ID))")
	if insertErr != nil {
		logging.Error.Log("Error while creating organization table %v - %v\n", tableName, insertErr.Error())
		_, dropError := client.Db.Exec("DROP SCHEMA `" + dbName + "`")
		if dropError != nil {
			logging.Error.Log("Error during rollback: %v\n", dropError.Error())
			return dropError
		}
		return insertErr
	}

	return nil
}

func (client *MySqlClient) SetupEnvironment(dbName string, tableName string, elements []string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Creating environment table %v\n", tableName)

	_, createErr := client.Db.Exec("CREATE TABLE IF NOT EXISTS `" + tableName + "` (ID int(11) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY (ID))")
	if createErr != nil {
		fmt.Printf("Error while creating environment table %v %v\n", tableName, createErr.Error())
		return createErr
	}

	for index, value := range elements {

		logging.Info.Log("Inserting artifact with index %v in environment table %v", index, tableName)

		insertErr := client.Insert(dbName, tableName, value)
		if insertErr != nil {
			logging.Error.Log("Error while inserting in table %v - %v\n", tableName, createErr.Error())
			_, dropError := client.Db.Exec("DROP TABLE `" + tableName + "`")
			if dropError != nil {
				logging.Error.Log("Error during rollback: %v\n", dropError.Error())
				return dropError
			}
			return insertErr
		}

	}

	return nil
}

func (client *MySqlClient) UpdateEnvironment(dbName string, tableName string, elements []string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Starting transaction %v\n", dbName)

	tx, startTransactionError := client.Db.Begin()
	if startTransactionError != nil {
		logging.Error.Log("Error starting transaction - %var name type\n", startTransactionError.Error())
		return startTransactionError
	}

	logging.Info.Log("Deleting artifacts from environment %v\n", tableName)

	_, deleteErr := tx.Exec("DELETE FROM `" + tableName + "`")
	if deleteErr != nil {
		logging.Error.Log("Error while deleting artifacts: %v\n", deleteErr.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			logging.Error.Log("Error during rollback: %v\n", rollbackError.Error())
			return rollbackError
		}

		return deleteErr
	}

	for index, value := range elements {

		logging.Info.Log("Inserting artifact with index %v in table %v\n", index, tableName)

		stmtIns, stmtInsErr := tx.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
		if stmtInsErr != nil {
			logging.Error.Log("Error while preparing insert statement for artifact in environment %v - %v\n", tableName, deleteErr.Error())
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				logging.Error.Log("Error during rollback: %v\n", rollbackError.Error())
				return rollbackError
			}
			fmt.Printf("Error: %v\n", stmtInsErr.Error())
			return stmtInsErr
		}

		defer stmtIns.Close()

		_, insErr := stmtIns.Exec(value)
		if insErr != nil {
			logging.Error.Log("Error while inserting artifact in environment %v - %v\n", tableName, deleteErr.Error())
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				logging.Error.Log("Error during rollback: %v\n", rollbackError.Error())
				return rollbackError
			}
			return insErr
		}

	}
	logging.Info.Log("Committing inserts\n")

	commitError := tx.Commit()
	if commitError != nil {
		logging.Error.Log("Error while committing update for environment %v - %v\n", tableName, deleteErr.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			logging.Error.Log("Error during rollback: %v\n", rollbackError.Error())
			return rollbackError
		}
		return commitError
	}

	return nil
}

func (client *MySqlClient) Insert(dbName string, tableName string, record string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Inserting record in table %v\n", tableName)

	stmtIns, err := client.Db.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
	if err != nil {
		logging.Error.Log("Error while preparing insert statment in table %v - %v\n", tableName, err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(record)
	if err != nil {
		logging.Error.Log("Error while inserting in organization %v - %v\n", tableName, err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) InsertOrReplace(dbName string, tableName string, name string, kind string, record string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Starting transaction\n")

	tx, startTransactionError := client.Db.Begin()
	if startTransactionError != nil {
		logging.Error.Log("Error starting transaction - %var name type\n", startTransactionError.Error())
		return startTransactionError
	}

	logging.Info.Log("Deleting record with name %v and kind %v from table %v\n", name, kind, tableName)

	stmtDelete, stmtError := tx.Prepare("DELETE FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if stmtError != nil {
		logging.Error.Log("Error while preparing delete statement for %v with name %v in environment %v in organization %v - %v\n", kind, name, tableName, name, startTransactionError.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		return stmtError
	}

	_, deleteError := stmtDelete.Exec()
	if deleteError != nil {
		logging.Error.Log("Error while deleting %v with name %v in environment %v in organization %v - %v\n", kind, name, tableName, dbName, startTransactionError.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		return deleteError
	}

	logging.Info.Log("Inserting %v with name %v in environment %v of organization %v\n", kind, name, tableName, dbName)

	stmtIns, err := tx.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
	if err != nil {
		logging.Error.Log("Error while preparing insert statement for %v with name %v in environment %v in organization %v - %v\n", kind, name, tableName, dbName, startTransactionError.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		return err
	}

	defer stmtIns.Close()

	_, insertError := stmtIns.Exec(record)
	if insertError != nil {
		logging.Error.Log("Error while inserting %v with name %v in environment %v in organization %v - %v\n", kind, name, tableName, name, startTransactionError.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		return insertError
	}

	logging.Info.Log("Committing insert\n")

	commitError := tx.Commit()
	if commitError != nil {
		logging.Error.Log("Error while committing insert of %v with name %v in environment %v in organization %v - %v\n", kind, name, tableName, name, startTransactionError.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			fmt.Printf("Error: %v\n", rollbackError.Error())
			return rollbackError
		}
		return commitError
	}

	return nil
}

func (client *MySqlClient) FindById(dbName string, tableName string, id int) (*Row, error) {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return nil, useDbErr
	}

	logging.Info.Log("Selecting record with id %v from table %v\n", id, tableName)

	stmtOut, err := client.Db.Prepare("SELECT * FROM `" + tableName + "` WHERE ID = ?")
	if err != nil {
		logging.Error.Log("Error preparing select statement for record with id %v in table %v - %v\n", id, tableName, err.Error())
		return nil, err
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	err = stmtOut.QueryRow(id).Scan(&resultId, &resultRecord)
	if err != nil {
		logging.Error.Log("Error selecting record with id %v from table %v - %v\n", id, tableName, err.Error())
		fmt.Printf("Error: %v\n", err.Error())
		return nil, err
	}

	return &Row{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) List(dbName string, tableName string, kind string) ([]Row, error) {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return []Row{}, useDbErr
	}

	logging.Info.Log("Selecting records with kind %v from table %v\n", kind, tableName)

	stmtOut, listErr := client.Db.Prepare("SELECT * FROM `" + tableName + "` WHERE Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if listErr != nil {
		logging.Error.Log("Error preparing select statement on table %v for records with kind %v - %v\n", tableName, kind, listErr.Error())
		return []Row{}, listErr
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	rows, err := stmtOut.Query()
	if err != nil {
		logging.Error.Log("Error selecting from table %v records with kind %v - %v\n", tableName, kind, err.Error())
		return []Row{}, err
	}
	defer rows.Close()

	var result []Row

	for rows.Next() {
		err := rows.Scan(&resultId, &resultRecord)
		if err != nil {
			logging.Error.Log("Error scanning select result - %v\n", err.Error())
			return []Row{}, err
		}

		result = append(result, Row{
			Id:     resultId,
			Record: resultRecord,
		})

	}

	return result, nil
}

func (client *MySqlClient) Update(dbName string, tableName string, id int, record string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Updating record with id %v in table %v\n", id, tableName)

	stmtIns, err := client.Db.Prepare("UPDATE `" + tableName + "` SET `Record` = ? WHERE ID = ?")
	if err != nil {
		logging.Error.Log("Error preparing update statement for record with id %v in table %v - %v\n", id, tableName, err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(record, id)
	if err != nil {
		logging.Error.Log("Error updating record with id %v in table %v - %v\n", id, tableName, err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) FindByNameAndKind(dbName string, tableName string, name string, kind string) (*Row, error) {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return nil, useDbErr
	}

	logging.Info.Log("Selecting records with kind %v and name %v from table %v\n", kind, name, tableName)

	stmtQuery, stmtError := client.Db.Prepare("SELECT * FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if stmtError != nil {
		logging.Error.Log("Error preparing select statement for record with name %v and kind %v in table %v - %v\n", name, kind, tableName, stmtError.Error())
		return nil, stmtError
	}

	defer stmtQuery.Close()

	var resultId int
	var resultRecord string

	queryError := stmtQuery.QueryRow().Scan(&resultId, &resultRecord)
	if queryError != nil {
		if queryError == sql.ErrNoRows {
			logging.Info.Log("No records found with name %v and kind %v in table %v\n", name, kind, tableName)
			return nil, nil
		} else {
			logging.Error.Log("Error selecting record with name %v and kind %v in table %v - %v\n", name, kind, tableName, queryError.Error())
			return nil, queryError
		}
	}

	return &Row{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) DeleteByNameAndKind(dbName string, tableName string, name string, kind string) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Deleting records with kind %v and name %v from table %v\n", kind, name, tableName)

	stmtIns, err := client.Db.Prepare("DELETE FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if err != nil {
		logging.Error.Log("Error preparing delete statement for records with name %v and kind %v on table %v - %v\n", name, kind, tableName, err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec()
	if err != nil {
		logging.Error.Log("Error deleting records with name %v and kind %v on table %v - %v\n", name, kind, tableName, err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) Delete(dbName string, tableName string, id int) error {

	logging.Info.Log("Using organization database %v\n", dbName)

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		logging.Error.Log("Error using database %v - %v\n", dbName, useDbErr.Error())
		return useDbErr
	}

	logging.Info.Log("Deleting record with id %v from table %v\n", id, tableName)

	stmtIns, err := client.Db.Prepare("DELETE FROM `" + tableName + "` WHERE ID = ?")
	if err != nil {
		logging.Error.Log("Error preparing delete statement for records with id %v on table %v - %v\n", id, tableName, err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(id)
	if err != nil {
		logging.Error.Log("Error deleting records with id %v on table %v - %v\n", id, tableName, err.Error())
		return err
	}

	return nil
}
