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
	"fmt"
	"os"
	"time"

	retry "github.com/avast/retry-go/v4"
	_ "github.com/databricks/databricks-sql-go"
	pgx "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	sf "github.com/snowflakedb/gosnowflake"
)

func Check(driver string, dsn string, tags map[string]string, retryTimes uint) Event {
	err := retry.Do(func() error {
		switch driver {
		case "snowflake":
			return checkSnowflake(dsn)
		case "databricks":
			return checkConnection("databricks", dsn)
		case "postgres":
			return checkPostgres(dsn)
		default:
			return fmt.Errorf("unknown driver: %v", driver)
		}
	}, retry.Attempts(retryTimes), retry.LastErrorOnly(true), retry.OnRetry(func(u uint, err error) {
		fmt.Fprintf(os.Stderr, "%s %d attempt: %s\n", time.Now().Format("03:04:05.000000"), u, err.Error())
	}))

	if err != nil {
		return NewEvent(NewResult(err, tags, retryTimes))
	}

	return NewEvent(NewResult(nil, tags, retryTimes))
}

func checkConnection(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return fmt.Errorf("can't open a connection to %s: %w", driver, err)
	}

	defer db.Close()

	if _, err := query(db, driver); err != nil {
		return fmt.Errorf("can't query %s: %w", driver, err)
	}

	return nil

}

func checkSnowflake(dsn string) error {
	var err error
	cfg, err := sf.ParseDSN(dsn)

	if err != nil {
		return fmt.Errorf("can't parse Snowflake DSN: %w", err)
	}

	validatedDsn, err := sf.DSN(cfg)
	if err != nil {
		return fmt.Errorf("can't construct Snowflake DSN: %w", err)
	}

	return checkConnection("snowflake", validatedDsn)
}

func checkPostgres(dsn string) error {
	cnf, err := pgx.ParseConfig(dsn)

	if err != nil {
		return fmt.Errorf("can't parse Postgres DSN: %w", err)
	}

	return checkConnection("postgres", cnf.ConnString())
}

func queryFor(driver string) string {
	dbs := map[string]string{
		"databricks": `SELECT 1;`,
		"postgres":   `SELECT * FROM information_schema.information_schema_catalog_name;`,
		"snowflake":  `SELECT * FROM information_schema.information_schema_catalog_name;`,
	}

	return dbs[driver]
}

func query(db *sql.DB, driver string) (string, error) {
	var name string
	var queryString = queryFor(driver)
	err := db.QueryRow(queryString).Scan(&name)
	return name, err
}
