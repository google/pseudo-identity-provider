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

// Package config implementes configuration and parameter construction.
package config

import "sync"

// Config represents the overall configuration.
type Config struct {
	AuthAction      AuthAction      `json:"auth_action" jsonschema:"title=Authorization Endpoint Configuration"`
	TokenAction     TokenAction     `json:"token_action" jsonschema:"title=Token Endpoint Configuration"`
	UserInfoAction  UserInfoAction  `json:"userinfo_action" jsonschema:"title=UserInfo Endpoint Configuration"`
	DiscoveryAction DiscoveryAction `json:"discovery_action" jsonschema:"title=Discovery Endpoint Configuration"`

	// Custom Parameter Config Entries.
	IDTokenConfig IDTokenConfig `json:"id_token_config" jsonschema:"title=ID Token Config"`
}

// AuthAction configures the authz endpoint.
type AuthAction struct {
	Action   string       `json:"action_type" jsonschema:"title=Endpoint Action,enum=redirect,enum=error,enum=block,default=redirect"`
	Redirect AuthRedirect `json:"redirect" jsonschema:"title=Redirect Config" jsonschema_extras:"hide=action_type !== redirect"`
	Error    Error        `json:"error" jsonschema:"title=Error Config" jsonschema_extras:"hide=action_type !== error"`
	// Block doesn't have any parameters.
}

// AuthRedirect configures a redirection.
type AuthRedirect struct {
	RedirectTarget     RedirectTarget `json:"redirect_target" jsonschema:"title=Redirect Target"`
	DefaultParamAction string         `json:"default_parameter_action" jsonschema:"title=Default Parameter Action,enum=passthrough,enum=omit"`
	Parameters         []Parameter    `json:"parameters" jsonschema:"title=Parameters"`
	UseHashFragment    bool           `json:"use_hash_fragment" jsonschema:"title=Use Hash Fragment"`
}

// RedirectTarget is the target of a redirection.
type RedirectTarget struct {
	UseCustomRedirectURI bool   `json:"use_custom_redirect_uri" jsonschema:"title=Use custom redirect"`
	Target               string `json:"target" jsonschema:"title=Target URL"`
	CustomKey            string `json:"custom_key" jsonschema:"title=Custom processor key"`
}

// Parameter is an input or output parameter to an endpoint that can be represented in a URL
// query parameter or form depending on the context.
type Parameter struct {
	ID        string   `json:"id" jsonschema:"title=Identifier of the parameter,default=id"`
	Action    string   `json:"action" jsonschema:"title=Parameter Action,enum=passthrough,enum=set,enum=omit,enum=random,enum=custom,default=passthrough"`
	Values    []string `json:"values" jsonschema:"title=Values,default=example_value" jsonschema_extras:"hide=action !== set"`
	CustomKey string   `json:"custom_key" jsonschema:"title=Custom Processor Key,default=test" jsonschema_extras:"hide=action !== custom"`
	JSONType  string   `json:"json_type" jsonschema:"title=JSON Value Type,enum=string,enum=array,enum=number,enum=boolean,enum=object,default=string"`
}

// Error represents an error to return.
type Error struct {
	ErrorCode    int    `json:"error_code" jsonschema:"title=Error Code"`
	ErrorContent string `json:"error_content" jsonschema:"title=Error Content"`
}

// TokenAction configures the Token endpoint.
type TokenAction struct {
	Action  string       `json:"action_type" jsonschema:"title=Token Endpoint Action,enum=respond,enum=error,enum=block"`
	Respond TokenRespond `json:"respond" jsonschema:"title=Response Config" jsonschema_extras:"hide=action_type !== respond"`
	// TODO Implement Forward TokenForward  `json:"forward"`
	Error Error `json:"error" jsonschema:"title=Error Config" jsonschema_extras:"hide=action_type !== 'error'"`
	// Block doesn't have any parameters.
}

// TokenRespond configures responding with JSON content.
type TokenRespond struct {
	Parameters []Parameter `json:"parameters" jsonschema:"title=Parameters"`
}

// IDTokenConfig configures IDToken responses.
type IDTokenConfig struct {
	Algorithm       string  `json:"alg" jsonschema:"title=JWT Signature Algorithm"`
	RemoveSignature bool    `json:"remove_signature" jsonschema:"title=Remove Signature"`
	UseWrongKey     bool    `json:"use_wrong_key" jsonschema:"title=Use Incorrect Key"`
	Claims          []Claim `json:"claims" jsonschema:"title=Claims"`
}

// Claim represents an IDToken claim.
type Claim struct {
	ID       string   `json:"id" jsonschema:"title=Claim ID"`
	Values   []string `json:"values" jsonschema:"title=Claim Values"`
	JSONType string   `json:"json_type" jsonschema:"title=JSON Type,enum=string,enum=array,enum=number,enum=boolean,enum=object,default=string"`
}

// DiscoveryAction configures the Discovery endpoint.
type DiscoveryAction struct {
	Action  string           `json:"action_type" jsonschema:"title=Discovery Endpoint Action,enum=respond,enum=error,enum=block"`
	Respond DiscoveryRespond `json:"respond" jsonschema:"title=Response Config" jsonschema_extras:"hide=action_type !== respond"`
	Error   Error            `json:"error" jsonschema:"title=Error Config" jsonschema_extras:"hide=action_type !== error"`
	// Block doesn't have any parameters.
}

// DiscoveryRespond configures the Discovery endpoint response of JSON content.
type DiscoveryRespond struct {
	Parameters []Parameter `json:"parameters" jsonschema:"title=Parameters"`
}

// UserInfoAction configures the UserInfo endpoint.
type UserInfoAction struct {
	Action  string          `json:"action_type" jsonschema:"title=User Info Endpoint Action,enum=respond,enum=error,enum=block"`
	Respond UserInfoRespond `json:"respond" jsonschema:"title=Response Config" jsonschema_extras:"hide=action_type !== respond"`
	Error   Error           `json:"error" jsonschema:"title=Error Config" jsonschema_extras:"hide=action_type !== error"`
	// Block doesn't have any parameters.
}

// UserInfoRespond configures the Userinfo endpoint response of JSON content.
type UserInfoRespond struct {
	Parameters []Parameter `json:"parameters" jsonschema:"title=Parameters"`
}

// DefaultConfig is the config present on first startup. It acts as a default OIDC IDP using
// authorization code flow and returning a static subject in the ID Token.
var DefaultConfig = Config{
	AuthAction: AuthAction{
		Action: "redirect",
		Redirect: AuthRedirect{
			DefaultParamAction: "passthrough",
			Parameters: []Parameter{
				{ID: "code", Action: "random", JSONType: "string"},
				{ID: "redirect_uri", Action: "omit", JSONType: "string"},
			},
		},
	},
	TokenAction: TokenAction{
		Action: "respond",
		Respond: TokenRespond{
			Parameters: []Parameter{
				{ID: "id_token", Action: "custom", CustomKey: "signed_token_id", JSONType: "string"},
				{ID: "access_token", Action: "random", JSONType: "string"},
				{ID: "refresh_token", Action: "random", JSONType: "string"},
				{ID: "expires_in", Action: "set", JSONType: "number", Values: []string{"3600"}},
				{ID: "token_type", Action: "set", Values: []string{"Bearer"}, JSONType: "string"},
			},
		},
	},
	UserInfoAction: UserInfoAction{
		Action: "respond",
		Respond: UserInfoRespond{
			Parameters: []Parameter{
				{ID: "sub", Action: "set", Values: []string{"12345abcde"}, JSONType: "string"},
				{ID: "email", Action: "set", Values: []string{"testsub@{{.Domain}}"}, JSONType: "string"},
			},
		},
	},
	DiscoveryAction: DiscoveryAction{
		Action: "respond",
		Respond: DiscoveryRespond{
			Parameters: []Parameter{
				{ID: "issuer", Action: "set", Values: []string{"https://{{.Domain}}"}, JSONType: "string"},
				{ID: "authorization_endpoint", Action: "set", Values: []string{"https://{{.Domain}}/oauth2/auth"}, JSONType: "string"},
				{ID: "token_endpoint", Action: "set", Values: []string{"https://{{.Domain}}/oauth2/token"}, JSONType: "string"},
				{ID: "userinfo_endpoint", Action: "set", Values: []string{"https://{{.Domain}}/oauth2/userinfo"}, JSONType: "string"},
				{ID: "jwks_uri", Action: "set", Values: []string{"https://{{.Domain}}/.well-known/jwks.json"}, JSONType: "string"},
				{ID: "subject_types_supported", Action: "set", JSONType: "array", Values: []string{"public"}},
				{ID: "id_token_signing_alg_values_supported", Action: "set", JSONType: "array", Values: []string{"RS256", "RS512", "ES256"}},
				{
					ID:       "response_types_supported",
					Action:   "set",
					JSONType: "array",
					Values: []string{
						"code", "code id_token", "id_token", "token id_token", "token", "token id_token code",
					},
				},
			},
		},
	},
	IDTokenConfig: IDTokenConfig{
		Algorithm: "RS256",
		Claims: []Claim{
			{ID: "iss", Values: []string{"https://{{.Domain}}"}, JSONType: "string"},
			{ID: "aud", Values: []string{"{{if .Session}}{{.Session.ClientID}}{{end}}"}, JSONType: "string"},
			{ID: "nonce", Values: []string{"{{if .Session}}{{.Session.Nonce}}{{end}}"}, JSONType: "string"},
			{ID: "iat", JSONType: "number", Values: []string{"{{.Time.Unix}}"}},
			{ID: "exp", JSONType: "number", Values: []string{"{{with $tomorrow := .Time.AddDate 0 0 1}}{{$tomorrow.Unix}}{{end}}"}},
			{ID: "sub", Values: []string{"12345abcde"}, JSONType: "string"},
		},
		RemoveSignature: false,
		UseWrongKey:     false,
	},
}

// Config storage.
var globalConfig Config
var configMutex sync.Mutex

func init() {
	globalConfig = DefaultConfig
}

// GetGlobalConfig get a copy of the global config. Thread-safe.
func GetGlobalConfig() *Config {
	configMutex.Lock()
	defer configMutex.Unlock()
	c := new(Config)
	*c = globalConfig
	return c
}

// SetGlobalConfig sets the global config. Thread-safe.
func SetGlobalConfig(config *Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	globalConfig = *config
}
