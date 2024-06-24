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

package config

import (
	"customidp/keys"
	sessionmgmt "customidp/session"
	"strings"

	"github.com/lestrrat-go/jwx/jwt"
)

func init() {
	RegisterCustomParam("signed_token_id", GenerateToken)
}

// GenerateToken creates a JWT token based on the IDTokenConfig.
func GenerateToken(input *sessionmgmt.RequestInput, config *Config) ([]string, error) {
	token := jwt.New()
	for _, claim := range config.IDTokenConfig.Claims {
		p := Parameter{
			ID:       claim.ID,
			Action:   "set",
			Values:   claim.Values,
			JSONType: claim.JSONType,
		}

		jsonVal, err := p.GetJSON(input)
		if err != nil {
			return nil, err
		}

		if jsonVal != nil {
			token.Set(claim.ID, jsonVal)
		}
	}

	signed, err := keys.SignToken(config.IDTokenConfig.Algorithm, token, config.IDTokenConfig.UseWrongKey)
	if err != nil {
		return nil, err
	}

	if config.IDTokenConfig.RemoveSignature {
		i := strings.LastIndex(signed, ".")
		if i >= 0 && i < len(signed)-1 {
			signed = signed[:i+1]
		}
	}

	return []string{signed}, nil
}
