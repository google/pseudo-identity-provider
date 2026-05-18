# Pseudo IdP `set_config` Schema Reference

Complete field reference for the `Psuedo Idp:set_config` MCP tool.

## Top-Level Fields

| Field | Type | Description |
|-------|------|-------------|
| `auth_action` | object | Authorization endpoint (`/oauth2/auth`) behavior |
| `token_action` | object | Token endpoint (`/oauth2/token`) behavior |
| `userinfo_action` | object | UserInfo endpoint (`/oauth2/userinfo`) behavior |
| `discovery_action` | object | Discovery endpoint (`/.well-known/openid-configuration`) behavior |
| `callback_action` | object | Callback endpoint (`/callback`) behavior |
| `id_token_config` | object | JWT ID Token generation settings |

All fields are optional — only provide what you want to change.

---

## `auth_action`

Controls the Authorization redirect endpoint.

```
auth_action
├── action_type: "redirect" | "error" | "block"  (default: "redirect")
├── error (visible when action_type === "error")
│   ├── error_code: integer    HTTP status code
│   └── error_content: string  Response body
└── redirect (visible when action_type === "redirect")
    ├── redirect_target
    │   ├── use_custom_redirect_uri: boolean  Use target URL instead of redirect_uri param
    │   ├── target: string                    Custom redirect URL
    │   └── custom_key: string                Custom processor key for target evaluation
    ├── default_parameter_action: "passthrough" | "omit"
    ├── use_hash_fragment: boolean             Put params in URL hash fragment (implicit flow)
    └── parameters: [ParameterConfig]
```

---

## `token_action`

Controls the Token endpoint JSON response.

```
token_action
├── action_type: "respond" | "error" | "block"
├── error: ErrorConfig
└── respond
    └── parameters: [ParameterConfig]
```

Default parameters returned: `id_token`, `access_token`, `refresh_token`, `expires_in`, `token_type`.

---

## `userinfo_action`

Controls the UserInfo endpoint JSON response.

```
userinfo_action
├── action_type: "respond" | "error" | "block"
├── error: ErrorConfig
└── respond
    └── parameters: [ParameterConfig]
```

Default parameters: `sub`, `email`.

---

## `discovery_action`

Controls the OIDC discovery document (`/.well-known/openid-configuration`).

```
discovery_action
├── action_type: "respond" | "error" | "block"
├── error: ErrorConfig
└── respond
    └── parameters: [ParameterConfig]
```

Default parameters include all standard OIDC discovery fields using `{{.Domain}}` templates.

---

## `callback_action`

Controls the `/callback` endpoint (and sub-paths). Useful for logging arbitrary redirects.

```
callback_action
├── action_type: "respond" | "redirect" | "error" | "block"
├── error: ErrorConfig
├── redirect
│   └── target: string
└── respond
    ├── body
    │   ├── action: "set" | "custom"
    │   ├── value: string     (when action === "set"; supports templates)
    │   └── custom_key: string
    └── headers: [{ key: string, value: string }]
```

---

## `id_token_config`

Controls JWT generation for the ID Token.

```
id_token_config
├── alg: string    "RS256"|"RS384"|"RS512"|"ES256"|"ES384"|"ES512"|"HS256"|"none"
├── remove_signature: boolean    Strip the JWT signature segment
├── use_wrong_key: boolean       Sign with key not in jwks.json
└── claims: [ClaimConfig]
```

### ClaimConfig

```
{
  "id": string,           Claim name (e.g., "iss", "sub", "exp", "nonce")
  "values": [string],     One or more values (supports Go templates)
  "json_type": "string" | "number" | "boolean" | "array" | "object"
}
```

Default claims: `iss`, `aud`, `nonce`, `iat`, `exp`, `sub`.

---

## ParameterConfig (used in endpoint respond/redirect arrays)

```
{
  "id": string,        Parameter name
  "action": "passthrough" | "set" | "omit" | "random" | "custom",
  "values": [string],  Required when action === "set"; supports Go templates
  "custom_key": string, Required when action === "custom"
  "json_type": "string" | "array" | "number" | "boolean" | "object"
}
```

---

## ErrorConfig

```
{
  "error_code": integer,    HTTP status code (e.g., 400, 401, 500)
  "error_content": string   HTTP response body
}
```

---

## Go Template Variables

Available in `values` fields with `action: "set"` and in claim `values`:

| Variable | Type | Description |
|----------|------|-------------|
| `{{.Domain}}` | string | IdP server domain |
| `{{.HTTPMethod}}` | string | HTTP method of the request |
| `{{.Path}}` | string | URL path |
| `{{.Proto}}` | string | Protocol (http/https) |
| `{{.Headers}}` | []Header | HTTP request headers |
| `{{.URLParams}}` | map | URL query parameters |
| `{{.FormParams}}` | map | POST form parameters |
| `{{.Session.Code}}` | string | Auth code |
| `{{.Session.Nonce}}` | string | OIDC nonce |
| `{{.Session.CodeChallenge}}` | string | PKCE code challenge |
| `{{.Session.CodeChallengeMethod}}` | string | PKCE method (e.g., S256) |
| `{{.Session.ClientID}}` | string | OAuth client ID |
| `{{.Session.RedirectURI}}` | string | Requested redirect URI |
| `{{.Time}}` | time.Time | Request timestamp (Go time.Time) |

### Template Examples

```
# Current domain
https://{{.Domain}}/oauth2/token

# Conditional nonce (only if session exists)
{{if .Session}}{{.Session.Nonce}}{{end}}

# Token expiry = tomorrow
{{with $tomorrow := .Time.AddDate 0 0 1}}{{$tomorrow.Unix}}{{end}}

# Token expiry = yesterday (expired)
{{with $yesterday := .Time.AddDate 0 0 -1}}{{$yesterday.Unix}}{{end}}

# Token issued at now
{{.Time.Unix}}
```

---

## Custom Processor Keys

Built-in custom processor registered in the Go backend:

| Key | Description |
|-----|-------------|
| `signed_token_id` | Generates and signs a JWT ID token using `id_token_config` settings |

Custom processors can be added by extending the Go backend source in `src/config/`.