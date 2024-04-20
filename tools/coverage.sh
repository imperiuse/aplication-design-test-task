#!/usr/bin/env sh

# FOR UNIT TEST ONLY
go test -short `go list ./... | grep -v tools` -coverprofile cover.out
go tool cover -func cover.out | grep total | awk '{print substr($3, 1, length($3)-1)}' | tee cover.score
