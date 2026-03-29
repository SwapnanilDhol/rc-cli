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
| **CLI surface** | `rc api …`, `rc projects …`, `rc offerings …`, etc. (no `internal` in path) | **`rc internal …`** and anything using `internal.EnsureAuthenticated` |
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

## When another tool embeds this CLI

Pass through **which backend** the user wants:

- **“Use my API key / v2 / public API”** → configure `apiKey` + `projectId`, use `rc api` or v2 subcommands.  
- **“Use my account / dashboard / all projects / internal”** → `rc login`, use `rc internal …`.

Never assume one credential satisfies both.
