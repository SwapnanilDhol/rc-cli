# RevenueCat CLI

A comprehensive command-line interface for RevenueCat, providing access to both the public v2 API and powerful internal APIs.

> **Disclaimer**: This is an unofficial CLI project made with love for developers to make their lives easier. It is not affiliated with, endorsed by, or connected to RevenueCat in any way.

## Installation

### Local checkout (quick)

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go build -o rc .
./rc login
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
cd rc-cli    # your clone directory
git pull
go install .
# restart the shell if you changed PATH; symlink/alias to rc stays valid
```

There is no separate release binary in this repo; **`go install .`** from an up-to-date clone is the supported global install path. (Publishing a Go module under `github.com/...` would allow `go install github.com/.../rc-cli@latest`; the current `go.mod` module path is `revenuecat-cli`, so remote `@latest` install is not configured.)

## Authentication

**There are two different APIs with different credentials.** Do not mix them. Full matrix: **[AGENTS.md](AGENTS.md)**.

| | **Internal (dashboard)** | **Public API v2** |
|--|--------------------------|---------------------|
| **Use when** | You want **all projects**, dashboard parity, `rc login` + `rc offerings` / `rc projects` / … | You want **documented** `api.revenuecat.com/v2` + **secret API key** |
| **How** | `./rc login` (email + password) | `./rc config` → `apiKey` + `projectId` |
| **In** `~/.revenuerc` | `email`, `password`, `authToken` | `apiKey`, `projectId` |

### Internal API (session — recommended for multi-project)

```bash
./rc login
# Enter your RevenueCat email and password
```

Credentials are stored in `~/.revenuerc` with automatic token refresh.

### Public v2 API (API key — per project)

```bash
./rc config
# Enter your API key and project ID
```

## Quick Commands

```bash
# Projects
./rc projects list              # List all projects
./rc projects use <id>          # Set default project
./rc projects create --name "My App"  # Create project

# Entitlements
./rc entitlements list
./rc entitlements create --identifier pro --name "Pro Tier"

# Offerings
./rc offerings list
./rc offerings create --identifier monthly --name "Monthly Offering"
./rc offerings update -o <offering_id> -m '{"tier":"pro"}'   # metadata + optional -n / -i (after rc login)

# Products
./rc products list --limit 2500

# Experiments (A/B Testing)
./rc experiments list
./rc experiments create --name "Price Test" --offering-a <id> --offering-b <id>

# Charts & Analytics
./rc charts overview            # Project overview
./rc charts trials             # Trial analytics
./rc charts revenue            # Revenue analytics

# Subscriber Lists
./rc lists list
./rc lists manifest

# Utilities
./rc utilities countries
./rc stores-status
```

## Multi-Project Support

```bash
# List all accessible projects
./rc projects list

# Set default project (all commands use this project)
./rc projects use <project-id>

# Projects show (current) marker when listed
```

## Full Command Reference

| Category | Commands |
|----------|----------|
| Auth | `login`, `logout` |
| Projects | `list`, `get`, `use`, `create` |
| Entitlements | `list`, `create`, `delete` |
| Offerings | `list`, `get`, `create`, `delete` |
| Products | `list` |
| Experiments | `list`, `get`, `create`, `pause`, `resume`, `stop` |
| Subscriber Lists | `list`, `get`, `manifest` |
| Charts | `overview`, `overview-all`, `trials`, `transactions`, `revenue` |
| Stores | `stores-status` |
| Utilities | `countries` |
| Team | `collaborators list`, `apikeys list`, `audit list` |
| Customers (v2) | `subscribers list`, `subscribers get`, `entitlements`, `subscriptions` |

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
