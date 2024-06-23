#!/bin/bash
GOARCH=wasm GOOS=js go build -o web/app.wasm ./frontend/cmd