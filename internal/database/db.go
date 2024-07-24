package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"

	"aniverse/internal/domain/types"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var dbType string
var postgres *sql.DB
var sqlite *sql.DB

func InitPostgres(connectionString string) error {
	var err error
	postgres, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	dbType = "postgresql"
	return postgres.Ping()
}

func InitSQLite(databaseFile string) error {
	var err error
	sqlite, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		return err
	}
	dbType = "sqlite"
	return sqlite.Ping()
}

func GetMedia(id string, fields []string) (interface{}, error) {
	switch dbType {
	case "postgresql":
		return getFromPostgres(id, fields)
	case "sqlite":
		return getFromSQLite(id, fields)
	default:
		return nil, errors.New("unsupported database type")
	}
}

func getFromPostgres(id string, fields []string) (interface{}, error) {
	var data []byte
	query := "SELECT * FROM anime WHERE id = $1"
	row := postgres.QueryRow(query, id)
	if err := row.Scan(&data); err != nil {
		query := "SELECT * FROM manga WHERE id = $1"
		row := postgres.QueryRow(query, id)
		if err := row.Scan(&data); err != nil {
			return nil, err
		}
		var manga types.Manga
		if err := json.Unmarshal(data, &manga); err != nil {
			return nil, err
		}
		filterFields(&manga, fields)
		return manga, nil
	}
	var anime types.Anime
	if err := json.Unmarshal(data, &anime); err != nil {
		return nil, err
	}
	filterFields(&anime, fields)
	return anime, nil
}

func getFromSQLite(id string, fields []string) (interface{}, error) {
	query := "SELECT * FROM anime WHERE id = ?"
	row := sqlite.QueryRow(query, id)
	var data []byte
	if err := row.Scan(&data); err != nil {
		query := "SELECT * FROM manga WHERE id = ?"
		row := sqlite.QueryRow(query, id)
		if err := row.Scan(&data); err != nil {
			return nil, err
		}
		var manga types.Manga
		if err := parseSQLiteData(data, &manga); err != nil {
			return nil, err
		}
		filterFields(&manga, fields)
		return manga, nil
	}
	var anime types.Anime
	if err := parseSQLiteData(data, &anime); err != nil {
		return nil, err
	}
	filterFields(&anime, fields)
	return anime, nil
}

func parseSQLiteData(data []byte, v interface{}) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		switch value.(type) {
		case string:
			// Assuming JSON fields are stored as strings
			var parsed interface{}
			if err := json.Unmarshal([]byte(value.(string)), &parsed); err == nil {
				raw[key] = parsed
			}
		}
	}
	parsedData, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(parsedData, v)
}

func filterFields(v interface{}, fields []string) {
	if len(fields) == 0 {
		return
	}
	value := reflect.ValueOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if !contains(fields, field.Name) {
			fieldValue := value.FieldByName(field.Name)
			if fieldValue.IsValid() {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
