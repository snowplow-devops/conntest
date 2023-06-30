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
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestQueryFor(t *testing.T) {
	res := queryFor("postgres")
	if !strings.Contains(res, "information_schema") {
		t.Fail()
	}
}

func TestMarshall(t *testing.T) {
	event := NewEvent(NewResult(nil, map[string]string{"lorem": "ipsum"}, 1))
	var unmarshaled Event
	marshaled, err := json.Marshal(event)
	res := json.Unmarshal(marshaled, &unmarshaled)

	if err != nil {
		t.Fail()
	}

	if reflect.DeepEqual(event, unmarshaled) {
		t.Log(unmarshaled, res)
		t.Fail()
	}
}

func TestCheckPostgres(t *testing.T) {
	result := Check("postgres", "pq://user:pass@localhost:5432/db", map[string]string{}, 1)

	if !containsError(result.Data.Messages, "can't parse Postgres DSN: cannot parse `pq://user:xxxxxx@localhost:5432/db`: failed to parse as DSN (invalid dsn)") {
		t.Fail()
	}
}

func TestCheckSnowflake(t *testing.T) {
	result := Check("snowflake", "incorrect@urb12345.us-east-1.snowflakecomputing.com/db", map[string]string{}, 1)

	if !containsError(result.Data.Messages, "can't parse Snowflake DSN: 260002: password is empty") {
		t.Fail()
	}
}

func TestCheckDatabricks(t *testing.T) {
	result := Check("databricks", "db://token:dapi12345@dbc-12345.cloud.databricks.com:443/sql/path", map[string]string{}, 1)

	if !containsError(result.Data.Messages, "can't open a connection to databricks: scheme db not recognized") {
		t.Fail()
	}
}

func TestCheckUnknownDriver(t *testing.T) {
	result := Check("unknown", "db://user:pass@localhost:123/db", map[string]string{}, 1)

	if !containsError(result.Data.Messages, "unknown driver") {
		t.Fail()
	}
}

func containsError(s []string, errorMessage string) bool {
	for _, v := range s {
		if strings.Contains(v, errorMessage) {
			return true
		}
	}
	return false
}
