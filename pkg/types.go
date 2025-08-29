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
	"time"

	uuid "github.com/google/uuid"
)

// Event represents a connection test result.
type Event struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	EmittedBy string      `json:"emittedBy"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewEvent creates a new Event.
func NewEvent(result Result) Event {
	name := "fabric:warehouse-connection-check"
	emittedBy := "conntest"
	version := 1

	return Event{uuid.New(), name, version, emittedBy, time.Now(), result}
}

// Result represents a connection test result.
type Result struct {
	Host     string            `json:"host"`
	Complete bool              `json:"complete"`
	Messages []string          `json:"messages"`
	Tags     map[string]string `json:"tags"`
	Attempts uint              `json:"attempts"`
}

// NewResult processes a connection test result.
func NewResult(host string, connError error, queryError error, tags map[string]string, attempts uint) Result {
	messages := []string{}

	if connError != nil {
		messages = append(messages, connError.Error())
	}

	if queryError != nil {
		messages = append(messages, queryError.Error())
	}

	return Result{host, connError == nil && queryError == nil, messages, tags, attempts}
}

// MultiResult represents multiple connection test results with summary.
type MultiResult struct {
	Results []Result `json:"results"`
	Summary Summary  `json:"summary"`
}

// Summary represents summary statistics for multiple connection tests.
type Summary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

// NewEventMultiple creates a new Event for multiple results.
func NewEventMultiple(results []Result) Event {
	succeeded := 0
	for _, r := range results {
		if r.Complete {
			succeeded++
		}
	}

	data := MultiResult{
		Results: results,
		Summary: Summary{
			Total:     len(results),
			Succeeded: succeeded,
			Failed:    len(results) - succeeded,
		},
	}

	name := "fabric:warehouse-connection-check"
	emittedBy := "conntest"
	version := 2

	return Event{uuid.New(), name, version, emittedBy, time.Now(), data}
}
