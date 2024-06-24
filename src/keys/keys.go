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

// Package keys implements crypto key handling.
package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
)

// SigningKey store a key for use
type SigningKey struct {
	Raw any
	Jwk jwk.Key
}

var (
	// Default key size for RSA keys.
	rsaKeySize = 2048

	// Map of keys presented as the IdP's keys.
	keys map[string]*SigningKey

	// Map of keys that are not presented as the IdP's keys.
	wrongKeys map[string]*SigningKey

	// Sync access to keys.
	keysMutex sync.Mutex
)

// SetupKeys sets up neccessary keys.
func SetupKeys() error {
	keysMutex.Lock()
	defer keysMutex.Unlock()

	keys = make(map[string]*SigningKey)
	wrongKeys = make(map[string]*SigningKey)
	registerNoneSigner()

	rsaKey, err := makeRSAKey()
	if err != nil {
		return err
	}

	keys["RSA"] = rsaKey

	wrongRsaKey, err := makeRSAKey()
	if err != nil {
		return err
	}

	wrongKeys["RSA"] = wrongRsaKey

	for _, alg := range []string{"ES256", "ES384", "ES512"} {
		ecdsaKey, err := makeECDSAKey(algToCurve(alg))
		if err != nil {
			return err
		}

		keys[alg] = ecdsaKey

		wrongECDSAKey, err := makeECDSAKey(algToCurve(alg))
		if err != nil {
			return err
		}

		wrongKeys[alg] = wrongECDSAKey
	}

	return nil
}

// GetJSONKeySet returns the IdP's key as a JSON Key Set.
func GetJSONKeySet() (string, error) {
	keysMutex.Lock()
	defer keysMutex.Unlock()

	keySet := jwk.NewSet()
	for _, key := range keys {
		pubKey, err := key.Jwk.PublicKey()
		if err != nil {
			return "", err
		}

		keySet.Add(pubKey)
	}

	buf, err := json.MarshalIndent(keySet, "", "  ")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// GetKey returns a Signing key object by algorithm type.
// if wrongKey is true, return a valid key, but not the one
// that is supposed to be used by the IdP.
func GetKey(alg string, wrongKey bool) *SigningKey {
	keysMutex.Lock()
	defer keysMutex.Unlock()

	keyMap := keys
	if wrongKey {
		keyMap = wrongKeys
	}

	key, ok := keyMap[alg]
	if !ok {
		return nil
	}

	return key
}

// registerNoneSigner registers a JWS signer for the None algorithm. It is a no-op.
func registerNoneSigner() {
	jws.RegisterSigner(jwa.NoSignature, NoneSignerFactory{signer: NoneSigner{}})
}

// NoneSignerFactory creates a NoneSigner.
type NoneSignerFactory struct {
	signer NoneSigner
}

// Create returns a NoneSigner.
func (sf NoneSignerFactory) Create() (jws.Signer, error) {
	return sf.signer, nil
}

// NoneSigner does no signing.
type NoneSigner struct {
}

// Sign is a no-op.
func (s NoneSigner) Sign([]byte, any) ([]byte, error) {
	return nil, nil
}

// Algorithm returns the NoSignature type.
func (s NoneSigner) Algorithm() jwa.SignatureAlgorithm {
	return jwa.NoSignature
}

// makeRSAKey makes an RSA Key.
func makeRSAKey() (*SigningKey, error) {
	raw, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		fmt.Printf("failed to generate new RSA private key: %s\n", err)
		return nil, err
	}

	key, err := jwk.New(raw)
	if err != nil {
		fmt.Printf("failed to create jwt key: %s\n", err)
		return nil, err
	}

	key.Set("use", "sig")
	key.Set("kid", "customidprsa")
	return &SigningKey{Raw: raw, Jwk: key}, nil
}

// makeECDSAKey made an ECDSA key.
func makeECDSAKey(curve elliptic.Curve) (*SigningKey, error) {
	raw, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Printf("failed to generate new ECDSA private key: %s\n", err)
		return nil, err
	}

	key, err := jwk.New(raw)
	if err != nil {
		fmt.Printf("failed to create jwt key: %s\n", err)
		return nil, err
	}

	key.Set("use", "sig")
	key.Set("kid", "customidpecdsa"+curve.Params().Name)
	return &SigningKey{Raw: raw, Jwk: key}, nil
}

func algToCurve(alg string) elliptic.Curve {
	switch alg {
	case "ES256":
		return elliptic.P256()
	case "ES384":
		return elliptic.P384()
	case "ES512":
		return elliptic.P521()
	}

	fmt.Printf("Invalid ECDSA alg specified %s", alg)
	return elliptic.P256()
}
