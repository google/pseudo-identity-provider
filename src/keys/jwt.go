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

package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// SignToken creates a signature for token content serialized to base64.
// Id wrongKey is true, it will use a valid key, but it won't be in the
// IdP's set of JSON keys.
func SignToken(alg string, token jwt.Token, wrongKey bool) (string, error) {
	var sigAlg jwa.SignatureAlgorithm
	if err := sigAlg.Accept(alg); err != nil {
		return "", fmt.Errorf("invalid algorithm %q: %s", alg, err)
	}

	var signed string
	var key any
	switch alg {
	case "RS256", "RS384", "RS512", "none":
		key = GetKey("RSA", wrongKey).Jwk
	case "ES256", "ES384", "ES512":
		key = GetKey(alg, wrongKey).Jwk
	case "HS256":
		// Use the RSA Public key as the HMAC Secret to make it easy to test for CVE-2016-5431.
		privKey := GetKey("RSA", wrongKey).Raw.(*rsa.PrivateKey)
		key = publicKeyToBytes(&privKey.PublicKey)
	default:
		return "", fmt.Errorf("specified key %s not supported", alg)
	}

	signedBytes, err := jwt.Sign(token, sigAlg, key)
	if err != nil {
		log.Printf("failed to sign token: %s", err)
		return "", fmt.Errorf("failed to sign token: %s", err)
	}
	signed = string(signedBytes)

	return signed, nil
}

// publicKeyToBytes gets a RSA Public key as a byte array.
func publicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Print(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}
