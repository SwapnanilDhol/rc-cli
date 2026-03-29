---
name: revenuecat-cli
description: >-
  Work with the rc-cli repo or embed RevenueCat behavior: choose Public API v2
  (Bearer api key, api.revenuecat.com/v2) vs Internal dashboard API (rc_auth_token
  cookie, app.revenuecat.com/internal/v1). Use when the user mentions RevenueCat,
  rc-cli, API keys, rc login, internal API, or Postman collections for this project.
---

# RevenueCat CLI — which API am I using?

## Two systems (orthogonal)

| Choose | Auth | Base | Typical need |
|--------|------|------|----------------|
| **V2 public API** | `apiKey` in `~/.revenuerc`, `Bearer` header | `api.revenuecat.com/v2` | Documented REST, automation with secret key |
| **Internal API** | `rc login` → `authToken` / cookie `rc_auth_token` | `app.revenuecat.com/internal/v1` | All projects, dashboard CRUD, `rc internal` |

They are **not** interchangeable. A project API key does **not** replace a dashboard session for internal routes.

## Commands
- **V2:** `rc api GET '/projects/...'`, `rc config` for `apiKey`
- **Internal:** `rc login`, `rc internal projects list`, etc.

## Postman
Import `postman/RevenueCat-API-v2.postman_collection.json`:
- Top folders **Developer API v2** → `apiKey`
- **Internal — dashboard session** → **Auth → Login** first, then requests use `rc_auth_token`

## Deep reference
Read **`AGENTS.md`** in the repo root for the full matrix and edge cases (`v1` vs `internal/v1` for `/developers/me`).
