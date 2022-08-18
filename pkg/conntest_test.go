package pkg

import (
	"errors"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	if status(nil) != "ok" {
		t.Fail()
	}

	if status(errors.New("lorem")) == "ok" {
		t.Fail()
	}
}

func TestQueryFor(t *testing.T) {
	res := queryFor("postgres")
	if !strings.Contains(res, "information_schema") {
		t.Fail()
	}
}

func TestMarshall(t *testing.T) {
	event := NewEvent(NewResult("lorem", nil, nil, map[string]string{"lorem": "ipsum"}))
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

func TestParseTags(t *testing.T) {
	expected := map[string]string{"dolor": "sit-amet", "lorem": "ipsum"}
	sample := "lorem=ipsum;dolor=sit-amet"
	actuals := []map[string]string{
		ParseTags(sample),
		ParseTags(";" + sample + ";"),
	}

	for i, actual := range actuals {
		if !reflect.DeepEqual(actual, expected) {
			t.Log("error at", i, ":", actual, expected)
			t.Fail()
		}
	}
}
