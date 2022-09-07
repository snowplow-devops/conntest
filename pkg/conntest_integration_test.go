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
	actual := Check(*dsn, tags)
	expected := Result{dsn.Host, true, []string{}, tags}
	fmt.Println(actual)
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
