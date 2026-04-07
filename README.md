# RevenueCat CLI

A comprehensive command-line interface for RevenueCat, providing access to both the public v2 API and powerful internal APIs.

> **Disclaimer**: This is an unofficial CLI project made with love for developers to make their lives easier. It is not affiliated with, endorsed by, or connected to RevenueCat in any way.

## Installation

### Homebrew (recommended for easy upgrades)

```bash
# Tap the repository
brew tap swapnanildhol/homebrew-rc-cli /Users/swapnanildhol/Desktop/cli-projects/homebrew-rc-cli

# Install
brew install rc-cli
```

To upgrade to a new version:
```bash
brew upgrade rc-cli
```

### Local checkout (quick)

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go build -o rc .
./rc internal login
```

### Install globally with Go (`revenuecat-cli` on your `PATH`)

Requires [Go](https://go.dev/dl/) 1.21+ (see `go.mod` for the exact toolchain).

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go install .
```

- The binary is installed as **`revenuecat-cli`** in **`$(go env GOPATH)/bin`** (often `~/go/bin`). Add that directory to your **`PATH`** (and restart the terminal) if `revenuecat-cli` is not found.
- To invoke it as **`rc`**, use either:
  - **Alias** (e.g. in `~/.zshrc`): `alias rc=revenuecat-cli`
  - **Symlink**: `ln -sf "$(go env GOPATH)/bin/revenuecat-cli" /usr/local/bin/rc`
    (use a directory that is already on your `PATH`; on Apple Silicon Homebrew, `/opt/homebrew/bin` is common.)

### Update when a new version is pushed to GitHub

```bash
brew upgrade rc-cli
```

Or with Go:
```bash
cd rc-cli    # your clone directory
git pull
go install .
# restart the shell if you changed PATH; symlink/alias to rc stays valid
```

## Version

Check the current version:
```bash
rc --version
```

## Authentication

**There are two different APIs with different credentials.** Do not mix them. Full matrix: **[AGENTS.md](AGENTS.md)**.

| | **Internal (dashboard)** | **Public API v2** |
|--|--------------------------|---------------------|
| **Use when** | You want **all projects**, dashboard parity, `rc internal login` + `rc internal offerings` / `rc internal projects` / … | You want **documented** `api.revenuecat.com/v2` + **secret API key** |
| **How** | `rc internal login` (email + password) | `rc config` → `apiKey` + `projectId` |
| **In** `~/.revenuerc` | `email`, `password`, `authToken` | `apiKey`, `projectId` |

### Internal API (session — recommended for multi-project)

```bash
rc internal login --email your@email.com --password yourpassword
# Or interactive (prompts for credentials):
rc internal login
```

Credentials are stored in `~/.revenuerc` with automatic token refresh.

### Public v2 API (API key — per project)

```bash
rc config
# Enter your API key and project ID
```

## Quick Commands

**Internal API** (session auth via `rc internal login`):

```bash
# Projects
rc internal projects list              # List all projects
rc internal projects use <id>          # Set default project
rc internal projects create --name "My App"  # Create project

# Entitlements
rc internal entitlements list
rc internal entitlements create --identifier pro --name "Pro Tier"
rc internal entitlements attach-products -e <id> --product-ids <id1>,<id2>

# Offerings
rc internal offerings list
rc internal offerings create --identifier monthly --name "Monthly Offering"
rc internal offerings update -o <offering_id> -m '{"tier":"pro"}'   # metadata + optional -n / -i

# Products
rc internal products list --limit 2500

# Experiments (A/B Testing)
rc internal experiments list
rc internal experiments create --name "Price Test" --offering-a <id> --offering-b <id>

# Charts & Analytics
rc internal charts overview            # Project overview
rc internal charts trials             # Trial analytics
rc internal charts revenue            # Revenue analytics

# Subscriber Lists
rc internal lists list
rc internal lists manifest

# Utilities
rc internal utilities countries
rc internal stores-status
```

**Public v2 API** (API key auth via `rc config`):

```bash
# Projects
rc projects list              # List projects (v2)
rc projects get <id>          # Get project details

# Entitlements
rc entitlements list         # List entitlements (v2)

# Offerings
rc offerings list            # List offerings (v2)
rc packages list             # List packages

# Products
rc products list --limit 2500

# Subscribers
rc subscribers list          # List customers
rc subscribers get <id>      # Get customer details
```

## Full Command Reference

| Category | Commands |
|----------|----------|
| Auth | `rc internal login`, `rc internal logout` |
| Projects | `rc internal projects list|get|use|create` |
| Entitlements | `rc internal entitlements list|create|delete|attach-products|detach-products` |
| Offerings | `rc internal offerings list|get|create|update|delete|duplicate|set-current|archive` |
| Products | `rc internal products list|create|update` |
| Experiments | `rc internal experiments list|get|create|pause|resume|stop|types` |
| Subscriber Lists | `rc internal lists list|get|manifest` |
| Charts | `rc internal charts overview|overview-all|trials|transactions|revenue` |
| Stores | `rc internal stores-status` |
| Utilities | `rc internal utilities countries` |
| Team | `rc internal collaborators list`, `rc internal apikeys list`, `rc internal audit list` |
| Customers (v2) | `rc subscribers list`, `rc subscribers get`, `rc entitlements`, `rc subscriptions` |

## Internal API Endpoints

Base URL: `https://app.revenuecat.com/internal/v1`

- Full CRUD on entitlements, offerings, experiments
- Access to all projects (not just one)
- Team management (collaborators, API keys)
- Audit logs
- Price experiments lifecycle (create, pause, resume, stop)
- Subscriber lists management
- Chart analytics (overview, trials, transactions, revenue)

## Public v2 API Endpoints

Base URL: `https://api.revenuecat.com/v2`

- Read-only access to projects, apps, customers, products, offerings, entitlements
- Requires API key authentication

## Claude Code Integration

This CLI integrates with Claude Code via a skill file. Once Claude Code is configured, you can use natural language to interact with RevenueCat:

```
# Example Claude Code prompts:
"List all my RevenueCat projects"
"Create a new entitlement called Pro in my project"
"Show me the trial analytics for the last 30 days"
"What's the store connection status?"
```

The skill is located at `~/.claude/skills/revenuecat/SKILL.md` and provides:
- Full command documentation
- API endpoint references
- Authentication instructions
- Usage examples

## Notes

- Internal API tokens auto-refresh when expired
- If refresh fails, run `rc login` again
- Customer IDs with special characters are automatically URL-encoded
- Some features require dashboard configuration (webhooks, promotions)

## License

MIT
