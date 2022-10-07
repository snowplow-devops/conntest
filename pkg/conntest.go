package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	retry "github.com/avast/retry-go/v4"
	_ "github.com/databricks/databricks-sql-go"
	_ "github.com/lib/pq"
	"github.com/snowflakedb/gosnowflake"
	"github.com/xo/dburl"
)

func DB(rawUri string) (*dburl.URL, error) {
	dsn, err := dburl.Parse(rawUri)

	if err == nil {
		return dsn, nil
	} else {
		return nil, errors.New("Can't parse DSN URI")
	}
}

func Check(uri dburl.URL, tags map[string]string, retryTimes uint) Event {
	var connErr, queryErr error
	gosnowflake.GetLogger().SetOutput(io.Discard)

	retry.Do(func() error {
		db, connErrN := connect(uri.String())
		_, queryErrN := query(db, uri.Driver)
		if connErr != nil {
			db.Close()
		}
		connErr = connErrN
		queryErr = queryErrN
		return queryErr
	}, retry.Attempts(retryTimes), retry.OnRetry(func(u uint, err error) { fmt.Fprintln(os.Stderr, "Retrying because of", err.Error()) }))

	return NewEvent(NewResult(uri.Host, connErr, queryErr, tags, retryTimes))
}

func queryFor(driver string) string {
	dbs := map[string]string{
		"databricks": `SELECT 1;`,
		"postgres":   `SELECT * FROM information_schema.information_schema_catalog_name;`,
		"snowflake":  `SELECT * FROM information_schema.information_schema_catalog_name;`,
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
