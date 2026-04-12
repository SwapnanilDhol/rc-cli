# RevenueCat CLI

A command-line interface for the **[RevenueCat Public V2 API](https://www.revenuecat.com/docs/api-v2)**.

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

```bash
rc config
# Enter your API key and project ID
```

Credentials are stored in `~/.revenuerc`.

## Commands

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

## Updating

```bash
brew upgrade swapnanildhol/tap/rc-cli   # Homebrew
git pull && go install .                # From source
```

## License

MIT