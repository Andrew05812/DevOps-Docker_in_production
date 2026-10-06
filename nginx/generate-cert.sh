#!/bin/sh
mkdir -p nginx/ssl

openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
  -keyout nginx/ssl/server.key \
  -out nginx/ssl/server.crt \
  -subj "/CN=devops-stack.local" \
  -addext "subjectAltName=DNS:devops-stack.local,DNS:localhost,IP:127.0.0.1"

chmod 600 nginx/ssl/server.key
