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

import "github.com/xo/dburl"

// Functions added below to support Databricks which is not currently supported in xo/dburl
// and xo/usql - the Databricks driver is still in beta
func RegisterDatabricks() {
	dburl.Register(dburl.Scheme{"databricks", GenDatabricks, 0, false, []string{"databricks"}, ""})
}

// GenDatabricks generates a databricks DSN from the passed URL.
// Format is here https://github.com/databricks/databricks-sql-go#usage
//
// databricks://:[your token]@[Workspace hostname][Endpoint HTTP Path]
func GenDatabricks(u *dburl.URL) (string, error) {
	host := u.Hostname()
	if host == "" {
		return "", dburl.ErrMissingHost
	}

	// add auth token, which is passed as the password
	user := ""
	if pass, _ := u.User.Password(); pass != "" {
		user += ":" + pass
	}

	return "databricks://" + user + "@" + host + u.Path, nil
}

func GetProtocols(name string) []string {
	return dburl.Protocols(name)
}

func SchemeDriverAndAliases(name string) (string, []string) {
	return dburl.SchemeDriverAndAliases(name)
}
