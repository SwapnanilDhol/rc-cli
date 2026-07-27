# RevenueCat CLI (`rc`)

Unofficial CLI covering **both** RevenueCat APIs in one binary.

| Command form | API | Base URL | Credential |
|---|---|---|---|
| `rc <cmd>` | Public v2 | `https://api.revenuecat.com/v2` | API key (`sk_…`) — `rc config` |
| `rc internal <cmd>` | Dashboard | `https://app.revenuecat.com/internal/v1` | Session cookie — `rc login` |

Both credentials live in `~/.revenuerc`.

## Agent skills

Skills live in [`.claude/skills/`](.claude/skills/) and are discovered automatically.
Start with `rc-cli-usage`.

## Rule for choosing an API

Reading catalog data works on either. **Every write goes through `rc internal`.**
Offerings metadata, packages, product creation, entitlement attach/detach,
experiments, and all charts exist only on the dashboard API.

## Agent conventions

- **Always pass `--json`.** Default output is emoji + ANSI text. `--json` puts one
  JSON document on stdout and progress on stderr. Errors are non-zero exits.
- **`-p` and `-o` accept names, not just IDs.** `rc internal offerings get -p "Habits"
  -o "Pro Paywall"`. Ambiguity is an error listing candidates, never a silent pick.
- **`--metadata` replaces the whole metadata object.** Use `--metadata-merge` to
  change individual keys without losing the rest.
- Offerings are **project-scoped, not app-scoped**. Only products belong to apps.

## Authentication

```bash
rc login                    # dashboard session (rc internal …)
rc config                   # public v2 API key + default project

rc login --email you@example.com --password "$RC_PASSWORD"
RC_EMAIL=… RC_PASSWORD=… rc login
```

Sessions refresh transparently on a 401; `rc login` is only needed once.

## Layout

```
api/              public v2 client (Bearer)
internal/         dashboard client (session cookie, 401 refresh + retry)
cmd/              v2 commands, rc login / rc logout, rc api
cmd/internalapi/  everything under `rc internal`, incl. rc internal api
.claude/skills/   agent skills
API.md            full endpoint reference for both APIs
```

## Escape hatches

```bash
rc api GET '/projects/{project_id}/customers' -q limit=20
rc internal api GET '/developers/me/projects/{project_id}/offerings'
```

`{project_id}` is substituted from the selected project.

## Editing this repo

The skills document a real command surface and go stale silently. After changing
flags or commands, verify against the binary rather than memory:

```bash
go build -o rc . && ./rc internal offerings update --help
```
