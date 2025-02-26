#!/bin/bash
set -e

echo "Compiling twigo..."
cd "$(dirname "$0")/.."
go build -o build/twigo main.go

echo "Done! Binary is at build/twigo" 