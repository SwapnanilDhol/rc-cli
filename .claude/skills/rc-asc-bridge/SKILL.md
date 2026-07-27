---
name: rc-asc-bridge
description: >-
  Keep a RevenueCat catalog in sync with App Store Connect using rc together with
  asc (asccli.sh). Use when creating a new in-app purchase or subscription end to
  end, when a purchase succeeds but no entitlement is granted, when auditing for
  drift between the two systems, or whenever a task spans both Apple and RevenueCat.
---

# rc ↔ asc bridge

Two CLIs, two systems, one catalog:

| CLI | System | Role |
|---|---|---|
| `asc` ([asccli.sh](https://asccli.sh)) | App Store Connect | **Source of truth.** Defines what users can actually buy. |
| `rc` | RevenueCat | Mirrors Apple's products, then layers entitlements and offerings on top. |

RevenueCat does not create products on Apple's side by itself — it references them.
So the order is always **App Store Connect first, RevenueCat second.**

## The join key

The two systems join on the product identifier:

```
App Store Connect   attributes.productId   "neon_pro_weekly_subscription"
RevenueCat          identifier             "neon_pro_weekly_subscription"
```

Apps join on the bundle ID:

```
RevenueCat  app.platform_identifier   "com.example.MyApp"   (on the type=app_store app)
ASC         attributes.bundleId       "com.example.MyApp"
```

That chain — RevenueCat project → its `app_store` app → bundle ID → ASC app ID — is
enough to go from a project name to an Apple app with no manual lookup.

## The failure this prevents

A product in App Store Connect that is missing from RevenueCat is **purchasable but
not recognised**: the user pays, Apple takes the money, and RevenueCat grants no
entitlement. It produces support tickets, not errors, so nothing surfaces it.

The reverse (in RevenueCat, not in ASC) is usually a harmless leftover — unless it is
attached to an entitlement.

## Audit an app

```bash
./scripts/rc-asc-reconcile.py --project "Neon"
./scripts/rc-asc-reconcile.py --project "Neon" --json
./scripts/rc-asc-reconcile.py --project "Neon" --strict   # also fail on RC-only
```

Exits non-zero when App Store Connect has products RevenueCat is missing, so it drops
straight into CI. It resolves the app automatically from the project name and prints
ready-to-run `rc internal products create` commands for anything missing.

## Create a subscription end to end

Requires `asc auth` and `rc login`.

```bash
# 1. Apple first — the group, then the subscription.
asc subscriptions groups list --app <asc_app_id>
asc subscriptions create \
  --group-id <group_id> \
  --product-id com.example.pro.weekly \
  --reference-name "Pro Weekly" \
  --subscription-period ONE_WEEK

# 2. Price it and localize it — RevenueCat can do neither.
#    --subscription-id takes the productId, the same key rc uses as --identifier.
asc subscriptions pricing summary --app <asc_app_id> --subscription-id com.example.pro.weekly
asc subscriptions pricing equalize --app <asc_app_id> --subscription-id com.example.pro.weekly \
  --base-price 4.99 --dry-run
asc subscriptions versions localizations list --version-id <version_id>

# 3. Register it in RevenueCat, matching --identifier to Apple's productId exactly.
rc internal apps list -p "My Project" --json          # take the type=app_store id
rc internal products create -p "My Project" \
  --app-id <rc_app_id> \
  --product-type subscription \
  --identifier com.example.pro.weekly \
  --name "Pro Weekly"

# 4. Entitlement, then offering.
rc internal products list -p "My Project" --json      # take the prod… id
rc internal entitlements attach-products -p "My Project" \
  -e <entitlement_id> --product-ids <product_id>
rc internal offerings update -p "My Project" -o "Default" --packages '[
  {"identifier":"weekly","display_name":"Weekly",
   "products":[{"product_id":"<product_id>"}]}
]'

# 5. Confirm the two systems agree.
./scripts/rc-asc-reconcile.py --project "My Project"
```

Step 3's `--identifier` must match step 1's `--product-id` character for character.
A typo here is exactly the silent failure described above.

## Which CLI creates the Apple product?

`rc internal apps app-store-products create` can also create a product in App Store
Connect, via RevenueCat's integration.

**Prefer `asc`.** It is talking to Apple directly with your own API key, and it can
set pricing, localizations, review screenshots, and offers — none of which the
RevenueCat path exposes. Use the `rc` path only for a quick throwaway product, and
only when that app's `are_asc_credentials_set_up` is `true` in
`rc internal apps list --json`.

## JSON output shapes

Both CLIs emit JSON, but neither is uniform — check the shape before indexing.

```bash
asc <cmd> --output json     # default is already json; --pretty to indent
rc  <cmd> --json
```

| Command | Shape |
|---|---|
| `asc apps list` | `{"data": [...]}`, fields under `attributes` |
| `asc subscriptions list` | `{"data": [...]}`, fields under `attributes` |
| `rc internal offerings list` | bare array |
| `rc internal products list` | `{"has_more":…, "products":[...]}` |
| `rc internal apps list` | bare array |

Known `rc` bug: a list that comes back **empty** emits `{"status":200}` rather than
`[]`, so `| jq '.[]'` fails on an empty catalog. Guard with
`jq 'if type=="array" then .[] else empty end'` until it is fixed.

## Related

- `rc-projects-workflow` — the RevenueCat side in full
- `rc-offerings-metadata` — paywall metadata and packages
- `asc docs` / `asc search <term>` — asc's own agent-oriented command discovery
- `asc install-skills` — asc's skill pack, for App Store Connect work beyond this bridge
