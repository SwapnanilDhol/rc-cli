---
name: rc-projects-workflow
description: >-
  Manage RevenueCat projects, entitlements, products, and apps via the dashboard
  API. Use when asked about project setup, creating or attaching entitlements,
  creating products, App Store Connect products, collaborators, API keys, audit
  logs, or store connection status.
---

# Projects, entitlements, products, apps

All of these are dashboard API commands (`rc login` required). Add `--json` for
parseable output.

## Entity model

```
Project
├── Apps            (iOS, Android, Amazon, Stripe …)  ← products belong here
├── Products        (each product belongs to exactly one app)
├── Entitlements    (project-scoped; products attach to them)
└── Offerings       (project-scoped; packages reference products)
```

Project-scoped: entitlements, offerings. App-scoped: products only.
`-p` accepts a project name or `proj…` ID.

## Projects

```bash
rc internal projects list --json
rc internal projects get -i <project_id> --json
rc internal projects use -i <project_id>            # persist as default
rc internal projects create -n "My App"
rc internal projects create -n "My App" --platforms apple_native,google_play
rc internal projects create -n "My App" --category finance
```

`create` sets the new project as the default automatically.
`--platforms` accepts: `apple_native`, `google_play`, `stripe`, `amazon`, and others.

## Entitlements

```bash
rc internal entitlements list --json
rc internal entitlements create -i pro -n "Pro Tier"
rc internal entitlements attach-products -e <entitlement_id> --product-ids <id1>,<id2>
rc internal entitlements detach-products -e <entitlement_id> --product-ids <id1>
rc internal entitlements archive -e <entitlement_id>
rc internal entitlements delete -e <entitlement_id>
```

`--product-ids` is comma-separated and takes product IDs (`prod…`), which come from
`rc internal products list --json`.

## Products

```bash
rc internal products list --json
rc internal products list --limit 2500 --json       # default is 100

rc internal products create \
  --app-id <app_id> \
  --product-type subscription \
  --identifier my.weekly \
  --name "Weekly Subscription"

rc internal products update --product-id <id> -n "New Name"
rc internal products update --product-id <id> -i new.identifier
```

Product types: `subscription`, `non_consumable_product`, `consumable_product`,
`non_renewing_subscription`.

The default `--limit 100` silently truncates larger catalogs — raise it before
concluding a product does not exist.

## Apps

```bash
rc internal apps list --json
rc internal apps subscription-groups -i <app_id> --json
rc internal apps app-store-products create \
  -i <app_id> \
  --subscription-group-id <group_id> \
  --identifier my.weekly \
  --name "Weekly" \
  --duration ONE_WEEK
```

`app-store-products create` also accepts `--subscription-group-name` (to create the
group inline instead of passing an existing ID), `--product-type`, and
`--in-app-purchase-type`.

## Admin and status

These take no flags beyond the globals:

```bash
rc internal collaborators list --json
rc internal apikeys list --json
rc internal audit list --json
rc internal stores-status --json
rc internal utilities countries --json
```

## Setting up a project from scratch

```bash
rc login
rc internal projects create -n "My App" --platforms apple_native
# create sets the new project as default

rc internal apps list --json                        # note the app_id
rc internal products create --app-id <app_id> --product-type subscription \
    --identifier my.weekly --name "Weekly"
rc internal products list --json                    # note the prod… id

rc internal entitlements create -i pro -n "Pro Tier"
rc internal entitlements list --json                # note the entitlement id
rc internal entitlements attach-products -e <entitlement_id> --product-ids <product_id>

rc internal offerings create -i default -n "Default"
# then attach packages and metadata — see the rc-offerings-metadata skill
```

Each step needs an ID produced by the previous one, so run them in order and read
the `--json` output rather than guessing identifiers.
