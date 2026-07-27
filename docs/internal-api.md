# rc internal CLI - Internal Dashboard API Commands

The `rc internal` subcommand group provides access to RevenueCat's **Internal Dashboard API** at `https://app.revenuecat.com/internal/v1`.

## Overview

This CLI is separate from the public `/v2` REST API. The confusion arises because:

| CLI Command Group | API Base URL | Authentication | Purpose |
|-------------------|--------------|----------------|---------|
| `rc` (public) | `https://api.revenuecat.com/v2` | API Key (`rc config`) | Customer-facing public API |
| `rc internal` | `https://app.revenuecat.com/internal/v1` | Session Cookie (`rc login`) | Dashboard operations |

### Key Differences

- **Public API (`rc`)**: Use for subscriber management, analytics, entitlements, webhooks, products, offerings - the documented v2 REST API. Requires API key.
- **Internal API (`rc internal`)**: Internal dashboard operations not available in the public API. Requires session login (email/password).

## Commands

### Projects

```
rc internal projects list
rc internal projects get --project-id <id>
rc internal projects use --project-id <id>
rc internal projects create --name <name> --category <cat> --platforms <plat>
```

**API Endpoints:**
- `GET /developers/me/projects`
- `GET /developers/me/projects/{project_id}`
- `POST /developers/me/projects`

### Entitlements

```
rc internal entitlements list
rc internal entitlements create --identifier <id> --name <name>
rc internal entitlements delete --entitlement-id <id>
rc internal entitlements archive --entitlement-id <id>
rc internal entitlements attach-products --entitlement-id <id> --product-ids <ids>
rc internal entitlements detach-products --entitlement-id <id> --product-ids <ids>
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/entitlements`
- `POST /developers/me/projects/{project_id}/entitlements`
- `DELETE /developers/me/projects/{project_id}/entitlements/{entitlement_id}`
- `POST /developers/me/projects/{project_id}/entitlements/{entitlement_id}/actions/archive`
- `POST /developers/me/projects/{project_id}/entitlements/{entitlement_id}/attach_products`
- `POST /developers/me/projects/{project_id}/entitlements/{entitlement_id}/detach_products`

### Offerings

```
rc internal offerings list [--platform <ios|android|macos>]
rc internal offerings get --offering-id <id>
rc internal offerings create --identifier <id> --name <name>
rc internal offerings update --offering-id <id> [--name <name>] [--identifier <id>] [--metadata <json>]
rc internal offerings delete --offering-id <id>
rc internal offerings duplicate --offering-id <id> --identifier <new-id> --name <new-name>
rc internal offerings set-current --offering-id <id>
rc internal offerings archive --offering-id <id>
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/offerings`
- `GET /developers/me/projects/{project_id}/offerings/{offering_id}`
- `POST /developers/me/projects/{project_id}/offerings`
- `PATCH /developers/me/projects/{project_id}/offerings/{offering_id}`
- `DELETE /developers/me/projects/{project_id}/offerings/{offering_id}`
- `POST /developers/me/projects/{project_id}/offerings/{offering_id}/duplicate`
- `POST /developers/me/projects/{project_id}/offerings/{offering_id}/actions/archive`

### Products

```
rc internal products list [--limit <n>]
rc internal products create --app-id <id> --product-type <type> --identifier <id> --name <name>
rc internal products update --product-id <id> [--product-type <type>] [--identifier <id>] [--name <name>]
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/products`
- `POST /developers/me/projects/{project_id}/apps/{app_id}/products`
- `PATCH /developers/me/projects/{project_id}/products/{product_id}`

### Apps

```
rc internal apps list
rc internal apps subscription-groups --app-id <id>
rc internal apps app-store-products create --app-id <id> --product-type <type> --identifier <id> --name <name> [--duration <dur>] [--subscription-group-id <id>]
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/apps`
- `GET /developers/me/projects/{project_id}/apps/{app_id}/subscription_groups`
- `POST /developers/me/projects/{project_id}/apps/{app_id}/app_store_products`

### Charts

```
rc internal charts overview [--sandbox] [--app-uuid <uuid>]
rc internal charts overview-all [--sandbox]
rc internal charts trials [--start-date <date>] [--end-date <date>] [--resolution <0|1|2>] [--sandbox] [--app-uuid <uuid>]
rc internal charts transactions [--start-date <date>] [--end-date <date>] [--resolution <0|1|2>] [--sandbox] [--app-uuid <uuid>]
rc internal charts revenue [--start-date <date>] [--end-date <date>] [--resolution <0|1|2>] [--sandbox] [--app-uuid <uuid>]
```

**API Endpoints:**
- `GET /developers/me/charts_v2/overview`
- `GET /developers/me/charts_v2/trials`
- `GET /developers/me/charts_v2/transactions`
- `GET /developers/me/charts_v2/revenue`

### Subscriber Lists

```
rc internal lists list [--limit <n>]
rc internal lists get --list-id <id>
rc internal lists manifest
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/subscriber_lists`
- `GET /developers/me/subscriber_lists/{list_id}`
- `GET /developers/me/subscriber_lists/manifest`

### Price Experiments

```
rc internal experiments list
rc internal experiments get --experiment-id <id>
rc internal experiments create --name <name> --offering-a <id> --offering-b <id> [--enrollment <pct>] [--type <type>]
rc internal experiments pause --experiment-id <id>
rc internal experiments resume --experiment-id <id>
rc internal experiments stop --experiment-id <id>
rc internal experiments types
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/price_experiments`
- `GET /developers/me/projects/{project_id}/price_experiments/{experiment_id}`
- `POST /developers/me/projects/{project_id}/price_experiments`
- `POST /developers/me/projects/{project_id}/price_experiments/{experiment_id}/pause`
- `POST /developers/me/projects/{project_id}/price_experiments/{experiment_id}/resume`
- `POST /developers/me/projects/{project_id}/price_experiments/{experiment_id}/stop`

### Other Commands

```
rc internal stores-status
rc internal collaborators list
rc internal apikeys list
rc internal audit list
rc internal utilities countries
```

**API Endpoints:**
- `GET /developers/me/projects/{project_id}/product_stores_statuses`
- `GET /developers/me/projects/{project_id}/collaborators`
- `GET /developers/me/projects/{project_id}/api_keys`
- `GET /developers/me/projects/{project_id}/audit_logs`
- `GET /utilities/countries`

## Authentication

Internal API commands require session cookie authentication (not API key):

```bash
# Login with email/password
rc login

# Then you can use internal commands
rc internal projects list
```

The session token is stored in your config (`~/.revenuerc`) and refreshed automatically.