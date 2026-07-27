---
name: rc-charts-analytics
description: >-
  Read RevenueCat revenue, trial, and transaction analytics, and manage A/B price
  experiments and subscriber lists, via the dashboard API. Use when asked about
  charts, revenue data, MRR, trial conversion, experiments, or subscriber lists.
---

# Charts, experiments, subscriber lists

Analytics exist **only** on the dashboard API (`rc login`). There is no public v2
charts endpoint and no `rc charts` command — use `rc internal charts`.

## Charts

```bash
rc internal charts overview --json
rc internal charts overview --sandbox --json
rc internal charts overview --app-uuid <app_uuid> --json
rc internal charts overview-all --json                    # across all projects

rc internal charts revenue --start-date 2026-01-01 --end-date 2026-12-31 --json
rc internal charts transactions --start-date 2026-01-01 --json
rc internal charts trials --start-date 2026-01-01 --json
```

Shared flags on `revenue`, `transactions`, and `trials`:
`--start-date`, `--end-date` (both `YYYY-MM-DD`), `--resolution`, `--sandbox`,
`--app-uuid`.

Resolution: `0` daily (default), `1` weekly, `2` monthly.

`overview` and `overview-all` take only `--sandbox` (and `--app-uuid` on `overview`) —
they have no date range.

## Experiments (A/B price tests)

```bash
rc internal experiments list --json
rc internal experiments types
rc internal experiments get -e <experiment_id> --json

rc internal experiments create -n "Price Test" -a <offering_a_id> -b <offering_b_id>
rc internal experiments create -n "Price Test" -a <id> -b <id> \
  --enrollment 50 \
  --type paywall_design \
  --primary-metric "Realized LTV per customer" \
  --notes "Testing new header"

rc internal experiments pause  -e <experiment_id>
rc internal experiments resume -e <experiment_id>
rc internal experiments stop   -e <experiment_id>
```

`-a` and `-b` require offering **IDs** (`ofrng…`), not names — get them from
`rc internal offerings list --json`.

Defaults: `--enrollment 100`, `--type paywall_design`,
`--primary-metric "Realized LTV per customer"`.

Experiment types (`rc internal experiments types` prints these):
`introductory_offer`, `free_trial_offer`, `paywall_design`, `price_point`,
`subscription_duration`, `subscription_ordering`, `other`.

`stop` is terminal — use `pause` if the experiment may resume.

## Subscriber lists

```bash
rc internal lists list --json
rc internal lists list --limit 50 --json
rc internal lists get -l <list_id> --json
rc internal lists manifest --json
```

## Anything not covered

Charts endpoints change more often than the CLI. If a chart you need has no command,
call it directly:

```bash
rc internal api GET '/developers/me/charts_v2/overview'
```
