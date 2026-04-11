# RevenueCat CLI

A command-line interface for the **[RevenueCat Public V2 API](https://www.revenuecat.com/docs/api-v2)** and the **Internal Dashboard API** (`app.revenuecat.com/internal/v1`).

## Two APIs, One CLI

| | **Public V2 API** | **Internal Dashboard API** |
|--|---|---|
| **Base URL** | `api.revenuecat.com/v2` | `app.revenuecat.com/internal/v1` |
| **Auth** | API key (`sk_…`) | Session cookie |
| **Use for** | Documented REST API, CI/CD, server-to-server | Full dashboard parity, all projects, CRUD |
| **Configure** | `rc config` | `rc login` |

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+

## Installation

### Homebrew

```bash
brew tap swapnanildhol/tap
brew install swapnanildhol/tap/rc-cli
```

### From source

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go build -o rc .
```

### Install globally

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go install .
```

The binary is installed as **`revenuecat-cli`** in **`$(go env GOPATH)/bin`**. Add that directory to your `PATH`, then invoke it as `rc` via alias or symlink:

```bash
# Alias (e.g. in ~/.zshrc)
alias rc=revenuecat-cli

# Or symlink
ln -sf "$(go env GOPATH)/bin/revenuecat-cli" /usr/local/bin/rc
```

## Authentication

### Public V2 API

```bash
rc config
# Enter your API key and project ID
```

Credentials are stored in `~/.revenuerc`.

### Internal Dashboard API

```bash
rc login
# Enter your RevenueCat email and password
```

Session token is stored in `~/.revenuerc` with auto-refresh.

## Public V2 Commands

### Subscribers

```bash
rc subscribers list                          # List subscribers
rc subscribers get <customer_id>             # Get subscriber details
rc subscribers entitlements <customer_id>    # Get subscriber entitlements
rc subscribers subscriptions <customer_id>   # Get subscriber subscriptions
```

### Products

```bash
rc products list                              # List products
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
rc entitlements list                         # List entitlements
```

### Charts

```bash
rc charts revenue                            # Revenue chart
rc charts overview                           # Overview chart
```

### Configuration

```bash
rc config                                    # Configure API key and project ID
rc config show                               # Show current configuration
rc config unset                              # Clear configuration
```

### Utilities

```bash
rc utilities countries                       # List supported countries
```

## Internal Dashboard Commands

All internal commands use `rc internal` and require `rc login` first.

### Projects

```bash
rc internal projects list              # List all accessible projects
rc internal projects use <id>          # Set default project
rc internal projects create --name "My App"  # Create project
```

### Entitlements

```bash
rc internal entitlements list
rc internal entitlements create --identifier pro --name "Pro Tier"
rc internal entitlements attach-products --entitlement-id <id> --product-ids <ids>
```

### Offerings

```bash
rc internal offerings list
rc internal offerings create --identifier monthly --name "Monthly Offering"
rc internal offerings update --offering-id <id> --name <name> --metadata '{"tier":"pro"}'
```

### Products

```bash
rc internal products list --limit 2500
rc internal products create --app-id <id> --product-type subscription --identifier <id> --name <name>
```

### Charts & Analytics

```bash
rc internal charts overview            # Project overview
rc internal charts trials              # Trial analytics
rc internal charts revenue             # Revenue analytics
```

### More Internal Commands

```bash
rc internal apps list
rc internal apps subscription-groups --app-id <id>
rc internal apps app-store-products create --app-id <id> ...
rc internal experiments list
rc internal experiments create --name "Price Test" --offering-a <id> --offering-b <id>
rc internal lists list
rc internal stores-status
rc internal collaborators list
rc internal apikeys list
rc internal audit list
rc internal utilities countries
```

## Updating

```bash
brew upgrade swapnanildhol/tap/rc-cli   # Homebrew
git pull && go install .                # From source
```

## License

MIT
