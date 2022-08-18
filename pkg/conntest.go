package pkg

import (
	"database/sql"
	"errors"
	"strings"

	retry "github.com/avast/retry-go/v4"
	_ "github.com/lib/pq"
	_ "github.com/snowflakedb/gosnowflake"
	"github.com/xo/dburl"
)

func ParseTags(raw string) map[string]string {
	splits := strings.Split(strings.Trim(raw, ";"), ";")
	tags := map[string]string{}

	for _, split := range splits {
		split := strings.Split(split, "=")
		tags[split[0]] = split[1]
	}

	return tags
}

func DB(rawUri string) (*dburl.URL, error) {
	dsn, err := dburl.Parse(rawUri)

	if err == nil {
		return dsn, nil
	} else {
		return nil, errors.New("Can't parse DSN URI")
	}
}

func Check(uri dburl.URL, tags map[string]string) Event {
	var connErr, queryErr error
	retry.Do(func() error {
		db, connErr := connect(uri.String())
		_, queryErr := query(db, uri.Driver)
		if connErr != nil {
			db.Close()
		}

		return queryErr
	})

	return NewEvent(NewResult(uri.Host, connErr, queryErr, tags))
}

func queryFor(driver string) string {
	dbs := map[string]string{
		"postgres":  `SELECT * FROM information_schema.information_schema_catalog_name;`,
		"snowflake": `SELECT * FROM information_schema.information_schema_catalog_name;`,
	}

	return dbs[driver]
}

func connect(uri string) (*sql.DB, error) {
	db, err := dburl.Open(uri)

	return db, err
}

func query(db *sql.DB, driver string) (string, error) {
	var name string
	var queryString = queryFor(driver)
	err := db.QueryRow(queryString).Scan(&name)

	return name, err
}
