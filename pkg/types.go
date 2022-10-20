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

type Event struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	EmittedBy string    `json:"emittedBy"`
	Timestamp time.Time `json:"timestamp"`
	Data      Result    `json:"data"`
}

func NewEvent(result Result) Event {
	name := "fabric:warehouse-connection-check"
	emittedBy := "conntest"
	version := 1

	return Event{uuid.New(), name, version, emittedBy, time.Now(), result}
}

type Result struct {
	Host     string            `json:"host"`
	Complete bool              `json:"complete"`
	Messages []string          `json:"messages"`
	Tags     map[string]string `json:"tags"`
	Attempts uint              `json:"attempts"`
}

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
