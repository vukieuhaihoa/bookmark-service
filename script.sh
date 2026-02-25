#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

find /Users/haihoa/Documents/learn/EBVN/Go/BE01/bookmark-management-ddd/bookmark-service -name "*.go" -exec sed -i '' 's|github.com/vukieuhaihoa/bookmark-management/internal|github.com/vukieuhaihoa/bookmark-service/internal|g' {} +

find /Users/haihoa/Documents/learn/EBVN/Go/BE01/bookmark-management-ddd/bookmark-service -name "*.go" -exec sed -i '' 's|github.com/vukieuhaihoa/bookmark-management/pkg|github.com/vukieuhaihoa/bookmark-libs/pkg|g' {} +

go mod tidy

make test