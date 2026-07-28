---
name: rc-cli-usage
description: >-
  Entry point for the rc CLI (RevenueCat). Covers the two APIs and two
  credentials, which one a given capability needs, --json output, project
  selection, and global flags. Use when asked about rc, rc-cli, RevenueCat
  auth, API keys, dashboard login, or when deciding which rc command to run.
---

# rc CLI

One binary, two APIs. Which one a command uses is always readable from the command.

| Command form | API | Base URL | Credential | Set up with |
|---|---|---|---|---|
| `rc <cmd>` | Public v2 | `api.revenuecat.com/v2` | API key (`sk_…`) | `rc config` |
| `rc internal <cmd>` | Dashboard | `app.revenuecat.com/internal/v1` | Session cookie | `rc login` |

Both credentials live in `~/.revenuerc`. They are independent — having one does not
give you the other.

## Rule for choosing

**Reading catalog data → either. Writing anything → `rc internal`.**

The public v2 API is read-mostly for catalog objects. Offerings metadata, packages,
experiments, charts, product creation, and entitlement mutation exist **only** under
`rc internal`. If a task involves changing an offering, a package, or a product, go
straight to `rc internal` — do not try the v2 command first.

## Always use --json for programmatic work

Every command takes a global `--json` flag. Without it, output is decorated text with
emoji and ANSI colour, intended for humans and expensive to parse. With it, stdout is
a single JSON document and progress messages go to stderr.

```bash
rc internal offerings list --json
rc internal projects list --json
```

Exit code is non-zero on any HTTP error, so `set -e` and `||` behave. Output is the
API's response verbatim, so a command's JSON type is stable: a list command returns
`[]` when empty, never a "none found" sentence and never an object.

## Setup

```bash
rc login                       # dashboard session — prompts for email + password
rc config                      # public v2 API key + default project

# Non-interactive login:
rc login --email you@example.com --password "$RC_PASSWORD"
RC_EMAIL=… RC_PASSWORD=… rc login
```

## Selecting a project

Almost every `rc internal` command is project-scoped. Set a default once:

```bash
rc internal projects list --json                 # find the ID
rc internal projects use -i proj1a2b3c4d         # persist as default
```

Or override per-command with `-p`. **`-p` accepts a project name as well as an ID** —
`rc internal offerings list -p "Habits"` works. An ambiguous name is an error, never a
silent first-match.

## Global flags

| Flag | Meaning |
|---|---|
| `--json` | Machine-readable output on stdout |
| `-p, --project-id` | Project ID **or name**, overriding the saved default |
| `--api-key` | Public v2 API key, overriding `rc config` |

## Command groups

```bash
# Public v2 (rc config)
rc projects list                 rc subscribers list
rc products list                 rc subscribers get <customer_id>
rc offerings list                rc entitlements list
rc apps list                     rc subscriptions list
rc api GET '/projects/{project_id}/customers'      # raw v2 passthrough

# Dashboard (rc login) — see the other skills for each of these
rc internal projects       rc internal offerings      rc internal products
rc internal entitlements   rc internal apps           rc internal charts
rc internal experiments    rc internal lists          rc internal audit
rc internal collaborators  rc internal apikeys        rc internal stores-status
rc internal utilities countries
```

## Escape hatch

Anything without a typed command is reachable directly, with `{project_id}`
substituted from the selected project:

```bash
rc internal api GET '/developers/me/projects/{project_id}/offerings'
rc internal api PATCH '/developers/me/projects/{project_id}/offerings_with_packages/ofrngXXXX' \
  -d '{"metadata":{"title":"Hello"}}'
rc api GET '/projects/{project_id}/customers' -q limit=20
```

## Related skills

- `rc-offerings-metadata` — offering metadata, packages, paywall config
- `rc-projects-workflow` — projects, entitlements, products, apps
- `rc-charts-analytics` — charts, experiments, subscriber lists
- `rc-asc-bridge` — App Store Connect sync via `asc`; **read this before creating
  any product**, since Apple must have it first

## Installing these skills globally

They load automatically inside this repo. For every other project:

```bash
rc install-skills
```
