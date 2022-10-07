package pkg

import (
	"testing"
)

func TestRegisterScheme(t *testing.T) {

	t.Log("Register Databricks scheme")
	RegisterDatabricks()

	t.Log("Get protocols for database scheme")
	protocols := GetProtocols("databricks")
	t.Log(protocols)

	if isElementExist(protocols, "snowplow") != false {
		t.Fail()
	}
	if isElementExist(protocols, "databricks") != true {
		t.Fail()
	}

	t.Log("Get scheme driver and aliases for scheme")
	t.Log(SchemeDriverAndAliases("databricks"))
}

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
