#!/usr/bin/env python3
"""Reconcile a RevenueCat project's product catalog against App Store Connect.

App Store Connect is the source of truth: a product that exists in ASC but not in
RevenueCat means users can complete a purchase that RevenueCat will not recognise,
so no entitlement is granted. That is the failure this script exists to catch.

The two systems join on the product identifier — ASC `attributes.productId` equals
RevenueCat `identifier` — and the app is matched automatically by bundle ID
(RevenueCat `platform_identifier` equals ASC `bundleId`).

Requires: `rc login` (RevenueCat dashboard) and `asc auth` (App Store Connect).

Usage:
  ./scripts/rc-asc-reconcile.py --project "Neon"
  ./scripts/rc-asc-reconcile.py --project 1a2b3c4d --json
  ./scripts/rc-asc-reconcile.py --project "Neon" --strict   # also fail on RC-only
"""

import argparse
import json
import shutil
import subprocess
import sys


def run(cmd):
    """Run a command and parse its JSON stdout."""
    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0:
        detail = (proc.stderr or proc.stdout).strip().splitlines()
        raise SystemExit(f"$ {' '.join(cmd)}\n{detail[0] if detail else 'failed'}")
    try:
        return json.loads(proc.stdout)
    except json.JSONDecodeError:
        raise SystemExit(f"$ {' '.join(cmd)}\nexpected JSON, got: {proc.stdout[:200]}")


def as_list(payload, *keys):
    """Both CLIs vary between a bare array and an object wrapping one."""
    if isinstance(payload, list):
        return payload
    for k in keys:
        if isinstance(payload.get(k), list):
            return payload[k]
    return []


def rc_products(rc, project):
    payload = run([rc, "internal", "products", "list", "-p", project,
                   "--limit", "2500", "--json"])
    return {
        p["identifier"]: p
        for p in as_list(payload, "products", "data")
        if p.get("identifier")
    }


def rc_bundle_id(rc, project):
    """Find the App Store app's bundle ID for a RevenueCat project."""
    apps = as_list(run([rc, "internal", "apps", "list", "-p", project, "--json"]),
                   "apps", "data")
    store_apps = [a for a in apps if a.get("type") == "app_store"]
    if not store_apps:
        raise SystemExit(f"project {project!r} has no App Store app in RevenueCat")
    bundle = store_apps[0].get("platform_identifier")
    if not bundle:
        raise SystemExit(
            f"project {project!r}: App Store app {store_apps[0].get('id')} has no bundle ID set "
            "in RevenueCat, so it cannot be matched to App Store Connect")
    return bundle


def asc_app_id(asc, bundle):
    apps = as_list(run([asc, "apps", "list", "--output", "json"]), "data", "apps")
    for a in apps:
        attrs = a.get("attributes", a)
        if attrs.get("bundleId") == bundle:
            return a.get("id"), attrs.get("name")
    raise SystemExit(f"no App Store Connect app found with bundle ID {bundle!r}")


def asc_products(asc, app_id):
    """Every purchasable product in ASC, keyed by productId."""
    out = {}
    for args, kind in ((["subscriptions", "list", "--app", app_id], "subscription"),
                       (["iap", "list", "--app", app_id], "iap")):
        for item in as_list(run([asc, *args, "--output", "json"]), "data"):
            attrs = item.get("attributes", {})
            pid = attrs.get("productId")
            if pid:
                out[pid] = {"kind": kind, "id": item.get("id"),
                            "name": attrs.get("name"), "state": attrs.get("state")}
    return out


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--project", required=True, help="RevenueCat project name or ID")
    ap.add_argument("--json", action="store_true", help="machine-readable output")
    ap.add_argument("--strict", action="store_true",
                    help="also exit non-zero when RevenueCat has products ASC does not")
    ap.add_argument("--rc", default="./rc", help="path to the rc binary")
    ap.add_argument("--asc", default="asc", help="path to the asc binary")
    args = ap.parse_args()

    rc = shutil.which(args.rc) or args.rc
    asc = shutil.which(args.asc) or args.asc
    for path, name in ((rc, "rc"), (asc, "asc")):
        if not shutil.which(path):
            raise SystemExit(f"{name} not found at {path!r}")

    bundle = rc_bundle_id(rc, args.project)
    app_id, app_name = asc_app_id(asc, bundle)
    rc_by_id = rc_products(rc, args.project)
    asc_by_id = asc_products(asc, app_id)

    missing_in_rc = sorted(set(asc_by_id) - set(rc_by_id))
    missing_in_asc = sorted(set(rc_by_id) - set(asc_by_id))
    matched = sorted(set(rc_by_id) & set(asc_by_id))

    if args.json:
        json.dump({
            "project": args.project,
            "bundle_id": bundle,
            "asc_app_id": app_id,
            "asc_app_name": app_name,
            "matched": matched,
            "missing_in_revenuecat": [
                {"product_id": p, **asc_by_id[p]} for p in missing_in_rc],
            "missing_in_app_store_connect": [
                {"product_id": p,
                 "revenuecat_id": rc_by_id[p].get("id"),
                 "product_type": rc_by_id[p].get("product_type"),
                 "entitlements": len(rc_by_id[p].get("entitlements") or [])}
                for p in missing_in_asc],
        }, sys.stdout, indent=2)
        print()
    else:
        print(f"{app_name}  ({bundle}, ASC app {app_id})")
        print(f"  RevenueCat {len(rc_by_id)}   App Store Connect {len(asc_by_id)}   "
              f"matched {len(matched)}\n")

        if missing_in_rc:
            print(f"MISSING IN REVENUECAT ({len(missing_in_rc)}) — "
                  "purchasable, but no entitlement will be granted:")
            for p in missing_in_rc:
                d = asc_by_id[p]
                print(f"    {p}  [{d['kind']}, {d['state']}]")
            print("\n  Fix with:")
            for p in missing_in_rc:
                d = asc_by_id[p]
                kind = "subscription" if d["kind"] == "subscription" else "non_consumable_product"
                print(f"    rc internal products create -p {args.project!r} "
                      f"--app-id <app_id> --product-type {kind} "
                      f"--identifier {p} --name {(d['name'] or p)!r}")
            print()

        if missing_in_asc:
            print(f"NOT IN APP STORE CONNECT ({len(missing_in_asc)}) — "
                  "stale RevenueCat entries, usually harmless:")
            for p in missing_in_asc:
                n = len(rc_by_id[p].get("entitlements") or [])
                warn = "  <-- attached to an entitlement" if n else ""
                print(f"    {p}  [{rc_by_id[p].get('product_type')}]{warn}")
            print()

        if not missing_in_rc and not missing_in_asc:
            print("  In sync.")

    if missing_in_rc:
        return 1
    if args.strict and missing_in_asc:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
