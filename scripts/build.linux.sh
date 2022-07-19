#!/usr/bin/env bash
cd ..
go build -o api-server -a -ldflags "-extldflags '-static' -X 'main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" ./cmd/api/main.go
chmod +x ./api-server
echo "api-server done."

go build -o web-server -a -ldflags "-extldflags '-static' -X 'main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.GitTag=$(git describe --tags --abbrev=0)'" ./cmd/web/main.go
chmod +x ./web-server
echo "web-server built."