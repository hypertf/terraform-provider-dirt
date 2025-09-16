#!/bin/bash

REPO="hypertf/terraform-provider-dirt"

echo "=== GitHub Actions Status ==="
curl -s "https://api.github.com/repos/$REPO/actions/runs?per_page=5" | jq -r '
.workflow_runs[] | 
"[\(.created_at | strptime("%Y-%m-%dT%H:%M:%SZ") | strftime("%H:%M:%S"))] \(.name) on \(.head_branch // "unknown"): \(.status) (\(.conclusion // "running"))"
'

echo ""
echo "=== Latest Releases ==="
curl -s "https://api.github.com/repos/$REPO/releases?per_page=3" | jq -r '
.[] | 
"[\(.published_at | strptime("%Y-%m-%dT%H:%M:%SZ") | strftime("%H:%M:%S"))] \(.tag_name): \(.name // "Unnamed")"
'

echo ""
echo "=== Available Tags ==="
curl -s "https://api.github.com/repos/$REPO/tags?per_page=5" | jq -r '.[] | .name'
