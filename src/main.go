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

// Package oidcservice runs the IdP in an AppEngine Container
package main

import (
	"html/template"
	"log"
	"os"

	idp "customidp/idp"

	"google.golang.org/appengine/v2"
)

// Static HTML file templates.
var templates = template.Must(template.ParseGlob("static/browser/*.html"))

// Main setups handlers and starts the service for AppEngine hosting.
func main() {
	idp.InitHandlers(templates, "static/browser")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	appengine.Main()
}

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}
