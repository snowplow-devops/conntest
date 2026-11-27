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
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: not running")
	}

	ctx := context.Background()
	dbC, dsnStr, err := SetupTestDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer dbC.Terminate(ctx)

	tags := map[string]string{}
	event, err := CheckDSNs([]string{dsnStr}, tags, 1)
	if err != nil {
		t.Fatal(err)
	}

	dsn, err := DB(dsnStr)
	if err != nil {
		t.Fatal(err)
	}

	expected := Result{dsn.Host, true, []string{}, tags, 1}

	if !reflect.DeepEqual(event.Data, expected) {
		t.Fail()
	}
}

func TestCheckPostgresSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: not running")
	}

	ctx := context.Background()
	dbC, dsnStr, err := SetupTestDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer dbC.Terminate(ctx)

	admin, err := sql.Open("postgres", dsnStr)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	for _, stmt := range []string{
		`CREATE SCHEMA existing`,
		`CREATE ROLE restricted LOGIN PASSWORD 'restricted'`,
	} {
		if _, err := admin.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	restrictedDSN := strings.Replace(dsnStr, "snowplow:snowplow@", "restricted:restricted@", 1)

	tests := []struct {
		name         string
		dsn          string
		schema       string
		wantComplete bool
		wantMessage  string
	}{
		{"existing schema, owner", dsnStr, "existing", true, ""},
		{"missing schema, owner of database", dsnStr, "atomic", false, "does not exist"},
		{"existing schema, no privileges", restrictedDSN, "existing", false, "lacks USAGE and/or CREATE"},
		{"missing schema, restricted", restrictedDSN, "atomic", false, "does not exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := CheckDSNs([]string{tt.dsn + "&schema=" + tt.schema}, map[string]string{}, 1)
			if err != nil {
				t.Fatal(err)
			}
			res := event.Data.(Result)
			if res.Complete != tt.wantComplete {
				t.Fatalf("complete = %v, want %v (messages: %v)", res.Complete, tt.wantComplete, res.Messages)
			}
			if tt.wantMessage != "" && (len(res.Messages) == 0 || !strings.Contains(res.Messages[0], tt.wantMessage)) {
				t.Errorf("messages = %v, want one containing %q", res.Messages, tt.wantMessage)
			}
		})
	}
}

func SetupTestDatabase(ctx context.Context) (testcontainers.Container, string, error) {
	var user = "snowplow"
	var password = "snowplow"
	var db = "snowplow"

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		// Postgres logs "ready" once for the temporary init server and again for the real one.
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       db,
		},
	}

	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		},
	)
	if err != nil {
		return nil, "", err
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port.Port(), db)

	return dbContainer, dsn, err
}
