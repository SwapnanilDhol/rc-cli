#!/usr/bin/env python3
"""Sequential RevenueCat public v2 GET crawl: all projects, then per-project resources.

Reads API key from ~/.revenuerc (apiKey). Writes markdown to the path passed as argv[1].

Does not print the API key. Customer listing uses a small limit to reduce PII in the output file.
"""

from __future__ import annotations

import json
import sys
import urllib.error
import urllib.parse
import urllib.request
from datetime import datetime, timezone
from pathlib import Path


BASE = "https://api.revenuecat.com/v2"


def load_api_key() -> str:
    cfg = Path.home() / ".revenuerc"
    if not cfg.is_file():
        sys.stderr.write("Missing ~/.revenuerc\n")
        sys.exit(1)
    data = json.loads(cfg.read_text(encoding="utf-8"))
    key = data.get("apiKey") or data.get("api_key")
    if not key:
        sys.stderr.write("No apiKey in ~/.revenuerc\n")
        sys.exit(1)
    return key


def api_get(key: str, path: str, params: dict[str, str] | None = None) -> tuple[int, str]:
    url = BASE + path
    if params:
        url += "?" + urllib.parse.urlencode(params)
    req = urllib.request.Request(
        url,
        headers={
            "Authorization": f"Bearer {key}",
            "Accept": "application/json",
        },
        method="GET",
    )
    try:
        with urllib.request.urlopen(req, timeout=120) as resp:
            raw = resp.read().decode("utf-8", errors="replace")
            return resp.status, raw
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8", errors="replace")
        return e.code, raw


def pretty_json(raw: str) -> str:
    if not raw.strip():
        return "_empty body_"
    try:
        return json.dumps(json.loads(raw), indent=2)
    except json.JSONDecodeError:
        return raw


def starting_after_from_next(next_page: str | None) -> str | None:
    if not next_page:
        return None
    q = urllib.parse.urlparse(next_page).query
    d = urllib.parse.parse_qs(q)
    vals = d.get("starting_after")
    return vals[0] if vals else None


def fetch_all_items(key: str, path: str, extra: dict[str, str] | None = None) -> tuple[int, list, str]:
    """Returns (http_status, items, raw_body_on_first_page). On non-200, items is []."""
    params = {"limit": "100"}
    if extra:
        params.update(extra)
    items: list = []
    first_raw = ""
    while True:
        status, raw = api_get(key, path, params)
        if first_raw == "":
            first_raw = raw
        if status != 200:
            return status, [], first_raw
        body = json.loads(raw)
        chunk = body.get("items") or []
        items.extend(chunk)
        nxt = starting_after_from_next(body.get("next_page"))
        if not nxt:
            break
        params = {"limit": "100", "starting_after": nxt}
        if extra:
            for k, v in extra.items():
                if k != "starting_after":
                    params[k] = v
    return 200, items, first_raw


def md_call(title: str, path: str, status: int, raw: str) -> str:
    body = pretty_json(raw) if status == 200 else raw
    return f"### {title}\n\n`GET {path}` → **HTTP {status}**\n\n```json\n{body}\n```\n\n"


def main() -> None:
    out_path = Path(sys.argv[1]) if len(sys.argv) > 1 else Path("docs/revenuecat-v2-api-dump.md")
    key = load_api_key()
    lines: list[str] = []

    lines.append("# RevenueCat API v2 — project crawl\n\n")
    lines.append(
        "> **Privacy:** This file may include customer identifiers, audit logs (emails, IPs, user agents), and catalog data. "
        "Keep it local or redact before sharing. Default output path is gitignored.\n\n"
    )
    lines.append(
        f"Generated **{datetime.now(timezone.utc).isoformat()}** via sequential public v2 GETs.\n\n"
    )
    lines.append(
        "Official reference: [Developer API v2](https://www.revenuecat.com/docs/api-v2).\n\n"
        "Re-run: `python3 tools/dump_v2_projects.py [out.md]`\n\n"
        "---\n\n## 1. List all projects\n\n"
    )

    proj_status, proj_raw = api_get(key, "/projects", {"limit": "100"})
    projects: list[dict] = []
    if proj_status == 200:
        parsed = json.loads(proj_raw)
        projects = list(parsed.get("items") or [])
        cursor = starting_after_from_next(parsed.get("next_page"))
        while cursor:
            st, raw = api_get(key, "/projects", {"limit": "100", "starting_after": cursor})
            if st != 200:
                lines.append(f"_Pagination error HTTP {st}_\n\n")
                break
            p2 = json.loads(raw)
            projects.extend(p2.get("items") or [])
            cursor = starting_after_from_next(p2.get("next_page"))

    lines.append(md_call("GET /projects (aggregated pages)", "/projects?limit=100…", 200, json.dumps({"items": projects})))

    for proj in projects:
        pid = proj.get("id")
        pname = proj.get("name", pid)
        if not pid:
            continue
        lines.append(f"---\n\n## Project `{pid}` — {pname}\n\n")

        lines.append(
            "> **Note:** `GET /projects/{project_id}` is **not** listed in the current [API v2](https://www.revenuecat.com/docs/api-v2) "
            "reference (project rows come from `GET /projects`). The call below is included for completeness.\n\n"
        )
        st, raw = api_get(key, f"/projects/{pid}")
        lines.append(md_call("Project detail (optional / non-doc)", f"/projects/{pid}", st, raw))

        def crawl_list(title: str, rel: str, extra: dict[str, str] | None = None):
            st2, items, raw0 = fetch_all_items(key, rel, extra)
            payload = json.dumps({"items": items}, indent=2) if st2 == 200 else raw0
            lines.append(md_call(f"{title} ({len(items)} items)", rel, st2, payload))
            return items

        apps = crawl_list("Apps", f"/projects/{pid}/apps")
        for app in apps:
            aid = app.get("id")
            if not aid:
                continue
            st, raw = api_get(key, f"/projects/{pid}/apps/{aid}")
            lines.append(md_call(f"App `{aid}`", f"/projects/{pid}/apps/{aid}", st, raw))
            st, raw = api_get(key, f"/projects/{pid}/apps/{aid}/public_api_keys")
            lines.append(md_call(f"App `{aid}` public API keys", f"/projects/{pid}/apps/{aid}/public_api_keys", st, raw))
            st, raw = api_get(key, f"/projects/{pid}/apps/{aid}/store_kit_config")
            lines.append(md_call(f"App `{aid}` StoreKit config", f"/projects/{pid}/apps/{aid}/store_kit_config", st, raw))

        products = crawl_list("Products", f"/projects/{pid}/products")
        for p in products:
            prid = p.get("id")
            if not prid:
                continue
            st, raw = api_get(key, f"/projects/{pid}/products/{prid}")
            lines.append(md_call(f"Product `{prid}`", f"/projects/{pid}/products/{prid}", st, raw))
        ents = crawl_list("Entitlements", f"/projects/{pid}/entitlements")
        for ent in ents:
            eid = ent.get("id")
            if not eid:
                continue
            st, raw = api_get(key, f"/projects/{pid}/entitlements/{eid}")
            lines.append(md_call(f"Entitlement `{eid}`", f"/projects/{pid}/entitlements/{eid}", st, raw))
            st2, raw2 = api_get(key, f"/projects/{pid}/entitlements/{eid}/products")
            lines.append(md_call(f"Entitlement `{eid}` products", f"/projects/{pid}/entitlements/{eid}/products", st2, raw2))

        offs = crawl_list("Offerings", f"/projects/{pid}/offerings")
        for off in offs:
            oid = off.get("id")
            if not oid:
                continue
            st, raw = api_get(key, f"/projects/{pid}/offerings/{oid}", {"expand": "package.product"})
            lines.append(
                md_call(f"Offering `{oid}` (expand=package.product)", f"/projects/{pid}/offerings/{oid}", st, raw)
            )
            pkgs = crawl_list(f"Packages in offering `{oid}`", f"/projects/{pid}/offerings/{oid}/packages")
            for pkg in pkgs:
                pkg_id = pkg.get("id")
                if not pkg_id:
                    continue
                st, raw = api_get(key, f"/projects/{pid}/packages/{pkg_id}", {"expand": "product"})
                lines.append(md_call(f"Package `{pkg_id}`", f"/projects/{pid}/packages/{pkg_id}", st, raw))
                st2, raw2 = api_get(key, f"/projects/{pid}/packages/{pkg_id}/products")
                lines.append(md_call(f"Package `{pkg_id}` products", f"/projects/{pid}/packages/{pkg_id}/products", st2, raw2))

        wh = crawl_list("Webhook integrations (v2 path)", f"/projects/{pid}/integrations/webhooks")
        for w in wh:
            wid = w.get("id")
            if not wid:
                continue
            st, raw = api_get(key, f"/projects/{pid}/integrations/webhooks/{wid}")
            lines.append(md_call(f"Webhook integration `{wid}`", f"/projects/{pid}/integrations/webhooks/{wid}", st, raw))

        st, raw = api_get(key, f"/projects/{pid}/offers")
        lines.append(md_call("Promotional offers (legacy path; may 404 in v2)", f"/projects/{pid}/offers", st, raw))

        for path, label, q in (
            (f"/projects/{pid}/collaborators", "Collaborators", None),
            (f"/projects/{pid}/audit_logs", "Audit logs", None),
            (f"/projects/{pid}/metrics/overview", "Metrics overview", None),
            # Project-level subscription/purchase lists require filters per v2 docs; omit noisy 400s.
        ):
            st, raw = api_get(key, path, q)
            qs = ("?" + urllib.parse.urlencode(q)) if q else ""
            lines.append(md_call(label, path + qs, st, raw))

        st, raw = api_get(
            key,
            f"/projects/{pid}/customers",
            {"limit": "10"},
        )
        lines.append(
            md_call(
                "Customers sample (limit=10; not fully paginated — avoids huge PII export)",
                f"/projects/{pid}/customers?limit=10",
                st,
                raw,
            )
        )

        vcs = crawl_list("Virtual currencies", f"/projects/{pid}/virtual_currencies")
        for vc in vcs:
            code = vc.get("code") or vc.get("id")
            if not code:
                continue
            st, raw = api_get(key, f"/projects/{pid}/virtual_currencies/{code}")
            lines.append(md_call(f"Virtual currency `{code}`", f"/projects/{pid}/virtual_currencies/{code}", st, raw))

        st, raw = api_get(key, f"/projects/{pid}/paywalls", {"limit": "50"})
        lines.append(md_call("Paywalls", f"/projects/{pid}/paywalls", st, raw))
        if st == 200:
            try:
                for pw in json.loads(raw).get("items") or []:
                    pwid = pw.get("id")
                    if not pwid:
                        continue
                    st2, raw2 = api_get(key, f"/projects/{pid}/paywalls/{pwid}")
                    lines.append(md_call(f"Paywall `{pwid}`", f"/projects/{pid}/paywalls/{pwid}", st2, raw2))
            except json.JSONDecodeError:
                pass

    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text("".join(lines), encoding="utf-8")
    print(f"Wrote {out_path}", file=sys.stderr)


if __name__ == "__main__":
    main()
