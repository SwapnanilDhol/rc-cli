---
name: revenuecat-cli
description: >-
  Work with the rc-cli repo or embed RevenueCat behavior: choose Public API v2
  (Bearer api key, api.revenuecat.com/v2) vs dashboard session (rc login,
  app.revenuecat.com/internal/v1). Offering metadata: rc offerings update or v2
  rc api POST. Use when the user mentions RevenueCat, rc-cli, API keys, rc login,
  internal API, offerings metadata, or Postman for this project.
---

# RevenueCat CLI — which API am I using?

## Two systems (orthogonal)

| Choose | Auth | Base | Typical need |
|--------|------|------|----------------|
| **V2 public API** | `apiKey` in `~/.revenuerc`, `Bearer` header | `api.revenuecat.com/v2` | Documented REST, automation with secret key |
| **Dashboard session** | `rc login` → cookie `rc_auth_token` | `app.revenuecat.com/internal/v1` (+ some `…/v1` routes) | All projects you can access, PUT/PATCH parity with the web app |

They are **not** interchangeable.

## Commands
- **V2:** `rc api GET '/projects/...'`, `rc config` for `apiKey`; update offering metadata e.g. `rc api POST '/projects/{project_id}/offerings/{id}' -d '{"metadata":{...}}'`
- **Session:** `rc login`, then e.g. `rc offerings list`, **`rc offerings update -o … -m '{"k":"v"}'`** (internal PUT; confirm in **Product catalog → Offerings**)

## Postman
Import `postman/RevenueCat-API-v2.postman_collection.json`:
- **Developer API v2** → `apiKey`
- **Internal — dashboard session** → **Auth → Login**, then **`Update offering (sample body: metadata)`** for a ready-made metadata payload

## Deep reference
Repo root **`AGENTS.md`**: full matrix, `/v1` vs `internal/v1`, offering metadata flow.
