#!/usr/bin/env python3
"""
Fetch all Graylog API resource docs and store them under docs/api-docs/.

Usage:
  python docs/fetch_api_docs.py --base-url https://graylog.internal.borkert.net --user admin --password YOUR_PASSWORD

Notes:
  - Requires the top-level API docs endpoint at <base-url>/api/api-docs/ (Swagger 1.2 style).
  - TLS verification is enabled by default; pass --insecure to skip (for self-signed lab certs).
"""

import argparse
import json
import os
import pathlib
import re
import sys
from typing import Any, Dict, List

import requests


def fetch(url: str, auth, verify: bool) -> Dict[str, Any]:
    resp = requests.get(url, auth=auth, verify=verify)
    resp.raise_for_status()
    return resp.json()


def sanitize_path(path: str) -> str:
    # Convert swagger path like /system/indices/index_sets/templates to system_indices_index_sets_templates
    return re.sub(r"[^A-Za-z0-9]+", "_", path.strip("/")).strip("_")


def main():
    parser = argparse.ArgumentParser(description="Fetch all Graylog API resource docs")
    parser.add_argument("--base-url", required=True, help="Graylog base URL (e.g., https://graylog.internal.borkert.net)")
    parser.add_argument("--user", required=True, help="Graylog username (or token)")
    parser.add_argument("--password", required=True, help="Graylog password (or 'token')")
    parser.add_argument("--out-dir", default="docs/api-docs", help="Output directory for JSON docs")
    parser.add_argument("--insecure", action="store_true", help="Skip TLS verification (self-signed)")
    args = parser.parse_args()

    base = args.base_url.rstrip("/")
    root_url = f"{base}/api/api-docs/"
    out_dir = pathlib.Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)

    verify = not args.insecure
    auth = (args.user, args.password)

    print(f"[+] Fetching index: {root_url}")
    index = fetch(root_url, auth, verify)
    apis: List[Dict[str, Any]] = index.get("apis", [])
    if not apis:
        print("[-] No APIs found in index; exiting.", file=sys.stderr)
        sys.exit(1)

    # Save index
    (out_dir / "api-index.json").write_text(json.dumps(index, indent=2))
    print(f"[+] Saved index to {out_dir / 'api-index.json'}")

    for api in apis:
        path = api.get("path")
        if not path:
            continue
        resource_url = f"{base}/api/api-docs{path}"
        fname = sanitize_path(path) or "root"
        target = out_dir / f"{fname}.json"
        try:
            print(f"[+] Fetching {resource_url}")
            data = fetch(resource_url, auth, verify)
            target.write_text(json.dumps(data, indent=2))
            print(f"[+] Saved {target}")
        except Exception as exc:  # noqa: BLE001
            print(f"[-] Failed to fetch {resource_url}: {exc}", file=sys.stderr)


if __name__ == "__main__":
    main()
