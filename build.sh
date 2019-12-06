#!/usr/bin/env bash 
set -xe

# install packages and dependencies
go get github.com/gin-gonic/contrib
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get golang.org/x/oauth2
go get github.com/dgrijalva/jwt-go

# build command
go build -o bin/application application.go
