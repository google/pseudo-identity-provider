---
name: pseudo-idp
description: >
  Use this skill whenever the user wants to do OAuth2 or OIDC security testing using the Pseudo Identity Provider (Pseudo IdP) MCP server. Triggers include: testing OAuth security, simulating identity provider attacks, testing CSRF defenses, JWT validation, nonce checks, token signature attacks, Mix-up attacks, PKCE testing, expired tokens, invalid issuers, or any scenario involving a fake/mock IdP. Also trigger when the user mentions "Pseudo IdP", "pseudo-identity-provider", "set_config", "get_config", "list_logs" in a security testing context, or asks to configure OAuth misbehavior for client testing. Always use this skill when the user wants to verify their OAuth/OIDC client's security posture by making the IdP behave incorrectly.
---

# Pseudo Identity Provider — OAuth/OIDC Security Testing

Pseudo IdP is a **fake OAuth2/OIDC Identity Provider** you control via MCP. It exposes real OIDC endpoints (Discovery, Authorization, Token, UserInfo) but lets you configure malformed or malicious behavior to verify your client's security checks.

## MCP Tools Available

Three tools are exposed by the `Psuedo Idp` MCP connector:

| Tool | Purpose |
|------|---------|
| `Psuedo Idp:get_config` | Fetch the current full configuration as JSON |
| `Psuedo Idp:set_config` | Update one or more endpoint configurations |
| `Psuedo Idp:list_logs` | View recent request/response logs |

**Always call `get_config` first** before making changes so you have the baseline and can restore it.

---

## Core Concepts

### Endpoints and Their Config Keys

| Endpoint | URL path | Config key in `set_config` |
|----------|----------|---------------------------|
| Authorization | `/oauth2/auth` | `auth_action` |
| Token | `/oauth2/token` | `token_action` |
| UserInfo | `/oauth2/userinfo` | `userinfo_action` |
| Discovery | `/.well-known/openid-configuration` | `discovery_action` |
| Callback | `/callback` | `callback_action` |
| ID Token (JWT) | (embedded in token response) | `id_token_config` |

### Endpoint Action Types

Each endpoint (except `id_token_config`) uses an `action_type` field:

- **`redirect`** (auth only) — Redirect the browser to a callback URI with configurable parameters
- **`respond`** — Return a successful JSON or HTTP response
- **`error`** — Return a specific HTTP error code
- **`block`** — Hang and time out (useful for testing timeout handling)

### Parameter Actions

Inside `parameters` arrays or `claims` arrays:

| Action | Behaviour |
|--------|-----------|
| `passthrough` | Keep the value from the inbound request |
| `set` | Set an explicit value (supports Go templates — see Level 2) |
| `omit` | Remove the parameter from the response |
| `random` | Set to a random base64 string |
| `custom` | Invoke a named custom Go processor (e.g., `signed_token_id` for JWT generation) |

---

## Workflow Pattern

```
1. get_config          ← snapshot baseline
2. set_config(...)     ← inject the misbehaviour
3. [run client flow]   ← trigger the OAuth flow from your client
4. list_logs           ← verify what the client sent and what IdP returned
5. set_config(...)     ← restore to baseline (or next test)
```

---

## Quick-Start Config Recipes

These are the most common test scenarios. For full JSON configs and more configs, see **`references/config-recipes.md`**.

### 1. Remove JWT Signature (CVE-2020-28042)
```json
{ "id_token_config": { "remove_signature": true } }
```
Verify your client rejects tokens with no signature segment.

### 2. alg:none Attack (CVE-2015-9235)
```json
{ "id_token_config": { "alg": "none" } }
```
Verify your client blocks the `none` algorithm.

### 3. Wrong Signing Key
```json
{ "id_token_config": { "use_wrong_key": true } }
```
Signs with a key not in `/.well-known/jwks.json`. Client must reject.

### 4. CSRF — Replace State Parameter
```json
{
  "auth_action": {
    "action_type": "redirect",
    "redirect": {
      "default_parameter_action": "passthrough",
      "parameters": [{ "id": "state", "action": "set", "values": ["tampered_state"] }],
      "redirect_target": { "use_custom_redirect_uri": false, "target": "", "custom_key": "" },
      "use_hash_fragment": false
    }
  }
}
```
Verify your client rejects mismatched `state` (Login CSRF protection).

### 5. Expired Token
Set the `exp` claim to yesterday:
```json
{
  "id_token_config": {
    "claims": [
      { "id": "exp", "values": ["{{with $yesterday := .Time.AddDate 0 0 -1}}{{$yesterday.Unix}}{{end}}"], "json_type": "number" }
    ]
  }
}
```

### 6. Invalid Nonce
```json
{
  "id_token_config": {
    "claims": [{ "id": "nonce", "values": ["wrong_nonce"], "json_type": "string" }]
  }
}
```

### 7. Invalid Issuer
```json
{
  "id_token_config": {
    "claims": [{ "id": "iss", "values": ["https://evil.example.com"], "json_type": "string" }]
  }
}
```

---

## Reading Logs

After running a test, call `list_logs` to inspect what happened. Logs show:
- Inbound request to each endpoint (method, headers, parameters)
- What the IdP returned

Use logs to confirm the client sent the right parameters (e.g., `code_verifier` for PKCE) and to debug unexpected behavior.

---

## Restoring to Baseline

After testing, restore the original config. Best practice: call `get_config` at the start and save the JSON; then pass it back verbatim via `set_config` to restore.

---

## Going Deeper

For advanced scenarios, full JSON schemas, Mix-up attack configs, PKCE bypass, implicit flow setup, and templating reference, read:
→ **`references/config-recipes.md`** — full config library with complete JSON configs
→ **`references/schema-reference.md`** — complete `set_config` field reference