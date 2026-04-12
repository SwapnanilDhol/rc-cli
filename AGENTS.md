# RevenueCat CLI (`rc`)

CLI for the **RevenueCat Public V2 API** at `https://api.revenuecat.com/v2`.

For the internal dashboard API, use the separate **`rc-internal`** CLI: https://github.com/SwapnanilDhol/rc-internal

## Agent Skills

Agent skills for this CLI are maintained in a separate repository: https://github.com/SwapnanilDhol/rc-cli-skills

## Authentication

```bash
rc config
# Enter your API key and project ID
```

Credentials stored in `~/.revenuerc`.

## Commands

### Subscribers
```bash
rc subscribers list                          # List subscribers
rc subscribers get <customer_id>             # Get subscriber details
rc subscribers entitlements <customer_id>    # Get subscriber entitlements
rc subscribers subscriptions <customer_id>    # Get subscriber subscriptions
```

### Products
```bash
rc products list                             # List products
```

### Offerings
```bash
rc offerings list                            # List offerings
```

### Apps
```bash
rc apps list                                # List apps
```

### Entitlements
```bash
rc entitlements list                         # List entitlements
```

### Charts
```bash
rc charts revenue                            # Revenue chart
rc charts overview                           # Overview chart
```

### Configuration
```bash
rc config show                               # Show current config
rc config unset                              # Clear config
```

### Utilities
```bash
rc utilities countries                       # List supported countries
```

### Direct API Access
```bash
rc api GET '/projects/{project_id}/offerings'
rc api POST '/projects/{project_id}/offerings' -d '{"identifier":"...","name":"..."}'
```

## Two APIs

| CLI | Auth | Base URL | Use for |
|-----|------|----------|---------|
| `rc` | API key (`sk_…`) | `api.revenuecat.com/v2` | Public V2 API, CI/CD |
| `rc-internal` | Session cookie | `app.revenuecat.com/internal/v1` | Dashboard parity |

## Notes

- All commands use Bearer token auth with the API key from `rc config`
- The V2 API is documented at https://www.revenuecat.com/docs/api-v2
- For internal dashboard commands (session-based auth), use the separate **`rc-internal`** CLI