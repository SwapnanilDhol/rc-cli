# RevenueCat CLI

A comprehensive command-line interface for RevenueCat, providing access to both the public v2 API and powerful internal APIs.

## Installation

```bash
git clone https://github.com/YOUR_USERNAME/rc-cli.git
cd rc-cli
go build -o rc .
./rc login
```

## Authentication

### Internal API (Recommended)

The internal API provides full CRUD operations and access to all projects:

```bash
./rc login
# Enter your RevenueCat email and password
```

Credentials are stored in `~/.revenuerc` with automatic token refresh.

### Public v2 API

For read-only access using an API key:

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

## Notes

- Internal API tokens auto-refresh when expired
- If refresh fails, run `rc login` again
- Customer IDs with special characters are automatically URL-encoded
- Some features require dashboard configuration (webhooks, promotions)

## License

MIT
