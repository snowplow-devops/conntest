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

func TestQueryForDatabricks(t *testing.T) {
	res := queryFor("databricks")
	if !strings.Contains(res, "SELECT 1;") {
		t.Fail()
	}
}

func TestQueryForPostgres(t *testing.T) {
	res := queryFor("postgres")
	if !strings.Contains(res, "information_schema") {
		t.Fail()
	}
}

func TestQueryForSnowflake(t *testing.T) {
	res := queryFor("snowflake")
	if !strings.Contains(res, "information_schema") {
		t.Fail()
	}
}

func TestDBSnowflakeValid(t *testing.T) {
	_, err := DB("snowflake://lorem:ipsum@abcdefg-ab01234.snowflakecomputing.com/lorem?account=ab01234&ocspFailOpen=true&protocol=https&region=eu-central-1&role=SNOWPLOW_LOADER_ROLE&schema=SNOWPLOW&validateDefaultParameters=true&warehouse=COMPUTE_WH")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestDBSnowflakeInvalidEscapeChar(t *testing.T) {
	_, err := DB("snowflake://lorem:ip%sum@abcdefg-ab01234.snowflakecomputing.com/lorem?account=ab01234&ocspFailOpen=true&protocol=https&region=eu-central-1&role=SNOWPLOW_LOADER_ROLE&schema=SNOWPLOW&validateDefaultParameters=true&warehouse=COMPUTE_WH")

	// // Password needs sanitising here
	// pass, _ := dsn.User.Password()
	// t.Log("DSN=" + dsn.String())
	// t.Log("PASSWORD=" + pass)

	if err != nil {
		t.Log(err)
	}
}

func TestDBSnowflakeInvalidAmpersand(t *testing.T) {
	_, err := DB("snowflake://lorem:i&psum@abcdefg-ab01234.snowflakecomputing.com/lorem?account=ab01234&ocspFailOpen=true&protocol=https&region=eu-central-1&role=SNOWPLOW_LOADER_ROLE&schema=SNOWPLOW&validateDefaultParameters=true&warehouse=COMPUTE_WH")

	// // Password needs sanitising here
	// pass, _ := dsn.User.Password()
	// t.Log("DSN=" + dsn.String())
	// t.Log("PASSWORD=" + pass)

	// "Can't parse DSN URI"
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestMarshall(t *testing.T) {
	event := NewEvent(NewResult("lorem", nil, nil, map[string]string{"lorem": "ipsum"}, 1))
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
