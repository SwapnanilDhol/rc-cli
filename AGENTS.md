# Agent / contributor context — RevenueCat CLI (`rc`)

This repo talks to **two different HTTP backends**. Picking the wrong auth is the #1 mistake. Read this before adding commands, docs, or integrations.

---

## The split (memorize this)

| | **Public Developer API v2** | **Internal (dashboard) API** |
|--|------------------------------|--------------------------------|
| **What it is** | Official, documented REST API | Same JSON the web app uses; **not** public docs |
| **Base URL** | `https://api.revenuecat.com/v2` | `https://app.revenuecat.com/internal/v1` (most routes) |
| **Auth** | **`Authorization: Bearer <secret API key>`** (`sk_…` from Project → API keys) | **Session cookie** `rc_auth_token=<token>` (from **`rc login`**) |
| **Config keys** (`~/.revenuerc`) | `apiKey`, `projectId` | `email`, `password`, `authToken` |
| **CLI surface** | `rc api …`, and other commands wired to **`api` package** / Bearer key | After **`rc login`**: dashboard session commands (`internal.EnsureAuthenticated`) e.g. **`rc offerings update`** (PUT offering + **metadata**), `rc projects …`, `rc entitlements …`, `rc experiments …` — see `API.md` / `rc --help` (this repo does not use a single `rc internal …` prefix everywhere) |
| **Project scope** | Secret keys are **per-project**; `GET /v2/projects` lists what **that key** can see | **Account/session** — list **all projects** you can access (`GET /internal/v1/developers/me/projects`) |
| **Postman** | Folder: **Developer API v2** — variable **`apiKey`** | Folder: **Internal — dashboard session** — run **Auth → Login**, variable **`rc_auth_token`** |

There is **no** single “RevenueCat API”: **v2 + API key** and **internal + session** are different products.

---

## Quick decision tree

1. **Need documented REST, server-to-server, CI, or Postman against `api.revenuecat.com`?**  
   → **v2 + API key.** Use `rc config`, then `rc api GET '/projects/{project_id}/…'`.

2. **Need dashboard parity, all projects, CRUD that matches the app, or `rc internal`?**  
   → **Internal + session.** Run `rc login` (email/password). Do **not** put an API key in `Authorization` for `internal/v1`.

3. **“Same” resource, two URLs?**  
   Often yes (e.g. offerings). Field names and verbs may differ. Prefer **v2** for stable automation; **internal** when the CLI already wraps it or you need an undocumented action.

---

## Same-origin `v1` (edge case)

Some routes (e.g. `GET /v1/developers/me`) live at **`https://app.revenuecat.com/v1`**, **not** under `internal/v1`. Wrong host → errors like **7117 Page not found**. See `API.md` (Same-origin `v1`).

---

## Files in this repo

- **`API.md`** — endpoint tables for both backends.  
- **`postman/RevenueCat-API-v2.postman_collection.json`** — v2 + Internal folders; regenerate with `go run ./tools/genpostman`.  
- **`.cursor/rules/revenuecat-cli-auth.mdc`** — Cursor rule: always prefer the correct auth for the base URL.

---

## Offerings: edit + **metadata** (dashboard / internal)

- **CLI (session):** `rc offerings update -o <offering_id> [-n "Display name"] [-i identifier] [-m '{"key":"value"}']` — uses `GET` → strip `packages` → `PUT` on `…/internal/v1/developers/me/projects/{project}/offerings/{id}`.  
- **Verify:** RevenueCat dashboard → **Product catalog → Offerings** → open the offering; refresh if needed. Metadata is part of the offering object (also visible in Paywalls / SDK-facing contexts per RevenueCat).  
- **Postman:** Internal folder → **Update offering (sample body: metadata)**.  
- **Public v2 alternative:** `rc api POST '/projects/{project_id}/offerings/{id}' -d '{"metadata":{...}}'` with **API key** (see [Developer API v2](https://www.revenuecat.com/docs/api-v2)).

---

## App Store Connect products (dashboard / internal)

This repo also supports creating App Store Connect products via the dashboard session.

- **Command:** `rc apps app-store-products create`
- **Uses:** internal/dashboard session cookie (run `rc login`)
- **Endpoint:** `POST https://app.revenuecat.com/internal/v1/developers/me/projects/{project_id}/apps/{app_id}/app_store_products`

For subscription-style products, the CLI validates `--duration` as one of:
`ONE_WEEK`, `ONE_MONTH`, `TWO_MONTHS`, `THREE_MONTHS`, `SIX_MONTHS`, `ONE_YEAR`.

Other fields are passed as provided:
`--product-type`, `--in-app-purchase-type` (passed through), `--identifier`, `--name`,
and for subscriptions: `--subscription-group-id` (+ optional `--subscription-group-name`).

---

## When another tool embeds this CLI

Pass through **which backend** the user wants:

- **“Use my API key / v2 / public API”** → configure `apiKey` + `projectId`, use `rc api` or v2 subcommands.  
- **“Use my account / dashboard / all projects / internal session”** → `rc login`, then session-backed commands (`rc offerings update`, `rc projects list`, … per `API.md`).

Never assume one credential satisfies both.
