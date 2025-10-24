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
	"net"
	"strconv"
	"strings"

	retry "github.com/avast/retry-go/v4"
	dbsql "github.com/databricks/databricks-sql-go"
	m2m "github.com/databricks/databricks-sql-go/auth/oauth/m2m"

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

// CheckDSNs parses and tests connections for all DSN strings.
// Returns an Event with results (version 1 for single DSN, version 2 for multiple).
func CheckDSNs(dsnStrings []string, tags map[string]string, retryTimes uint) (Event, error) {
	results := make([]Result, 0, len(dsnStrings))

	for _, dsnString := range dsnStrings {
		var result Result

		if strings.HasPrefix(dsnString, "git+ssh://") {
			gitURL, keyFile, host, err := parseGitDSN(dsnString)
			if err != nil {
				return Event{}, err
			}
			result = checkGit(gitURL, keyFile, host, tags, retryTimes)
		} else {
			dsn, err := DB(dsnString)
			if err != nil {
				return Event{}, err
			}
			result = checkSingle(*dsn, tags, retryTimes)
		}

		results = append(results, result)
	}

	// Return version 1 for single DSN, version 2 for multiple
	if len(results) == 1 {
		return NewEvent(results[0]), nil
	}
	return NewEventMultiple(results), nil
}

// checkSingle tests a connection to a single database and returns a Result.
func checkSingle(uri dburl.URL, tags map[string]string, retryTimes uint) Result {
	var connErr, queryErr error
	gosnowflake.GetLogger().SetOutput(io.Discard)

	if strings.HasPrefix(uri.DSN, "bigquery") {
		// Do some Bigquery specific check
		Logger.Info(fmt.Sprintf("Connecting to Bigquery with %s", uri.DSN))
		_, connErr := gorm.Open(bigquery.Open(uri.DSN), &gorm.Config{})
		if connErr != nil {
			Logger.Info(fmt.Sprintf("Connection error %s", connErr.Error()))
		} else {
			Logger.Info("Connection acquired")
		}
		// Set the query error same as the connection error to force an error response
		queryErr = connErr
	} else {
		// Do some non-Bigquery checks
		retry.Do(func() error {
			Logger.Info("Connection attempt")
			db, connErrN := connect(uri)
			if connErrN != nil {
				Logger.Info(fmt.Sprintf("Connection error %s", connErrN.Error()))
			} else {
				Logger.Info("Connection acquired")
			}

			Logger.Info("Query attempt")
			if db != nil {
				_, queryErrN := query(db, uri.Driver)
				if queryErrN != nil {
					Logger.Info(fmt.Sprintf("Query error %s", queryErrN.Error()))
				} else {
					Logger.Info("Query actioned")
				}
				queryErr = queryErrN
			} else {
				queryErr = connErrN
			}

			if connErr == nil && db != nil {
				db.Close()
			}

			connErr = connErrN
			return queryErr
		}, retry.Attempts(retryTimes), retry.OnRetry(func(u uint, err error) {
			Logger.Info(fmt.Sprintf("Retrying because of %s", err.Error()))
		}))
	}

	return NewResult(uri.Host, connErr, queryErr, tags, retryTimes)
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
	if uri.Scheme == "databricks" {
		host, portString, err := net.SplitHostPort(uri.Host)
		if err != nil {
			host = uri.Host
			portString = "443"
		}
		port, err := strconv.Atoi(portString)
		if err != nil {
			return nil, fmt.Errorf("failed to convert port to integer: %w", err)
		}

		// If the connection is using PAT authentication
		if uri.User.Username() == "token" {
			token, hasToken := uri.User.Password()
			if !hasToken || token == "" {
				return nil, errors.New("databricks PAT authentication requires a token in DSN")
			}
			return sql.Open("databricks", fmt.Sprintf("token:%s@%s:%d%s", token, host, port, uri.Path))
		}

		// If the connection is using OAuth M2M authentication
		clientSecret, hasClientSecret := uri.User.Password()
		if !hasClientSecret || clientSecret == "" {
			return nil, errors.New("databricks OAuth M2M authentication requires a client secret in DSN")
		}
		authenticator := m2m.NewAuthenticator(
			uri.User.Username(),
			clientSecret,
			host,
		)
		connector, err := dbsql.NewConnector(
			dbsql.WithServerHostname(host),
			dbsql.WithHTTPPath(uri.Path),
			dbsql.WithPort(port),
			dbsql.WithAuthenticator(authenticator),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OAuth M2M connector: %w", err)
		}
		return sql.OpenDB(connector), nil
	}

	return dburl.Open(uri.String())
}

func query(db *sql.DB, driver string) (string, error) {
	var name string
	var queryString = queryFor(driver)
	err := db.QueryRow(queryString).Scan(&name)

	return name, err
}

// checkGit tests a connection to a Git repository and returns a Result.
func checkGit(gitURL, keyPath, host string, tags map[string]string, retryTimes uint) Result {
	var gitErr error

	retry.Do(func() error {
		Logger.Info(fmt.Sprintf("Git connection attempt to %s", gitURL))
		err := connectToGitRepo(gitURL, keyPath)
		if err != nil {
			Logger.Info(fmt.Sprintf("Git connection error %s", err.Error()))
		} else {
			Logger.Info("Git clone successful")
		}
		gitErr = err
		return gitErr
	}, retry.Attempts(retryTimes), retry.OnRetry(func(u uint, err error) {
		Logger.Info(fmt.Sprintf("Retrying because of %s", err.Error()))
	}))

	return NewResult(host, gitErr, nil, tags, retryTimes)
}
