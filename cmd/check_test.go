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

package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func TestTagsVar(t *testing.T) {
	str := "lorem=ipsum;dolor=sit-amet"
	tags := tagsVar{}
	tags.Set(str)

	result := tags.String()

	// Parse the input string to get expected key-value pairs
	expectedPairs := make(map[string]string)
	for _, pair := range strings.Split(str, ";") {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			expectedPairs[parts[0]] = parts[1]
		}
	}

	// Parse the output string to get actual key-value pairs
	actualPairs := make(map[string]string)
	for _, pair := range strings.Split(result, ";") {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			actualPairs[parts[0]] = parts[1]
		}
	}

	// Compare the maps
	if !reflect.DeepEqual(expectedPairs, actualPairs) {
		t.Errorf("Expected pairs: %v, Got pairs: %v", expectedPairs, actualPairs)
		t.Fail()
	}
}
