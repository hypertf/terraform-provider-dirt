#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


REPO="hypertf/terraform-provider-dirt"
TAG="v0.1.2"
MAX_ATTEMPTS=30
SLEEP_INTERVAL=10

echo "Polling for GitHub release $TAG..."
echo "Repository: $REPO"
echo "Max attempts: $MAX_ATTEMPTS (${MAX_ATTEMPTS}0 seconds total)"
echo "Checking every $SLEEP_INTERVAL seconds"
echo ""

for i in $(seq 1 $MAX_ATTEMPTS); do
    echo "Attempt $i/$MAX_ATTEMPTS..."
    
    # Check if release exists
    response=$(curl -s "https://api.github.com/repos/$REPO/releases/tags/$TAG")
    
    if echo "$response" | grep -q '"tag_name"'; then
        echo "✅ Release $TAG found!"
        echo "Release URL: $(echo "$response" | grep '"html_url"' | head -1 | cut -d'"' -f4)"
        echo "Published at: $(echo "$response" | grep '"published_at"' | cut -d'"' -f4)"
        echo "Assets:"
        echo "$response" | grep '"browser_download_url"' | cut -d'"' -f4 | sed 's/^/  - /'
        exit 0
    elif echo "$response" | grep -q "Not Found"; then
        echo "❌ Release $TAG not found yet"
    else
        echo "⚠️  Unexpected response"
    fi
    
    if [ $i -lt $MAX_ATTEMPTS ]; then
        echo "Waiting $SLEEP_INTERVAL seconds..."
        sleep $SLEEP_INTERVAL
        echo ""
    fi
done

echo "❌ Release $TAG was not found after ${MAX_ATTEMPTS} attempts"
echo "You may need to check GitHub Actions or wait longer for the release to be created"
exit 1
