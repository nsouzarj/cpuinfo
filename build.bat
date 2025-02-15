@echo off
go clean -modcache
go mod tidy
go mod download
go build -v 