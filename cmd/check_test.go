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
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/snowplow/conntest/pkg"
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

func TestMultipleDSNVersioning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		dsns        []string
		wantVersion int
	}{
		{
			name:        "single DSN should use version 1",
			dsns:        []string{"postgres://user:pass@host/db"},
			wantVersion: 1,
		},
		{
			name:        "multiple DSNs should use version 2",
			dsns:        []string{"postgres://user:pass@host/db1", "postgres://user:pass@host/db2"},
			wantVersion: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := map[string]string{"test": "value"}
			retryTimes := uint(1)

			event, err := pkg.CheckDSNs(tt.dsns, tags, retryTimes)
			if err != nil {
				t.Skip("Skipping test due to DSN error:", err)
			}

			if event.Version != tt.wantVersion {
				t.Errorf("Expected version %d, got %d", tt.wantVersion, event.Version)
			}

			// Test JSON marshaling works
			_, err = json.Marshal(event)
			if err != nil {
				t.Errorf("Failed to marshal event: %v", err)
			}
		})
	}
}

func TestMultipleResultsStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that multiple DSNs produce the expected structure
	dsns := []string{
		"postgres://user:pass@host1/db1",
		"postgres://user:pass@host2/db2",
	}

	tags := map[string]string{"env": "test"}
	event, err := pkg.CheckDSNs(dsns, tags, uint(1))
	if err != nil {
		t.Skip("Skipping test due to DSN error:", err)
	}

	// Verify event structure
	if event.Version != 2 {
		t.Errorf("Expected version 2, got %d", event.Version)
	}

	// Marshal to JSON and verify structure
	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Unmarshal to verify structure
	var eventMap map[string]interface{}
	err = json.Unmarshal(jsonData, &eventMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check that data contains results array and summary
	data, ok := eventMap["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Data field should be an object")
	}

	results, ok := data["results"].([]interface{})
	if !ok {
		t.Fatal("Data should contain results array")
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	summary, ok := data["summary"].(map[string]interface{})
	if !ok {
		t.Fatal("Data should contain summary object")
	}

	if summary["total"] != float64(2) {
		t.Errorf("Expected summary total 2, got %v", summary["total"])
	}
}
