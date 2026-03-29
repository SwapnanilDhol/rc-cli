# RevenueCat API Reference

This document covers all APIs accessible by the RevenueCat CLI.

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
| GET | `/projects/{id}/offerings` | List offerings |
| GET | `/projects/{id}/offerings/{id}` | Get offering details |
| GET | `/projects/{id}/offerings/{id}/packages` | Get packages in offering |
| GET | `/projects/{id}/products` | List products |
| GET | `/projects/{id}/products/{id}` | Get product details |
| GET | `/projects/{id}/packages/{id}` | Get package details |
| GET | `/projects/{id}/packages/{id}/products` | Get products in package |
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
| GET | `/developers/me/projects/{project_id}/offerings/{id}` | Get offering |
| POST | `/developers/me/projects/{project_id}/offerings` | Create offering |
| POST | `/developers/me/projects/{project_id}/offerings/{id}/duplicate` | Duplicate offering |
| PUT | `/developers/me/projects/{project_id}/offerings/{id}` | Update offering |
| PATCH | `/developers/me/projects/{project_id}/offerings/{id}` with `{"is_current": true}` | Set offering as current/default |
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
rc packages list          → GET /projects/{id}/offerings/{id}/packages
rc packages get           → GET /projects/{id}/packages/{id}
rc packages products     → GET /projects/{id}/packages/{id}/products
rc entitlements list     → GET /projects/{id}/entitlements
rc entitlements get      → GET /projects/{id}/entitlements/{id}
rc entitlements products  → GET /projects/{id}/entitlements/{id}/products
rc webhooks list         → GET /projects/{id}/webhooks
rc webhooks events       → GET /projects/{id}/webhooks/events
```

### Internal API Commands
```
rc login                  → POST /v1/developers/login
rc logout                 → Clear config
rc internal projects list → GET /developers/me/projects
rc internal projects get  → GET /developers/me/projects/{id}
rc internal entitlements list    → GET /developers/me/projects/{id}/entitlements
rc internal entitlements create  → POST /developers/me/projects/{id}/entitlements
rc internal entitlements delete  → DELETE /developers/me/projects/{id}/entitlements/{id}
rc internal offerings list       → GET /developers/me/projects/{id}/offerings
rc internal offerings create     → POST /developers/me/projects/{id}/offerings
rc internal offerings delete    → DELETE /developers/me/projects/{id}/offerings/{id}
rc internal products list        → GET /developers/me/projects/{id}/products
rc internal apps list            → GET /developers/me/projects/{id}/apps
rc internal collaborators list   → GET /developers/me/projects/{id}/collaborators
rc internal apikeys list         → GET /developers/me/projects/{id}/api_keys
rc internal audit list           → GET /developers/me/projects/{id}/audit_logs
```

---

## Notes

- **Internal API advantages**: Full CRUD operations, access to all projects, team management
- **Public API limitations**: Read-only for most operations, requires per-project API keys
- **Token expiry**: Internal API tokens expire; CLI auto-refreshes when possible
- **API key types**: Secret keys (full access) vs publishable keys (limited)
