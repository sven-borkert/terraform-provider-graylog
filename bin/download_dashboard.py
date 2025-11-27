#!/usr/bin/env python3
"""
Download a Graylog dashboard via API and save it as JSON.
Usage: python download_dashboard.py [dashboard_id] [output_file]

If no arguments provided, downloads the PacketBeat dashboard.
"""

import json
import sys
import requests
from datetime import datetime
from urllib3.exceptions import InsecureRequestWarning

# Suppress SSL warnings for self-signed certs
requests.packages.urllib3.disable_warnings(category=InsecureRequestWarning)

# Configuration
GRAYLOG_URL = "https://graylog.internal.borkert.net/api"
USERNAME = "admin"
PASSWORD = "1JTxlqJArpswjUqcsboLoQ=="

# PacketBeat dashboard ID (from terraform output)
DEFAULT_DASHBOARD_ID = "6928b7ca524162443155b69b"


def get_dashboard(dashboard_id: str) -> dict:
    """Fetch dashboard from Graylog API."""
    url = f"{GRAYLOG_URL}/views/{dashboard_id}"
    response = requests.get(
        url,
        auth=(USERNAME, PASSWORD),
        headers={"X-Requested-By": "dashboard-downloader"},
        verify=False,
    )
    response.raise_for_status()
    return response.json()


def get_search(search_id: str) -> dict:
    """Fetch search from Graylog API."""
    url = f"{GRAYLOG_URL}/views/search/{search_id}"
    response = requests.get(
        url,
        auth=(USERNAME, PASSWORD),
        headers={"X-Requested-By": "dashboard-downloader"},
        verify=False,
    )
    response.raise_for_status()
    return response.json()


def save_dashboard(data: dict, filename: str):
    """Save dashboard JSON to file with pretty formatting."""
    with open(filename, "w") as f:
        json.dump(data, f, indent=2, sort_keys=True)
    print(f"Saved dashboard to {filename}")


def main():
    dashboard_id = sys.argv[1] if len(sys.argv) > 1 else DEFAULT_DASHBOARD_ID

    # Generate filename with timestamp if not provided
    if len(sys.argv) > 2:
        output_file = sys.argv[2]
    else:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_file = f"dashboard_{timestamp}.json"

    print(f"Downloading dashboard {dashboard_id}...")
    dashboard = get_dashboard(dashboard_id)

    print(f"Title: {dashboard.get('title', 'Unknown')}")
    print(f"ID: {dashboard.get('id', 'Unknown')}")

    save_dashboard(dashboard, output_file)

    # Also print state keys for quick inspection
    if "state" in dashboard:
        print(f"\nState keys: {list(dashboard['state'].keys())}")
        for state_id, state_data in dashboard["state"].items():
            if "widgets" in state_data:
                print(f"  Widgets in state '{state_id}': {len(state_data['widgets'])}")


if __name__ == "__main__":
    main()
