#!/bin/bash

# Очистка неиспользуемых зависимостей
echo "Running go mod tidy..."
go mod tidy
go run main.go