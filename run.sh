# !#/bin/bash

go build
go install
forge new example
cd example
forge init
forge build