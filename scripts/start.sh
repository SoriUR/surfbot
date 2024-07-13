#!/bin/bash

cd "$(dirname "$0")"
cd ..
go mod tidy
go run main.go