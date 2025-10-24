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
	"io"
	"log/slog"
	"os"
)

// simpleHandler is a custom slog handler that formats logs to match current output format.
type simpleHandler struct {
	w io.Writer
}

func (h *simpleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *simpleHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := r.Time.Format("03:04:05.000000")
	_, err := fmt.Fprintf(h.w, "%s %s\n", timestamp, r.Message)
	return err
}

func (h *simpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *simpleHandler) WithGroup(name string) slog.Handler {
	return h
}

// Logger is the package-level logger for conntest.
var Logger *slog.Logger

func init() {
	handler := &simpleHandler{w: os.Stderr}
	Logger = slog.New(handler)
}
