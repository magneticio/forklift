package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	List(dbName string, tableName string) ([]Row, error)
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

		// fmt.Printf("Accessing DB: %v - %v - %v - %v \n", config.Sql.User, config.Sql.Password, host, config.Sql.Database)

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
			fmt.Printf("Error in schema drop: %v\n", dropError.Error())
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

func (client *MySqlClient) SetupEnvironment(dbName string, tableName string, elements []string) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	_, createErr := client.Db.Exec("CREATE TABLE IF NOT EXISTS `" + tableName + "` (ID int(11) NOT NULL AUTO_INCREMENT, Record mediumtext, PRIMARY KEY (ID))")
	if createErr != nil {
		fmt.Printf("Error during create: %v\n", createErr.Error())
		return createErr
	}

	for _, value := range elements {

		fmt.Println("Value:", value)

		insertErr := client.Insert(dbName, tableName, value)
		if insertErr != nil {
			fmt.Printf("Error during insert of %v - %v\n", value, insertErr.Error())
			_, dropError := client.Db.Exec("DROP TABLE `" + tableName + "`")
			if dropError != nil {
				fmt.Printf("Error in table drop: %v\n", dropError.Error())
				return dropError
			}
			return insertErr
		}

	}

	return nil
}

func (client *MySqlClient) UpdateEnvironment(dbName string, tableName string, elements []string) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	tx, startTransactionError := client.Db.Begin()
	if startTransactionError != nil {
		fmt.Printf("Error: %v\n", startTransactionError.Error())
		return startTransactionError
	}

	_, deleteErr := tx.Exec("DELETE FROM `" + tableName + "`")
	if deleteErr != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			fmt.Printf("Error: %v\n", rollbackError.Error())
			return rollbackError
		}
		fmt.Printf("Error: %v\n", deleteErr.Error())
		return deleteErr
	}

	for _, value := range elements {

		stmtIns, stmtInsErr := tx.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
		if stmtInsErr != nil {
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				fmt.Printf("Error: %v\n", rollbackError.Error())
				return rollbackError
			}
			fmt.Printf("Error: %v\n", stmtInsErr.Error())
			return stmtInsErr
		}

		defer stmtIns.Close()

		_, insErr := stmtIns.Exec(value)
		if insErr != nil {
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				fmt.Printf("Error: %v\n", rollbackError.Error())
				return rollbackError
			}
			fmt.Printf("Error: %v\n", insErr.Error())
			return insErr
		}

	}

	commitError := tx.Commit()
	if commitError != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			fmt.Printf("Error: %v\n", rollbackError.Error())
			return rollbackError
		}
		fmt.Printf("Error: %v\n", commitError.Error())
		return commitError
	}

	return nil
}

func (client *MySqlClient) Insert(dbName string, tableName string, record string) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(record)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func (client *MySqlClient) InsertOrReplace(dbName string, tableName string, name string, kind string, record string) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	tx, startTransactionError := client.Db.Begin()
	if startTransactionError != nil {
		fmt.Printf("Error: %v\n", startTransactionError.Error())
		return startTransactionError
	}

	stmtDelete, stmtError := tx.Prepare("DELETE FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if stmtError != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		fmt.Printf("Error: %v\n", stmtError.Error())
		return stmtError
	}

	_, deleteError := stmtDelete.Exec()
	if deleteError != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		fmt.Printf("Error: %v\n", deleteError.Error())
		return deleteError
	}

	stmtIns, err := tx.Prepare("INSERT INTO `" + tableName + "` ( Record ) VALUES( ? )")
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		fmt.Printf("Error: %v\n", err.Error())
		return err
	}

	defer stmtIns.Close()

	_, insertError := stmtIns.Exec(record)
	if insertError != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return rollbackError
		}
		fmt.Printf("Error: %v\n", err.Error())
		return insertError
	}

	commitError := tx.Commit()
	if commitError != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			fmt.Printf("Error: %v\n", rollbackError.Error())
			return rollbackError
		}
		fmt.Printf("Error: %v\n", commitError.Error())
		return commitError
	}

	return nil
}

func (client *MySqlClient) FindById(dbName string, tableName string, id int) (*Row, error) {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return nil, useDbErr
	}

	stmtOut, err := client.Db.Prepare("SELECT * FROM `" + tableName + "` WHERE ID = ?")
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

	return &Row{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) List(dbName string, tableName string) ([]Row, error) {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return []Row{}, useDbErr
	}

	stmtOut, listErr := client.Db.Prepare("SELECT * FROM `" + tableName + "`")
	if listErr != nil {
		fmt.Printf("Error: %v\n", listErr.Error())
		return []Row{}, listErr
	}

	defer stmtOut.Close()

	var resultId int
	var resultRecord string

	rows, err := stmtOut.Query()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return []Row{}, err
	}
	defer rows.Close()

	var result []Row

	for rows.Next() {
		err := rows.Scan(&resultId, &resultRecord)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
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

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("UPDATE `" + tableName + "` SET `Record` = ? WHERE ID = ?")
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

func (client *MySqlClient) FindByNameAndKind(dbName string, tableName string, name string, kind string) (*Row, error) {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return nil, useDbErr
	}

	stmtQuery, stmtError := client.Db.Prepare("SELECT * FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
	if stmtError != nil {
		fmt.Printf("Error: %v\n", stmtError.Error())
		return nil, stmtError
	}

	defer stmtQuery.Close()

	var resultId int
	var resultRecord string

	queryError := stmtQuery.QueryRow().Scan(&resultId, &resultRecord)
	if queryError != nil {
		if queryError == sql.ErrNoRows {
			// fmt.Printf("No rows\n")
			return nil, nil
		} else {
			fmt.Printf("Query error: %v\n", queryError.Error())
			return nil, queryError
		}
	}

	return &Row{
		Id:     resultId,
		Record: resultRecord,
	}, nil
}

func (client *MySqlClient) DeleteByNameAndKind(dbName string, tableName string, name string, kind string) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("DELETE FROM `" + tableName + "` WHERE Record LIKE '%\"name\":\"" + name + "\"%' AND Record LIKE '%\"kind\":\"" + kind + "\"%'")
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

func (client *MySqlClient) Delete(dbName string, tableName string, id int) error {

	_, useDbErr := client.Db.Exec("USE `" + dbName + "`")
	if useDbErr != nil {
		fmt.Printf("Error: %v\n", useDbErr.Error())
		return useDbErr
	}

	stmtIns, err := client.Db.Prepare("DELETE FROM `" + tableName + "` WHERE ID = ?")
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
