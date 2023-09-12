#!/bin/bash
set -eu

mkdir -p bin
go build -o bin/central
exec sudo ./bin/central
