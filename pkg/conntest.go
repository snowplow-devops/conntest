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
	"time"

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
		fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection attempt ")
		db, connErrN := connect(uri.String())
		if connErrN != nil {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection error "+connErrN.Error())
		} else {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Connection acquired")
		}

		fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query attempt")
		_, queryErrN := query(db, uri.Driver)
		if queryErrN != nil {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query error "+queryErrN.Error())
		} else {
			fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Query actioned")
		}

		if connErr != nil {
			db.Close()
		}

		connErr = connErrN
		queryErr = queryErrN
		return queryErr
	}, retry.Attempts(retryTimes), retry.OnRetry(func(u uint, err error) {
		fmt.Fprintln(os.Stderr, time.Now().Format("03:04:05.000000")+" Retrying because of "+err.Error())
	}))

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
