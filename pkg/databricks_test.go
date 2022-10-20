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
