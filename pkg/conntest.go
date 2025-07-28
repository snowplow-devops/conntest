/*
 * Copyright (c) 2022 Snowplow Analytics Ltd. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
 */

package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v4"
	//nolint:gosec // This is a valid import for Databricks
	_ "github.com/databricks/databricks-sql-go"
	//nolint:gosec // This is a valid import for Postgres
	_ "github.com/lib/pq"
	"github.com/snowflakedb/gosnowflake"
	"github.com/xo/dburl"
	"gorm.io/driver/bigquery"

	//nolint:gosec // This is a valid import for BigQuery
	_ "gorm.io/driver/bigquery/driver"
	"gorm.io/gorm"
)

// DB parses a raw URI and returns a dburl.URL.
func DB(rawURI string) (*dburl.URL, error) {
	dsn, err := dburl.Parse(rawURI)

	if err == nil {
		return dsn, nil
	}
	return nil, errors.New("failed to parse DSN URI")
}

// Check tests a connection to a database and returns an Event.
func Check(uri dburl.URL, tags map[string]string, retryTimes uint) Event {
	var connErr, queryErr error
	gosnowflake.GetLogger().SetOutput(io.Discard)

	if strings.HasPrefix(uri.DSN, "bigquery") {
		// Do some Bigquery specific check
		fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connecting to Bigquery with "+uri.DSN)
		_, connErr := gorm.Open(bigquery.Open(uri.DSN), &gorm.Config{})
		if connErr != nil {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection error "+connErr.Error())
		} else {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection acquired")
		}
		// Set the query error same as the connection error to force an error response
		queryErr = connErr
	} else {
		// Do some non-Bigquery checks
		retry.Do(func() error {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection attempt ")
			db, connErrN := connect(uri)
			if connErrN != nil {
				fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection error "+connErrN.Error())
			} else {
				fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection acquired")
			}

			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query attempt")
			if db != nil {
				_, queryErrN := query(db, uri.Driver)
				if queryErrN != nil {
					fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query error "+queryErrN.Error())
				} else {
					fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query actioned")
				}
				queryErr = queryErrN
			} else {
				queryErr = connErrN
			}

			if connErr != nil && db != nil {
				db.Close()
			}

			connErr = connErrN
			return queryErr
		}, retry.Attempts(retryTimes), retry.OnRetry(func(u uint, err error) {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Retrying because of "+err.Error())
		}))
	}

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

func connect(uri dburl.URL) (*sql.DB, error) {
	fmt.Printf("Connecting to: %+v\n", uri)

	if strings.HasPrefix(uri.Scheme, "databricks") {
		return sql.Open("databricks", fmt.Sprintf("%s@%s%s?%s", uri.User, uri.Host, uri.Path, uri.RawQuery))
	}
	return dburl.Open(uri.String())
}

func query(db *sql.DB, driver string) (string, error) {
	var name string
	var queryString = queryFor(driver)
	err := db.QueryRow(queryString).Scan(&name)

	return name, err
}
