---
name: rc-cli
description: >-
  Work with the rc-cli repository for the RevenueCat Public V2 API at
  api.revenuecat.com/v2. Use API key authentication via `rc config`. Covers
  subscribers, products, offerings, apps, entitlements, charts, configuration,
  and utilities. Use when user mentions rc-cli, RevenueCat API, API key auth,
  or public V2 API commands.
---

# RevenueCat CLI — Public V2 API

A CLI for the **RevenueCat Public V2 API** at `https://api.revenuecat.com/v2`.

## Authentication

```bash
rc config
# Enter your API key and project ID
```

Credentials are stored in `~/.revenuerc`.

## Commands

### Subscribers

```bash
rc subscribers list                          # List subscribers
rc subscribers get <customer_id>              # Get subscriber details
rc subscribers entitlements <customer_id>    # Get subscriber entitlements
rc subscribers subscriptions <customer_id>    # Get subscriber subscriptions
```

### Products

```bash
rc products list                             # List products
```

### Offerings

```bash
rc offerings list                             # List offerings
```

### Apps

```bash
rc apps list                                 # List apps
```

### Entitlements

```bash
rc entitlements list                          # List entitlements
```

### Charts

```bash
rc charts revenue                             # Revenue chart
rc charts overview                            # Overview chart
```

### Configuration

```bash
rc config                                    # Configure API key and project ID
rc config show                               # Show current configuration
rc config unset                              # Clear configuration
```

### Utilities

```bash
rc utilities countries                        # List supported countries
```

### API (direct V2 access)

```bash
rc api GET '/projects/{project_id}/offerings'
rc api POST '/projects/{project_id}/offerings' -d '{"identifier":"...","name":"..."}'
```

## Notes

- All commands use Bearer token auth with the API key from `rc config`
- The V2 API is documented at https://www.revenuecat.com/docs/api-v2
- For internal dashboard commands (session-based auth), use the separate `rc-internal` CLI