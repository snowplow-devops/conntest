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

func TestDB_DatabricksOAuth(t *testing.T) {
	tests := []struct {
		name    string
		rawURI  string
		wantErr bool
		check   func(*testing.T, string, error)
	}{
		{
			name:   "valid OAuth M2M DSN",
			rawURI: "databricks://client123:secret456@dbc-abc123.cloud.databricks.com/sql/1.0/endpoints/xyz789",
			check: func(t *testing.T, uri string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				dsn, parseErr := DB(uri)
				if parseErr != nil {
					t.Errorf("Expected successful parsing, got %v", parseErr)
				}
				if dsn.Scheme != "databricks" {
					t.Errorf("Expected scheme 'databricks', got '%s'", dsn.Scheme)
				}
				if dsn.User.Username() != "client123" {
					t.Errorf("Expected client ID 'client123', got '%s'", dsn.User.Username())
				}
				if secret, _ := dsn.User.Password(); secret != "secret456" {
					t.Errorf("Expected client secret 'secret456', got '%s'", secret)
				}
			},
		},
		{
			name:   "regular databricks DSN should still work",
			rawURI: "databricks://token:abc123@dbc-abc123.cloud.databricks.com/sql/1.0/endpoints/xyz789",
			check: func(t *testing.T, uri string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				dsn, parseErr := DB(uri)
				if parseErr != nil {
					t.Errorf("Expected successful parsing, got %v", parseErr)
				}
				if dsn.Scheme != "databricks" {
					t.Errorf("Expected scheme 'databricks', got '%s'", dsn.Scheme)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.check != nil {
				tt.check(t, tt.rawURI, nil)
			}
		})
	}
}
