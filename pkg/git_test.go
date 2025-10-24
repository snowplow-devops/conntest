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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseGitDSN(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		wantGitURL  string
		wantKeyFile string
		wantHost    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "github standard",
			dsn:         "git+ssh://git@github.com/user/repo.git?keyfile=/path/to/key",
			wantGitURL:  "git@github.com:user/repo.git",
			wantKeyFile: "/path/to/key",
			wantHost:    "github.com",
			wantErr:     false,
		},
		{
			name:        "gitlab with group",
			dsn:         "git+ssh://git@gitlab.com/group/project.git?keyfile=/home/user/.ssh/id_rsa",
			wantGitURL:  "git@gitlab.com:group/project.git",
			wantKeyFile: "/home/user/.ssh/id_rsa",
			wantHost:    "gitlab.com",
			wantErr:     false,
		},
		{
			name:        "bitbucket",
			dsn:         "git+ssh://git@bitbucket.org/company/repo.git?keyfile=/keys/deploy",
			wantGitURL:  "git@bitbucket.org:company/repo.git",
			wantKeyFile: "/keys/deploy",
			wantHost:    "bitbucket.org",
			wantErr:     false,
		},
		{
			name:        "custom host",
			dsn:         "git+ssh://git@git.example.com/org/repo.git?keyfile=/etc/ssh/key",
			wantGitURL:  "git@git.example.com:org/repo.git",
			wantKeyFile: "/etc/ssh/key",
			wantHost:    "git.example.com",
			wantErr:     false,
		},
		{
			name:        "nested path",
			dsn:         "git+ssh://git@github.com/org/team/repo.git?keyfile=/path/key",
			wantGitURL:  "git@github.com:org/team/repo.git",
			wantKeyFile: "/path/key",
			wantHost:    "github.com",
			wantErr:     false,
		},
		{
			name:        "without .git extension",
			dsn:         "git+ssh://git@github.com/user/repo?keyfile=/path/key",
			wantGitURL:  "git@github.com:user/repo",
			wantKeyFile: "/path/key",
			wantHost:    "github.com",
			wantErr:     false,
		},
		{
			name:        "with multiple query params",
			dsn:         "git+ssh://git@github.com/user/repo.git?keyfile=/path/key&other=value",
			wantGitURL:  "git@github.com:user/repo.git",
			wantKeyFile: "/path/key",
			wantHost:    "github.com",
			wantErr:     false,
		},
		{
			name:        "missing keyfile parameter",
			dsn:         "git+ssh://git@github.com/user/repo.git",
			wantErr:     true,
			errContains: "keyfile query parameter",
		},
		{
			name:        "empty keyfile value",
			dsn:         "git+ssh://git@github.com/user/repo.git?keyfile=",
			wantErr:     true,
			errContains: "keyfile query parameter",
		},
		{
			name:        "invalid URL format",
			dsn:         "not-a-valid-url",
			wantErr:     true,
		},
		{
			name:        "malformed git URL",
			dsn:         "git+ssh://github.com/user/repo.git?keyfile=/path/key",
			wantGitURL:  "git@github.com:user/repo.git",
			wantKeyFile: "/path/key",
			wantHost:    "github.com",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGitURL, gotKeyFile, gotHost, err := parseGitDSN(tt.dsn)

			// Check error cases
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if gotGitURL != tt.wantGitURL {
				t.Errorf("GitURL = %q, want %q", gotGitURL, tt.wantGitURL)
			}

			if gotKeyFile != tt.wantKeyFile {
				t.Errorf("KeyFile = %q, want %q", gotKeyFile, tt.wantKeyFile)
			}

			if gotHost != tt.wantHost {
				t.Errorf("Host = %q, want %q", gotHost, tt.wantHost)
			}

			// Verify format: should have colon, not slash after host
			if strings.Contains(gotGitURL, "git@") {
				if !strings.Contains(gotGitURL, ":") {
					t.Errorf("GitURL missing colon separator: %q (should be git@host:path format)", gotGitURL)
				}
				// Ensure it doesn't have slash immediately after host
				atIndex := strings.Index(gotGitURL, "@")
				if atIndex >= 0 && atIndex+1 < len(gotGitURL) {
					afterAt := gotGitURL[atIndex+1:]
					colonIndex := strings.Index(afterAt, ":")
					slashIndex := strings.Index(afterAt, "/")
					if slashIndex >= 0 && (colonIndex < 0 || slashIndex < colonIndex) {
						t.Errorf("GitURL has slash before colon: %q (wrong format, should use colon)", gotGitURL)
					}
				}
			}
		})
	}
}

func TestConnectToGitRepo_InvalidKeyPath(t *testing.T) {
	err := connectToGitRepo("git@github.com:snowplow/conntest.git", "/nonexistent/key/path")
	if err == nil {
		t.Error("Expected error for nonexistent key file, got nil")
	}
	if err != nil && !strings.HasPrefix(err.Error(), "SSH key file not accessible") {
		t.Errorf("Expected 'SSH key file not accessible' error, got: %v", err)
	}
}

func TestConnectToGitRepo_InvalidKeyFormat(t *testing.T) {
	// Create a temporary file with invalid SSH key content
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "invalid_key")
	err := os.WriteFile(keyPath, []byte("not a valid ssh key"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}

	err = connectToGitRepo("git@github.com:snowplow/conntest.git", keyPath)
	if err == nil {
		t.Error("Expected error for invalid key format, got nil")
	}
}

func TestCheckGit(t *testing.T) {
	tags := map[string]string{"test": "true"}
	result := checkGit("git@github.com:snowplow/conntest.git", "/nonexistent/key", "github.com", tags, 1)

	if result.Complete {
		t.Error("Expected Complete=false for nonexistent key, got true")
	}
	if len(result.Messages) == 0 {
		t.Error("Expected error messages, got none")
	}
	if result.Host != "github.com" {
		t.Errorf("Expected host=github.com, got: %s", result.Host)
	}
}

func TestCheckGitSingle(t *testing.T) {
	tags := map[string]string{"env": "test"}
	result := checkGit("git@gitlab.com:example/repo.git", "/fake/key", "gitlab.com", tags, 1)

	if result.Complete {
		t.Error("Expected Complete=false for invalid connection, got true")
	}
	if result.Host != "gitlab.com" {
		t.Errorf("Expected host=gitlab.com, got: %s", result.Host)
	}
	if result.Tags["env"] != "test" {
		t.Error("Expected tags to be preserved in result")
	}
}

func TestCheckGit_HostExtraction(t *testing.T) {
	testCases := []struct {
		gitURL       string
		expectedHost string
	}{
		{"git@github.com:user/repo.git", "github.com"},
		{"git@gitlab.com:group/project.git", "gitlab.com"},
		{"git@bitbucket.org:company/repo.git", "bitbucket.org"},
		{"git@custom-git.example.com:org/repo.git", "custom-git.example.com"},
	}

	for _, tc := range testCases {
		result := checkGit(tc.gitURL, "/nonexistent/key", tc.expectedHost, nil, 1)
		if result.Host != tc.expectedHost {
			t.Errorf("For URL %s, expected host=%s, got: %s", tc.gitURL, tc.expectedHost, result.Host)
		}
	}
}
