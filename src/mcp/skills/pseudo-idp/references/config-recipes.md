# Pseudo IdP Attack Recipes

Full JSON configurations for OAuth2/OIDC security testing scenarios.

## Table of Contents

1. [JWT Signature Attacks](#1-jwt-signature-attacks)
2. [Token Claim Attacks](#2-token-claim-attacks)
3. [CSRF / State Parameter Attacks](#3-csrf--state-parameter-attacks)
4. [Mix-Up Attacks](#4-mix-up-attacks)
5. [Endpoint-Level Attacks](#5-endpoint-level-attacks)
6. [PKCE Testing](#6-pkce-testing)
7. [Implicit Flow Setup](#7-implicit-flow-setup)
8. [Callback Endpoint Testing](#8-callback-endpoint-testing)
9. [Restore to Default Config](#9-restore-to-default-config)

---

## 1. JWT Signature Attacks

### 1a. Remove Signature — CVE-2020-28042
Strips the signature segment from the JWT. Client must reject.
```json
{ "id_token_config": { "remove_signature": true } }
```

### 1b. alg:none — CVE-2015-9235
Sets the algorithm header to `none`. Client must reject.
```json
{ "id_token_config": { "alg": "none" } }
```

### 1c. Wrong Signing Key
Signs with a valid key type but one NOT listed in `/.well-known/jwks.json`. Client must reject.
```json
{ "id_token_config": { "use_wrong_key": true } }
```

### 1d. Key Confusion — CVE-2016-5431 (HS256 with RSA pubkey as secret)
Uses HS256, which causes some libraries to accept the RSA public key as the HMAC secret.
```json
{ "id_token_config": { "alg": "HS256" } }
```

---

## 2. Token Claim Attacks

All of these modify `id_token_config.claims`. When setting specific claims, include the full existing claims array with your modification, or only include the claim(s) you're changing (the server merges by claim ID).

### 2a. Expired Token
```json
{
  "id_token_config": {
    "claims": [
      {
        "id": "exp",
        "values": ["{{with $yesterday := .Time.AddDate 0 0 -1}}{{$yesterday.Unix}}{{end}}"],
        "json_type": "number"
      }
    ]
  }
}
```

### 2b. Invalid Nonce
```json
{
  "id_token_config": {
    "claims": [{ "id": "nonce", "values": ["wrong_nonce_value"], "json_type": "string" }]
  }
}
```

### 2c. Wrong Issuer
```json
{
  "id_token_config": {
    "claims": [{ "id": "iss", "values": ["https://attacker.example.com"], "json_type": "string" }]
  }
}
```

### 2d. Wrong Audience
```json
{
  "id_token_config": {
    "claims": [{ "id": "aud", "values": ["wrong_client_id"], "json_type": "string" }]
  }
}
```

### 2e. Missing Subject
```json
{
  "id_token_config": {
    "claims": [{ "id": "sub", "values": [""], "json_type": "string" }]
  }
}
```

---

## 3. CSRF / State Parameter Attacks

### 3a. Replace State with Tampered Value
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "default_parameter_action": "passthrough",
      "parameters": [{ "id": "state", "action": "set", "values": ["tampered_state_value"] }],
      "redirect_target": { "use_custom_redirect_uri": false, "target": "", "custom_key": "" },
      "use_hash_fragment": false
    }
  }
}
```

### 3b. Omit State Entirely
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "default_parameter_action": "passthrough",
      "parameters": [{ "id": "state", "action": "omit" }],
      "redirect_target": { "use_custom_redirect_uri": false, "target": "", "custom_key": "" },
      "use_hash_fragment": false
    }
  }
}
```

### 3c. Random State (replay-like noise)
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "default_parameter_action": "passthrough",
      "parameters": [{ "id": "state", "action": "random" }],
      "redirect_target": { "use_custom_redirect_uri": false, "target": "", "custom_key": "" },
      "use_hash_fragment": false
    }
  }
}
```

---

## 4. Mix-Up Attacks

Mix-up attacks redirect the client to send an auth code to the wrong IdP.

### 4a. Simple Mix-Up (redirect auth to honest IdP)
Replace `https://honestidp/auth` with the real IdP's auth endpoint URL.
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "redirect_target": {
        "use_custom_redirect_uri": true,
        "target": "https://honestidp/auth",
        "custom_key": ""
      },
      "default_parameter_action": "passthrough",
      "parameters": [
        { "id": "client_id", "action": "set", "values": ["clientidathonestidp"] },
        { "id": "nonce", "action": "omit" }
      ],
      "use_hash_fragment": false
    }
  },
  "token_action": { "action_type": "block" }
}
```

### 4b. Chosen Nonce Mix-Up
Attacker pre-fetches a nonce from the honest IdP and injects it into the Pseudo IdP flow.
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "redirect_target": {
        "use_custom_redirect_uri": true,
        "target": "https://honestidp/auth",
        "custom_key": ""
      },
      "default_parameter_action": "passthrough",
      "parameters": [
        { "id": "client_id", "action": "set", "values": ["clientidathonestidp"] },
        { "id": "nonce", "action": "set", "values": ["attacker_chosen_nonce"] }
      ],
      "use_hash_fragment": false
    }
  },
  "token_action": { "action_type": "block" }
}
```

---

## 5. Endpoint-Level Attacks

### 5a. Block Token Endpoint (timeout simulation)
```json
{ "token_action": { "action_type": "block" } }
```

### 5b. Error from Authorization Endpoint (e.g., 401)
```json
{
  "auth_action": {
    "action_type": "error",
    "error": { "error_code": 401, "error_content": "Unauthorized" }
  }
}
```

### 5c. Error from Discovery Endpoint
```json
{
  "discovery_action": {
    "action_type": "error",
    "error": { "error_code": 500, "error_content": "Internal Server Error" }
  }
}
```

### 5d. Omit Token Endpoint URL from Discovery
Forces the client to fall back to default or fail discovery-based setup.
```json
{
  "discovery_action": {
    "action_type": "respond",
    "respond": {
      "parameters": [
        { "id": "token_endpoint", "action": "omit" }
      ]
    }
  }
}
```

---

## 6. PKCE Testing

PKCE (Proof Key for Code Exchange) parameters arrive at the Token endpoint. Use `list_logs` to verify:
- `code_verifier` is present in the token request
- The server's `code_challenge` was set during the auth request

### 6a. Verify PKCE Parameters Were Sent
Run a normal flow, then call `list_logs` and check:
```
Token request body should contain: code_verifier=<value>
Auth request URL should contain: code_challenge=<value>&code_challenge_method=S256
```

### 6b. Test PKCE Downgrade (omit code_challenge from auth response passthrough)
Modify the auth action to omit `code_challenge` from the session, then verify the client rejects the token exchange or requires it.

---

## 7. Implicit Flow Setup

Configure the Authorization endpoint to return tokens directly in the redirect (hash fragment).
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "redirect_target": { "use_custom_redirect_uri": false, "target": "", "custom_key": "" },
      "default_parameter_action": "passthrough",
      "parameters": [
        { "id": "id_token", "action": "custom", "custom_key": "signed_token_id" },
        { "id": "access_token", "action": "random" },
        { "id": "token_type", "action": "set", "values": ["Bearer"] },
        { "id": "expires_in", "action": "set", "values": ["3600"] }
      ],
      "use_hash_fragment": true
    }
  }
}
```

---

## 8. Callback Endpoint Testing

The `/callback` endpoint can receive arbitrary redirects — useful for testing authorization code leakage.

### 8a. Log Inbound Requests at /callback
Default respond action — just call `list_logs` after triggering to see what arrived.
```json
{
  "callback_action": {
    "action_type": "respond",
    "respond": {
      "body": { "action": "set", "value": "received" },
      "headers": [{ "key": "Content-Type", "value": "text/plain" }]
    }
  }
}
```

### 8b. Redirect from Callback
Useful for open redirect testing on the client side.
```json
{
  "callback_action": {
    "action_type": "redirect",
    "redirect": { "target": "https://attacker.example.com/steal" }
  }
}
```

---

## 9. Restore to Default Config

After testing, restore the safe default. Call `get_config` before starting tests and store the result. To restore, pass the original JSON back via `set_config`. The default safe baseline includes:

- `auth_action.action_type`: `redirect`, passthrough all parameters, random `code`
- `token_action.action_type`: `respond`, returns `id_token` via `signed_token_id` custom processor
- `id_token_config.alg`: `RS256`, `remove_signature`: false, `use_wrong_key`: false
- `discovery_action.action_type`: `respond` with templated endpoint URLs
- `userinfo_action.action_type`: `respond` with `sub` and `email`

A minimal "reset to safe" call:
```json
{
  "id_token_config": {
    "alg": "RS256",
    "remove_signature": false,
    "use_wrong_key": false
  },
  "auth_action": { "action_type": "redirect" },
  "token_action": { "action_type": "respond" },
  "discovery_action": { "action_type": "respond" },
  "userinfo_action": { "action_type": "respond" }
}
```