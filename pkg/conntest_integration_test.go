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
	"fmt"
	"reflect"
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

	dsn, err := DB(dsnStr)
	if err != nil {
		t.Fatal(err)
	}

	tags := map[string]string{}
	actual := Check(*dsn, tags, 1)
	expected := Result{dsn.Host, true, []string{}, tags, 1}

	if !reflect.DeepEqual(actual.Data, expected) {
		t.Fail()
	}
}

func SetupTestDatabase(ctx context.Context) (testcontainers.Container, string, error) {
	var user = "snowplow"
	var password = "snowplow"
	var db = "snowplow"

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
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
