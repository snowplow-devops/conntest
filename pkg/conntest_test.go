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

func TestExtractPostgresSchema(t *testing.T) {
	tests := []struct {
		name       string
		rawURI     string
		wantSchema string
		wantInDSN  []string
		notInDSN   []string
	}{
		{
			name:       "schema is extracted and stripped, TLS params kept",
			rawURI:     "postgres://user:pass@db.example.com:5432/snowplow?schema=atomic&sslmode=verify-full&sslrootcert=%2Ftmp%2Fca.pem",
			wantSchema: "atomic",
			wantInDSN:  []string{"sslmode=verify-full", "sslrootcert=", "db.example.com"},
			notInDSN:   []string{"schema"},
		},
		{
			name:       "no schema param leaves DSN untouched",
			rawURI:     "postgres://user:pass@db.example.com:5432/snowplow?sslmode=disable",
			wantSchema: "",
			wantInDSN:  []string{"sslmode=disable"},
		},
		{
			name:       "non-postgres drivers are ignored",
			rawURI:     "snowflake://user:pass@acct.snowflakecomputing.com/db?schema=atomic",
			wantSchema: "",
			wantInDSN:  []string{"schema=atomic"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := DB(tt.rawURI)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			cleaned, schema, err := extractPostgresSchema(dsn)
			if err != nil {
				t.Fatalf("extract: %v", err)
			}
			if schema != tt.wantSchema {
				t.Errorf("schema = %q, want %q", schema, tt.wantSchema)
			}
			for _, s := range tt.wantInDSN {
				if !strings.Contains(cleaned.String(), s) {
					t.Errorf("expected %q in DSN %q", s, cleaned.String())
				}
			}
			for _, s := range tt.notInDSN {
				if strings.Contains(cleaned.String(), s) || strings.Contains(cleaned.DSN, s) {
					t.Errorf("did not expect %q in DSN %q / %q", s, cleaned.String(), cleaned.DSN)
				}
			}
		})
	}
}
