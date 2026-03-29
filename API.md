# RevenueCat API Reference

This document covers all APIs accessible by the RevenueCat CLI.

### Two APIs (why it looks duplicated)

RevenueCat exposes **two different HTTP APIs** to programmatic tooling. They talk about the **same dashboard concepts** (projects, offerings, packages, customers) but use **different URLs, auth, and sometimes different JSON field names**. They are not two “versions” of one API; they are two backends.

| | **Public v2** (documented API) | **Internal** (dashboard API) |
|--|-------------------------------|------------------------------|
| **Base** | `https://api.revenuecat.com/v2` | `https://app.revenuecat.com/internal/v1` |
| **Official docs** | Yes — [Developer API v2](https://www.revenuecat.com/docs/api-v2) | No — paths match what the **web app** calls after you log in |
| **Auth** | **API key** in `Authorization: Bearer …` | **Session** from `rc login` (cookie / token in config) |
| **Typical `rc` commands** | `rc offerings …`, `rc packages …`, `rc subscribers …`, `rc projects …` | `rc internal offerings …`, `rc internal projects …`, etc. |

**Rule of thumb:** If the command is under `rc internal`, it hits **internal**. Otherwise most data commands use **public v2** with your project API key. Open **Authentication** below for where each credential lives in `~/.revenuerc`.

---

## Authentication

### Public v2 API (`api.revenuecat.com/v2`)
- **Auth**: Bearer token (API key)
- **Config**: `~/.revenuerc` stores `apiKey`
- **Usage**: `rc config` to set API key

### Internal API (`app.revenuecat.com/internal/v1`)
- **Auth**: Session cookie (`rc_auth_token`)
- **Login**: Email + password via `POST /v1/developers/login`
- **Refresh**: `POST /v1/developers/login/refresh-token` with existing cookie
- **Config**: `~/.revenuerc` stores `email`, `password`, `authToken`
- **Auto-refresh**: CLI automatically refreshes expired tokens

---

## Public v2 API Endpoints

**Base URL**: `https://api.revenuecat.com/v2`

Full schemas and permissions are in the [Developer API v2](https://www.revenuecat.com/docs/api-v2) reference. Below lists the routes most relevant to **offerings** and **packages**; other resources are summarized.

### Offerings & packages (v2)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/{id}/offerings` | List offerings (`starting_after`, `limit`, `expand` = `items.package` or `items.package.product`) |
| POST | `/projects/{id}/offerings` | Create offering (body: `lookup_key`, `display_name`, optional `metadata`) |
| GET | `/projects/{id}/offerings/{id}` | Get offering (`expand` = `package` or `package.product`) |
| POST | `/projects/{id}/offerings/{id}` | Update offering (`display_name`, `is_current`, `metadata`) |
| DELETE | `/projects/{id}/offerings/{id}` | Delete offering and attached packages |
| POST | `/projects/{id}/offerings/{id}/actions/archive` | Archive offering |
| POST | `/projects/{id}/offerings/{id}/actions/unarchive` | Unarchive offering |
| GET | `/projects/{id}/offerings/{id}/packages` | List packages in offering (`starting_after`, `limit`, `expand`) |
| POST | `/projects/{id}/offerings/{id}/packages` | Create package (`lookup_key`, `display_name`, optional `position`) |
| GET | `/projects/{id}/packages/{id}` | Get package (`expand` = `product`) |
| POST | `/projects/{id}/packages/{id}` | Update package (`display_name`, `position`) |
| DELETE | `/projects/{id}/packages/{id}` | Delete package |
| GET | `/projects/{id}/packages/{id}/products` | List products attached to package |
| POST | `/projects/{id}/packages/{id}/actions/attach_products` | Attach products to package |
| POST | `/projects/{id}/packages/{id}/actions/detach_products` | Detach products from package |

**Note:** Public v2 uses **POST** for offering and package updates. The internal dashboard API uses **PUT**/**PATCH** for the same concepts (see internal Offerings table below).

### Other v2 endpoints (summary)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | List all projects |
| GET | `/projects/{id}` | Get project details |
| GET | `/projects/{id}/apps` | List apps in project |
| GET | `/projects/{id}/apps/{id}` | Get app details |
| GET | `/projects/{id}/customers` | List customers (supports `limit`, `starting_after`, `products`, `platform`) |
| GET | `/projects/{id}/customers/{id}` | Get customer details |
| GET | `/projects/{id}/customers/{id}/active_entitlements` | Get customer's active entitlements |
| GET | `/projects/{id}/customers/{id}/subscriptions` | Get customer's subscriptions |
| GET | `/projects/{id}/products` | List products |
| GET | `/projects/{id}/products/{id}` | Get product details |
| GET | `/projects/{id}/entitlements` | List entitlements |
| GET | `/projects/{id}/entitlements/{id}` | Get entitlement details |
| GET | `/projects/{id}/entitlements/{id}/products` | Get products in entitlement |
| GET | `/projects/{id}/webhooks` | List webhooks |
| GET | `/projects/{id}/webhooks/events` | Get webhook events |
| GET | `/projects/{id}/offers` | List promotional offers |

---

## Internal API Endpoints

**Base URL**: `https://app.revenuecat.com/internal/v1`

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/developers/login` | Login with email/password, returns auth token |
| POST | `/v1/developers/login/refresh-token` | Refresh auth token (with existing cookie) |

### Developer/Account
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects` | List all projects for current user |
| POST | `/developers/me/projects` | Create project |

### Entitlements (CRUD)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/entitlements` | List entitlements |
| GET | `/developers/me/projects/{project_id}/entitlements/{id}` | Get entitlement |
| POST | `/developers/me/projects/{project_id}/entitlements` | Create entitlement |
| PUT | `/developers/me/projects/{project_id}/entitlements/{id}` | Update entitlement |
| DELETE | `/developers/me/projects/{project_id}/entitlements/{id}` | Delete entitlement |
| POST | `/developers/me/projects/{project_id}/entitlements/{id}/actions/archive` | Archive entitlement |
| GET | `/developers/me/projects/{project_id}/entitlements/{id}/products` | Get products in entitlement |

### Offerings (CRUD)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/offerings` | List offerings |
| GET | `/developers/me/projects/{project_id}/offerings?platform=IOS` | List offerings filtered by platform |
| GET | `/developers/me/projects/{project_id}/offerings/{id}` | Get offering (response includes nested `packages` when present) |
| POST | `/developers/me/projects/{project_id}/offerings` | Create offering (`identifier`, `display_name` in CLI) |
| POST | `/developers/me/projects/{project_id}/offerings/{id}/duplicate` | Duplicate offering (not exposed in public v2; dashboard + `rc internal offerings duplicate`) |
| PUT | `/developers/me/projects/{project_id}/offerings/{id}` | Update offering (e.g. identifier, display name, metadata; not wrapped in CLI) |
| PATCH | `/developers/me/projects/{project_id}/offerings/{id}` | Partial update; CLI sends `{"is_current": true}` to set default offering |
| POST | `/developers/me/projects/{project_id}/offerings/{id}/actions/archive` | Archive offering |
| DELETE | `/developers/me/projects/{project_id}/offerings/{id}` | Delete offering |

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/products` | List products |
| GET | `/developers/me/projects/{project_id}/products?limit=2500` | List products with limit (max 2500) |
| GET | `/developers/me/projects/{project_id}/product_stores_statuses` | Get store connection status |

### Apps
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/apps` | List apps |

### Collaborators
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/collaborators` | List collaborators (team members) |
| GET | `/developers/me/projects/{project_id}/collaborators/{id}` | Get collaborator |
| POST | `/developers/me/projects/{project_id}/collaborators` | Add collaborator |
| PUT | `/developers/me/projects/{project_id}/collaborators/{id}` | Update collaborator |
| DELETE | `/developers/me/projects/{project_id}/collaborators/{id}` | Remove collaborator |

### API Keys
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/api_keys` | List API keys |
| POST | `/developers/me/projects/{project_id}/api_keys` | Create API key |
| DELETE | `/developers/me/projects/{project_id}/api_keys/{id}` | Delete API key |

### Audit Logs
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/audit_logs` | Get audit logs |

### Webhooks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/webhooks` | List webhooks |
| POST | `/developers/me/projects/{project_id}/webhooks` | Create webhook |
| PUT | `/developers/me/projects/{project_id}/webhooks/{id}` | Update webhook |
| DELETE | `/developers/me/projects/{project_id}/webhooks/{id}` | Delete webhook |
| POST | `/developers/me/projects/{project_id}/webhooks/{id}/test` | Test webhook |
| GET | `/developers/me/projects/{project_id}/webhooks/events` | Get webhook events |

### Price Experiments (A/B Testing)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/price_experiments` | List price experiments |
| GET | `/developers/me/projects/{project_id}/price_experiments/{id}` | Get experiment details |
| POST | `/developers/me/projects/{project_id}/price_experiments` | Create price experiment |
| POST | `/developers/me/projects/{project_id}/price_experiments/{id}/pause` | Pause experiment |
| POST | `/developers/me/projects/{project_id}/price_experiments/{id}/resume` | Resume experiment |
| POST | `/developers/me/projects/{project_id}/price_experiments/{id}/stop` | Stop experiment |

### Subscriber Lists
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/subscriber_lists` | List subscriber lists |
| GET | `/developers/me/projects/{project_id}/subscriber_lists?limit=700` | List with limit |
| GET | `/developers/me/subscriber_lists/{id}` | Get subscriber list details |
| GET | `/developers/me/subscriber_lists/manifest` | Get manifest of all lists |

### Charts V2
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/charts_v2/overview` | Get project overview analytics |
| GET | `/developers/me/charts_v2/overview` (no app_uuid) | Get overview for all projects |
| GET | `/developers/me/charts_v2/trials` | Get trial analytics |
| GET | `/developers/me/charts_v2/transactions` | Get transaction analytics |
| GET | `/developers/me/charts_v2/revenue` | Get revenue analytics |

### Utilities
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/utilities/countries` | List supported countries |

### Promotional Offers
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/promotions` | List promotions |
| GET | `/developers/me/projects/{project_id}/intro_offers` | List introductory offers |

### Paywalls (dashboard)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/developers/me/projects/{project_id}/paywalls` | List paywall configs; query `include_localizations` (`true` on Paywalls page, `false` when bundled with Offerings) |

---

## Product catalog (offerings, entitlements, products, virtual currencies)

Paths below use placeholders only (`{project_id}`, `{offering_id}`, etc.). The dashboard route prefix is `/projects/{project_id}/product-catalog/…` for offerings, products, entitlements, and virtual currencies.

### Offerings — internal API (`/internal/v1`)

Used by **`rc internal offerings …`** (session cookie). The top-level **`rc offerings …`** commands use **public v2** with the API key (`cmd/products.go`).

| Operation | Method | Path | Body / query |
|-----------|--------|------|----------------|
| List | GET | `/developers/me/projects/{project_id}/offerings` | Optional query `platform` (e.g. `IOS`, `ANDROID`) |
| Get | GET | `/developers/me/projects/{project_id}/offerings/{offering_id}` | Returns offering with nested `packages` (same shape as `internal.Offering` in code) |
| Create | POST | `/developers/me/projects/{project_id}/offerings` | JSON: `identifier`, `display_name` (public v2 equivalent field is `lookup_key`) |
| Update | PUT | `/developers/me/projects/{project_id}/offerings/{offering_id}` | Full/partial update as sent by the dashboard (not in CLI) |
| Set current | PATCH | `/developers/me/projects/{project_id}/offerings/{offering_id}` | JSON: `is_current`: `true` (all this subcommand sends today) |
| Duplicate | POST | `/developers/me/projects/{project_id}/offerings/{offering_id}/duplicate` | JSON: `identifier`, `display_name`, `packages_only` (bool) |
| Archive | POST | `/developers/me/projects/{project_id}/offerings/{offering_id}/actions/archive` | Empty JSON object `{}` |
| Delete | DELETE | `/developers/me/projects/{project_id}/offerings/{offering_id}` | — |

**Packages:** Creating, editing, or removing packages from the dashboard ultimately uses **public v2** package endpoints (Bearer key), e.g. `POST …/offerings/{offering_id}/packages`, `POST …/packages/{package_id}`, attach/detach actions. The CLI exposes **read-only** v2 helpers: `rc packages list|get|products` and `rc offerings packages` (see [CLI Commands](#cli-commands-by-api)).

**Dashboard:** opening **Offerings** also triggers `GET …/paywalls?include_localizations=false` (see Paywalls above). Navigating to **Product catalog → Offerings → `{offering_id}`** loads the same offering **GET** as above (no separate “packages-only” internal route required for read; list/detail responses include packages when applicable).

**Duplicate offering** exists on **internal** only; there is no v2 `duplicate` route in the public reference.

### Entitlements — internal API (`/internal/v1`)

| Operation | Method | Path | Body / query |
|-----------|--------|------|----------------|
| List | GET | `/developers/me/projects/{project_id}/entitlements` | — |
| Get | GET | `/developers/me/projects/{project_id}/entitlements/{entitlement_id}` | — (used by dashboard; use public v2 `rc entitlements get` for read with API key) |
| Create | POST | `/developers/me/projects/{project_id}/entitlements` | JSON: `identifier`, `display_name` |
| Update | PUT | `/developers/me/projects/{project_id}/entitlements/{entitlement_id}` | JSON fields per dashboard (not wrapped in CLI yet) |
| Archive | POST | `/developers/me/projects/{project_id}/entitlements/{entitlement_id}/actions/archive` | Empty JSON `{}` |
| Delete | DELETE | `/developers/me/projects/{project_id}/entitlements/{entitlement_id}` | — |

CLI commands: `rc entitlements list|create|archive|delete` (internal auth).

### Products — internal API (`/internal/v1`)

| Operation | Method | Path | Body / query |
|-----------|--------|------|----------------|
| List | GET | `/developers/me/projects/{project_id}/products` | Query `limit` (dashboard/CLI often use large limits) |

The CLI only implements **list** (`rc products list`). Product **creation** in RevenueCat is normally tied to App Store / Play / other stores (import / sync), not a generic “create SKU” wizard in the same shape as offerings. If the dashboard exposes extra mutations, capture them via DevTools when you attach or refresh store products.

### Virtual currencies — **public v2 API** (`https://api.revenuecat.com/v2`)

Configured under **Product catalog → Virtual currencies**. These use a **secret API key** (Bearer), not the dashboard session cookie. See [Virtual Currency](https://docs.revenuecat.com/docs/offerings/virtual-currency) and [Developer API v2](https://www.revenuecat.com/docs/api-v2).

**Project configuration (define currencies):**

| Operation | Method | Path |
|-----------|--------|------|
| List | GET | `/projects/{project_id}/virtual_currencies` |
| Create | POST | `/projects/{project_id}/virtual_currencies` |
| Get | GET | `/projects/{project_id}/virtual_currencies/{virtual_currency_code}` |
| Update | POST | `/projects/{project_id}/virtual_currencies/{virtual_currency_code}` |
| Delete | DELETE | `/projects/{project_id}/virtual_currencies/{virtual_currency_code}` |
| Archive | POST | `/projects/{project_id}/virtual_currencies/{virtual_currency_code}/actions/archive` |
| Unarchive | POST | `/projects/{project_id}/virtual_currencies/{virtual_currency_code}/actions/unarchive` |

**Per-customer balances / ledger** (same v2 base):

| Operation | Method | Path |
|-----------|--------|------|
| List balances | GET | `/projects/{project_id}/customers/{customer_id}/virtual_currencies` |
| Create transaction | POST | `/projects/{project_id}/customers/{customer_id}/virtual_currencies/transactions` |
| Update balance (no transaction row) | POST | `/projects/{project_id}/customers/{customer_id}/virtual_currencies/update_balance` |

**CLI:** virtual currency commands are not in `rc` yet; use v2 with `rc config` API key or a small script.

### Capturing create/delete in the browser

If automation shows an empty Offerings/table or missing XHR, refresh the tab, confirm you are still logged in, or export a **HAR** while you click **New offering**, save, and delete/archive. Offering CRUD often appears as **`/internal/v1/.../offerings`** (session cookie); **package** create/update and product attach/detach are **public v2** (`api.revenuecat.com/v2`, Bearer). Compare captured URLs to the tables above and the v2 **Offerings & packages** section.

---

## Dashboard-observed API paths (session cookie)

Observed from browser network traffic on `app.revenuecat.com` (no project identifiers or credentials recorded here). Used for extending the CLI; paths may change without notice.

### Internal (`https://app.revenuecat.com/internal/v1`)

| Method | Endpoint | Where seen |
|--------|----------|------------|
| GET | `/developers/me/projects` | All-projects overview |
| GET | `/developers/me/collaborations/pending` | All-projects overview |
| GET | `/developers/me/dashboard_notifications` | All-projects overview |
| GET | `/rc_capital/onboarding/status` | All-projects overview (returned 404 for observer) |
| GET | `/developers/me/projects/{project_id}/offerings` | Product catalog → Offerings |
| GET | `/developers/me/projects/{project_id}/offerings/{offering_id}` | Product catalog → offering detail (`…/product-catalog/offerings/{offering_id}`) |
| GET | `/developers/me/projects/{project_id}/paywalls` | Offerings (`include_localizations=false`); Paywalls page (`include_localizations=true`) |

### Same-origin `v1` (`https://app.revenuecat.com/v1`)

The web app also calls JSON APIs under `/v1/developers/me/...` (session cookie), parallel to some `internal/v1` charts routes the CLI uses.

**Note:** These paths are **not** mounted at `https://app.revenuecat.com/internal/v1/developers/me/...`. Calling e.g. `GET …/internal/v1/developers/me` returns application error **7117** (“Page not found”). Use `https://app.revenuecat.com/v1/developers/me` for the rows below; use `…/internal/v1/developers/me/projects` (and the rest of the internal table) for catalog-style JSON the CLI uses with the same cookie.

| Method | Endpoint | Where seen |
|--------|----------|------------|
| GET | `/developers/me` | After login |
| GET | `/developers/me/billing/info` | After login |
| GET | `/developers/me/transactions` | All-projects overview |
| GET | `/developers/me/charts_v2/overview` | All-projects overview (e.g. `sandbox_mode`, `v3`) |
| GET | `/developers/me/charts_v2/customers_new` | All-projects overview (e.g. `start_date`, `end_date`, `resolution`, `is_sparkline`) |
| GET | `/developers/me/charts_v2/revenue` | All-projects overview (same style params) |
| GET | `/developers/me/charts_v2/trials` | All-projects overview |
| GET | `/developers/me/charts_v2/actives` | All-projects overview |
| GET | `/developers/me/charts_v2/mrr` | All-projects overview |

### Dashboard SPA routes (for mapping screens → traffic)

Patterns use `{project_id}` as in the URL bar (format may vary by account).

| Path pattern |
|--------------|
| `/overview` (all projects) |
| `/projects/{project_id}/overview` |
| `/projects/{project_id}/charts` (may redirect to a chart sub-path such as `charts/revenue` with query params) |
| `/projects/{project_id}/customers` |
| `/projects/{project_id}/product-catalog/offerings` |
| `/projects/{project_id}/product-catalog/products` |
| `/projects/{project_id}/product-catalog/entitlements` |
| `/projects/{project_id}/product-catalog/virtual-currencies` |
| `/projects/{project_id}/paywalls` |
| `/projects/{project_id}/ads` |
| `/projects/{project_id}/targeting` |
| `/projects/{project_id}/experiments` |
| `/projects/{project_id}/web-home` |
| `/projects/{project_id}/customer-center` |
| `/projects/{project_id}/integrations` |
| `/projects/{project_id}/project-settings` |

For **Apps & providers** and nested **Project settings** items, open each section and read the URL bar or network tab—the path segment may not match the sidebar label exactly.

Some screens lazy-load data; capture XHR **after** a short wait on each route, or export a HAR for complete coverage.

---

## CLI Commands by API

### Public v2 API Commands
```
rc projects list          → GET /projects
rc projects get           → GET /projects/{id}
rc apps list              → GET /projects/{id}/apps
rc apps get               → GET /projects/{id}/apps/{id}
rc subscribers list       → GET /projects/{id}/customers
rc subscribers get        → GET /projects/{id}/customers/{id}
rc subscribers search     → GET /projects/{id}/customers/{id}
rc subscribers entitlements → GET /projects/{id}/customers/{id}/active_entitlements
rc subscribers subscriptions → GET /projects/{id}/customers/{id}/subscriptions
rc products list          → GET /projects/{id}/products
rc products get           → GET /projects/{id}/products/{id}
rc offerings list         → GET /projects/{id}/offerings
rc offerings get          → GET /projects/{id}/offerings/{id}
rc offerings packages     → GET /projects/{id}/offerings/{id}/packages
rc packages list          → GET /projects/{id}/offerings/{id}/packages
rc packages get           → GET /projects/{id}/packages/{id}
rc packages products     → GET /projects/{id}/packages/{id}/products
rc entitlements list     → GET /projects/{id}/entitlements
rc entitlements get      → GET /projects/{id}/entitlements/{id}
rc entitlements products  → GET /projects/{id}/entitlements/{id}/products
rc webhooks list         → GET /projects/{id}/webhooks (legacy path; v2 integrations use `rc api` below)
rc webhooks events       → GET /projects/{id}/webhooks/events
```

**Full v2 surface (all methods and paths from the official reference):** use the generic HTTP helper (aliases `rc v2`, `rc http`):

```
rc api <METHOD> <path> [-q key=value ...] [-d '<json>' | --data-file file]
```

- Paths are rooted at `https://api.revenuecat.com/v2` (pass only the path, e.g. `/projects/{project_id}/customers`).
- `{project_id}` or `{{project_id}}` is replaced by the configured default project (or `-p`), same as other `rc` commands.
- Other path segments (`customer_id`, `paywall_id`, …) are literal in the path string.

**Postman:** import `postman/RevenueCat-API-v2.postman_collection.json`. It includes **Developer API v2** (variables `baseUrl`, `apiKey`, …) and a top-level folder **Internal — dashboard session** (`internalBaseUrl`, `authBaseUrl`, `rc_email`, `rc_password`, `rc_auth_token`). Run **Auth → Login** in Postman to set `rc_auth_token` from the login JSON. Regenerate from repo root:

```bash
go run ./tools/genpostman
```

Friendly subcommands above remain for common read flows; **writes** and any route not wrapped yet go through `rc api`.

### Internal API Commands (after `rc login`; top-level `rc …`, not `rc internal …`)
```
rc login                  → POST /v1/developers/login
rc logout                 → Clear config
rc projects list          → GET /developers/me/projects
rc projects get           → GET /developers/me/projects/{id}
rc entitlements list      → GET /developers/me/projects/{id}/entitlements
rc entitlements create    → POST /developers/me/projects/{id}/entitlements
rc entitlements delete    → DELETE /developers/me/projects/{id}/entitlements/{id}
rc offerings list         → GET /developers/me/projects/{id}/offerings
rc offerings get          → GET /developers/me/projects/{id}/offerings/{id}
rc offerings create       → POST /developers/me/projects/{id}/offerings
rc offerings update       → PUT /developers/me/projects/{id}/offerings/{id} (display name, identifier, metadata; see AGENTS.md)
rc offerings delete       → DELETE /developers/me/projects/{id}/offerings/{id}
rc offerings duplicate    → POST /developers/me/projects/{id}/offerings/{id}/duplicate
rc offerings set-current  → PATCH /developers/me/projects/{id}/offerings/{id}
rc offerings archive      → POST /developers/me/projects/{id}/offerings/{id}/actions/archive
rc products list          → GET /developers/me/projects/{id}/products
rc apps list              → GET /developers/me/projects/{id}/apps
rc collaborators list     → GET /developers/me/projects/{id}/collaborators
rc apikeys list           → GET /developers/me/projects/{id}/api_keys
rc audit list             → GET /developers/me/projects/{id}/audit_logs
```

---

## Notes

- **Internal API advantages**: Full CRUD operations, access to all projects, team management
- **Public API limitations**: Read-only for most operations, requires per-project API keys
- **Token expiry**: Internal API tokens expire; CLI auto-refreshes when possible
- **API key types**: Secret keys (full access) vs publishable keys (limited)
