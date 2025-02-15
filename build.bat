@echo off
go clean -modcache
go mod tidy
go mod download
CGO_ENABLED=1 GOOS=windows go build -v 