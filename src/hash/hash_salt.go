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
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// Use "go run" in the "hash" directory to generate BCrypt Hashes.
func main() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Errorf("failed to setup terminal: %v", err)
		return
	}
	t := term.NewTerminal(os.Stdin, "")

	pass, err := t.ReadPassword("Enter a password:")
	if err != nil {
		fmt.Errorf("failed to get password: %v", err)
		term.Restore(fd, oldState)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	term.Restore(fd, oldState)

	fmt.Println("Hash: " + string(hash))
}
