# RevenueCat CLI

An unofficial command-line interface for RevenueCat, covering **both** the
[Public V2 API](https://www.revenuecat.com/docs/api-v2) and the dashboard API —
in one binary.

| Command form | API | Credential |
|---|---|---|
| `rc <cmd>` | `api.revenuecat.com/v2` | API key (`sk_…`) via `rc config` |
| `rc internal <cmd>` | `app.revenuecat.com/internal/v1` | Email + password via `rc login` |

Reading catalog data works on either. Writing — offering metadata, packages,
product creation, entitlement changes, experiments, charts — is dashboard-only,
so it lives under `rc internal`.

> This project is not affiliated with or endorsed by RevenueCat. The dashboard API
> is undocumented and can change without notice. Use `rc login` only with your own
> account.

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+

## Installation

```bash
brew tap swapnanildhol/tap
brew install swapnanildhol/tap/rc-cli
```

Or from source:

```bash
git clone https://github.com/SwapnanilDhol/rc-cli.git
cd rc-cli
go build -o rc .
```

`go install .` installs the binary as **`revenuecat-cli`** in `$(go env GOPATH)/bin`.
Alias or symlink it to `rc`:

```bash
alias rc=revenuecat-cli
# or
ln -sf "$(go env GOPATH)/bin/revenuecat-cli" /usr/local/bin/rc
```

## Authentication

```bash
rc config     # public v2 API key + default project
rc login      # dashboard session, for `rc internal` commands
```

Both are stored in `~/.revenuerc`. Sessions refresh automatically on expiry.

Non-interactive:

```bash
rc login --email you@example.com --password "$RC_PASSWORD"
RC_EMAIL=… RC_PASSWORD=… rc login
```

## JSON output

Every command takes `--json`, which puts a single JSON document on stdout and
progress messages on stderr. Errors exit non-zero.

```bash
rc internal offerings list --json | jq '.[].identifier'
```

## Public v2 commands

```bash
rc projects list                             rc apps list
rc products list                             rc entitlements list
rc offerings list                            rc subscriptions list
rc subscribers list
rc subscribers get <customer_id>
rc subscribers entitlements <customer_id>
rc subscribers subscriptions <customer_id>

rc config                                    # configure key + project
rc config show | rc config unset

rc api GET '/projects/{project_id}/customers' -q limit=20    # raw passthrough
```

## Dashboard commands

```bash
rc internal projects list | get -i <id> | use -i <id> | create -n "Name"
rc internal offerings list | get -o <ref> | create -i <id> -n "Name"
rc internal offerings update -o <ref> --metadata-merge '{"title":"Pro"}'
rc internal offerings duplicate | set-current | archive | delete -o <ref>
rc internal entitlements list | create | archive | delete
rc internal entitlements attach-products -e <id> --product-ids <a>,<b>
rc internal products list --limit 2500 | create | update
rc internal apps list | subscription-groups -i <app_id>
rc internal charts overview | revenue | transactions | trials   # ⚠ see note below
rc internal experiments list | create | pause | resume | stop
rc internal lists list | get -l <id> | manifest
rc internal collaborators list | apikeys list | audit list
rc internal stores-status | utilities countries

rc internal api GET '/developers/me/projects/{project_id}/offerings'
```

Run `rc internal <group> --help` for the full flag set on any group.

> **`rc internal charts` currently returns HTTP 404.** The dashboard API is
> undocumented and its paths move: `charts_v2` no longer resolves, and experiments
> had to be repointed from `price_experiments` to `experiments`. If a command
> 404s, capture the real request from the browser's network tab and update the
> path — the CLI is not at fault.

### Names instead of IDs

`-p/--project-id` and `-o/--offering-id` accept a human name, identifier, or ID:

```bash
rc internal offerings get -p "Habits" -o "Pro Paywall" --json
```

A name matching more than one object errors with the candidates listed.

### Offering metadata: replace vs merge

`--metadata` **replaces the entire metadata object**. To change one key and keep
the rest, use `--metadata-merge`:

```bash
rc internal offerings update -o <ref> --metadata-merge '{"title":"Unlock Pro"}'
rc internal offerings update -o <ref> --metadata-merge '{"headerImageURL":null}'  # delete a key
```

## App Store Connect

RevenueCat mirrors Apple's products — it doesn't create them. So the order is
always **App Store Connect first, RevenueCat second**, using
[`asc`](https://asccli.sh) for the Apple side:

```bash
brew install asc && asc auth   # if you don't have it
```

The two systems join on the product identifier (ASC `productId` == RevenueCat
`identifier`) and on the bundle ID (ASC `bundleId` == RevenueCat
`platform_identifier`).

A product that exists in App Store Connect but not in RevenueCat is
**purchasable but unrecognised** — the user pays and gets no entitlement. It
shows up as support tickets, never as an error. Audit for it with:

```bash
./scripts/rc-asc-reconcile.py --project "My App"
./scripts/rc-asc-reconcile.py --project "My App" --json
```

It resolves the Apple app from the project name, diffs both catalogs, prints
ready-to-run fix commands, and exits non-zero when RevenueCat is missing
something Apple sells — so it drops straight into CI.

See the `rc-asc-bridge` skill for the full end-to-end create flow.

## Agent skills

[`.claude/skills/`](.claude/skills/) holds skills that make common tasks
deterministic for coding agents:

| Skill | Covers |
|---|---|
| `rc-cli-usage` | Start here — two APIs, `--json`, project selection |
| `rc-offerings-metadata` | Offering metadata, paywall config, packages |
| `rc-projects-workflow` | Projects, entitlements, products, apps |
| `rc-charts-analytics` | Charts, experiments, subscriber lists |
| `rc-asc-bridge` | Keeping RevenueCat in sync with App Store Connect |

They are picked up automatically when an agent works **inside this repo**. To use
them in *any* project, install them globally:

```bash
rc install-skills            # → ~/.claude/skills/
rc install-skills --list     # show what's bundled
rc install-skills --force    # overwrite existing copies
rc install-skills --dir DIR  # somewhere else
```

The skills are embedded in the binary, so this works from a Homebrew install with
no checkout. Re-run it after upgrading `rc` to pick up skill changes.

## Documentation

[`API.md`](API.md) is the full endpoint reference for both APIs.

## Updating

```bash
brew upgrade swapnanildhol/tap/rc-cli
git pull && go install .
```

## License

MIT
