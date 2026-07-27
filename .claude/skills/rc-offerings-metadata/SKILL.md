---
name: rc-offerings-metadata
description: >-
  Read and update RevenueCat offering metadata — paywall title, subtitle,
  features, colour scheme, reviews, package presentation — and attach packages,
  via the dashboard API. Use when asked to change offering metadata, paywall
  configuration, package setup, or anything about a specific offering.
---

# Offering metadata

All offering writes go through the dashboard API: `rc internal offerings …`
(requires `rc login`). The public v2 API cannot write offering metadata.

## Entity model — read this before running anything

**Offerings belong to a project, not to an app.** A project contains apps
(iOS, Android, …) and, separately, one catalog of offerings, products, and
entitlements shared across those apps.

So a request like *"update the offerings metadata on one project in one app"*
resolves to: pick the **project**, then pick the **offering**. There is no app
dimension and no `--app-id` flag on any offering command. Apps only matter for
products (`rc internal products create --app-id …`).

## The one thing that will silently destroy data

`--metadata` **replaces the entire metadata object.** Passing
`--metadata '{"title":"New"}'` deletes `subtitle`, `features`, `colorScheme`,
`reviews`, and everything else.

Use `--metadata-merge` to change some keys and keep the rest. It GETs the current
metadata, applies your keys on top, and writes the result back.

```bash
# WRONG if the offering already has metadata — wipes every other key
rc internal offerings update -o ofrngXXXX --metadata '{"title":"Unlock Pro"}'

# RIGHT — changes title, preserves everything else
rc internal offerings update -o ofrngXXXX --metadata-merge '{"title":"Unlock Pro"}'

# Delete a single key
rc internal offerings update -o ofrngXXXX --metadata-merge '{"headerImageURL":null}'
```

Reach for `--metadata` only when deliberately replacing the whole object, e.g.
configuring a brand-new offering.

## Deterministic procedure

For *"update the offering metadata on project P"*, run exactly this:

```bash
# 1. Point at the project. -p accepts a name or an ID; skip if the default is right.
rc internal projects list --json

# 2. Read what is there now. Never write metadata without reading it first.
rc internal offerings list -p "P" --json

# 3. Write. -o accepts an offering ID, identifier, or display name.
rc internal offerings update -p "P" -o "Pro Paywall" \
  --metadata-merge '{"title":"Unlock Pro","subtitle":"All premium features"}' --json

# 4. Confirm.
rc internal offerings get -p "P" -o "Pro Paywall" --json
```

Step 2 is not optional: metadata is a free-form object with no server-side schema,
so the only way to know the existing key names is to look.

## Reference identifiers

`-o` and `-p` both accept either the ID or the human name:

| Flag | Accepts |
|---|---|
| `-p, --project-id` | `proj…` ID, or project name |
| `-o, --offering-id` | `ofrng…` ID, offering identifier, or display name |

A name matching more than one object is an error listing the candidates — it never
picks one for you.

## Commands

```bash
rc internal offerings list --json
rc internal offerings list --platform IOS --json     # IOS, ANDROID, MACOS, WINDOWS, LINUX, tvOS
rc internal offerings get -o <ref> --json
rc internal offerings create -i weekly -n "Weekly"   # only -i and -n exist
rc internal offerings update -o <ref> -n "New Name"
rc internal offerings update -o <ref> -i new-identifier
rc internal offerings duplicate -o <ref> -i new-id -n "New Name"
rc internal offerings duplicate -o <ref> -i new-id -n "New Name" --packages-only
rc internal offerings set-current -o <ref>
rc internal offerings archive -o <ref>
rc internal offerings delete -o <ref>
```

## Attaching packages

`--packages` uses **snake_case** keys (`display_name`, `product_id`), unlike the
metadata object, which is free-form and conventionally camelCase.

```bash
rc internal offerings update -o <ref> --packages '[
  {
    "identifier": "my.weekly",
    "display_name": "Weekly",
    "products": [{ "product_id": "prodXXXX" }]
  }
]'
```

Get `product_id` values from `rc internal products list --json`.

## Metadata schema

Metadata is free-form: the server stores whatever object you send. The keys below
are the convention this project's paywalls read — they are **not** enforced, so
confirm against `offerings get --json` before assuming a key exists.

| Field | Type | Description |
|---|---|---|
| `title` | string | Paywall title |
| `subtitle` | string | Paywall subtitle |
| `features` | string[] | Feature bullets |
| `colorScheme.dark` / `.light` | object | `{ primary, background, secondaryText }` hex colours |
| `packages` | object[] | Presentation per package (see below) |
| `packages[].rcIdentifier` | string | Matches the package identifier (`weekly`, `annual`, `lifetime`) |
| `packages[].type` | string | Package type |
| `packages[].title` / `.subtitle` | string | Display strings |
| `packages[].shouldBeGivenProminance` | boolean | Feature prominently (note the spelling) |
| `primarySelectedPackageIndex` | number | Default selection, 0-indexed |
| `headerImageURL` | string | Header image |
| `privacyPolicyURL`, `termsAndConditionsURL` | string | Legal links |
| `reviews` | object[] | `{ review, reviewer, stars }` |

`metadata.packages` is presentation only. It does **not** attach products — that is
`--packages` (snake_case, above). The two are different things with similar names.

### Full replacement example

```bash
rc internal offerings update -o <ref> --metadata '{
  "title": "Unlock Pro",
  "subtitle": "Get access to all premium features",
  "features": ["Unlimited habits", "Cloud sync", "Custom themes"],
  "colorScheme": {
    "dark":  { "primary": "#6366F1", "background": "#1E1B4B", "secondaryText": "#A5B4FC" },
    "light": { "primary": "#4F46E5", "background": "#FFFFFF", "secondaryText": "#4338CA" }
  },
  "packages": [
    { "rcIdentifier": "weekly",   "type": "weekly",   "title": "Weekly",   "subtitle": "Try Pro this week", "shouldBeGivenProminance": false },
    { "rcIdentifier": "annual",   "type": "annual",   "title": "Yearly",   "subtitle": "Best value",        "shouldBeGivenProminance": true },
    { "rcIdentifier": "lifetime", "type": "lifetime", "title": "Lifetime", "subtitle": "One-time unlock" }
  ],
  "primarySelectedPackageIndex": 1,
  "privacyPolicyURL": "https://example.com/privacy",
  "termsAndConditionsURL": "https://example.com/terms",
  "reviews": [{ "review": "Great app!", "reviewer": "User", "stars": 5 }]
}'
```
