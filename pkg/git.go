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
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// parseGitDSN parses a git+ssh DSN and extracts the repository URL, key file path, and host.
func parseGitDSN(dsnString string) (gitURL, keyFile, host string, err error) {
	parsed, err := url.Parse(dsnString)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse git DSN: %w", err)
	}

	// Extract keyfile parameter from query string
	keyFile = parsed.Query().Get("keyfile")
	if keyFile == "" {
		return "", "", "", fmt.Errorf("git+ssh DSN requires keyfile query parameter")
	}

	// Host is already parsed by url.Parse
	host = parsed.Host

	// Reconstruct the git URL without the query parameters
	// Convert git+ssh:// to git@ format for go-git
	// Standard Git SSH URL format uses colon, not slash: git@github.com:user/repo.git
	gitURL = fmt.Sprintf("git@%s:%s", host, strings.TrimPrefix(parsed.Path, "/"))

	return gitURL, keyFile, host, nil
}

// connectToGitRepo tests a connection to a Git repository using SSH authentication.
// It performs a shallow clone (depth=1) to a temporary directory and cleans up afterwards.
func connectToGitRepo(repoURL, keyPath string) error {
	// Validate key file exists
	if _, err := os.Stat(keyPath); err != nil {
		return fmt.Errorf("SSH key file not accessible: %w", err)
	}

	// Set up SSH authentication with the provided key file
	auth, err := ssh.NewPublicKeysFromFile("git", keyPath, "")
	if err != nil {
		return fmt.Errorf("failed to load SSH key: %w", err)
	}

	// Accept any host key for connection testing
	// This is similar to SSH's StrictHostKeyChecking=no
	// Appropriate for a connection testing tool that doesn't persist data
	auth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

	// Create a temporary directory for the clone
	tempDir, err := os.MkdirTemp("", "conntest-git-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Perform a shallow clone to test the connection
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      repoURL,
		Auth:     auth,
		Depth:    1,
		Progress: os.Stderr,
	})

	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}
