#!/bin/bash

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"
echo "Script dir: $SCRIPT_DIR"
echo "Project root: $PROJECT_ROOT"

cd "$PROJECT_ROOT"

# # Autenticarse si no lo has hecho
# docker login

# Construir la imagen multi-arquitectura
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7,linux/386 \
  -f docker/Dockerfile \
  -t rueda1208/workspace-telegraf:1.0.0 \
  -t rueda1208/workspace-telegraf:latest \
  --push .
