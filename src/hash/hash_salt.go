// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package hash_salt generates BCrypt Hashes for password config.
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// Use "go run" in the "hash" directory to generate BCrypt Hashes.
func main() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("failed to setup terminal: %v", err)
		return
	}
	t := term.NewTerminal(os.Stdin, "")

	t.SetPrompt("Enter a username: ")
	username, err := t.ReadLine()
	if err != nil {
		fmt.Printf("failed to get password: %v", err)
		term.Restore(fd, oldState)
		return
	}

	pass, err := t.ReadPassword("Enter a password: ")
	if err != nil {
		fmt.Printf("failed to get password: %v", err)
		term.Restore(fd, oldState)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	hashStr := string(hash)
	dockerEscaped := strings.ReplaceAll(hashStr, "$", "$$")

	urlPart := username + ":" + pass
	basicAuth := base64.URLEncoding.EncodeToString([]byte(urlPart))

	term.Restore(fd, oldState)

	fmt.Println("\n\033[1mFormat for docker-compose.yml:\033[0m")
	fmt.Println("environment:")
	fmt.Printf("  - LOG_USERNAME=${LOG_USERNAME:-%s}\n", username)
	fmt.Printf("  - LOG_PASSWORD=${LOG_PASSWORD:-%s}\n", dockerEscaped)

	fmt.Println("\n\033[1mFormat for command line or .env file:\033[0m")
	fmt.Printf("LOG_USERNAME=%s\n", username)
	fmt.Printf("LOG_PASSWORD='%s'\n", hashStr)

	fmt.Println("\n\033[1mFormat for AppEngine app.yaml:\033[0m")
	fmt.Println("env_variables:")
	fmt.Printf("  LOG_USERNAME: %s\n", username)
	fmt.Printf("  LOG_PASSWORD: %s\n", hashStr)

	fmt.Println("\n\033[1mFormat for Auth HTTP header (MCP Clients):\033[0m")
	fmt.Printf("Authorization: Basic %s\n\n", basicAuth)
}
