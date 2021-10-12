#!/bin/sh

# ===
# Generate Docs
# ===
swag init --parseDependency --parseInternal

# ===
# Obtain a public key
# ===
mkdir -p keys

# Get the public key of Auth Service
curl http://api.hiringo.tech/auth/public-key -o keys/public.pem