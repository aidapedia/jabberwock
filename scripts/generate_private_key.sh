#!/bin/bash
set -euo pipefail

SECRET_JSON="files/config/secret.json"

echo "Generating RSA private key..."
openssl genrsa -out private_key.pem 2048

# Write the key into files/config/secret.json
if command -v jq >/dev/null 2>&1; then
  echo "jq detected; updating $SECRET_JSON"
  jq --rawfile pk private_key.pem '.auth.private_key = $pk' "$SECRET_JSON" > "$SECRET_JSON.tmp"
  mv "$SECRET_JSON.tmp" "$SECRET_JSON"
else
  echo "jq not found; using Python to update $SECRET_JSON"
  python3 - "$SECRET_JSON" "private_key.pem" << 'PY'
import sys, json
secret_path = sys.argv[1]
pk_path = sys.argv[2]
with open(secret_path, 'r') as f: data = json.load(f)
with open(pk_path, 'r') as f: pk = f.read()
# ensure structure
if 'auth' not in data or not isinstance(data['auth'], dict):
    data['auth'] = {}
data['auth']['private_key'] = pk
with open(secret_path, 'w') as f: json.dump(data, f, indent=4)
print("Updated", secret_path)
PY
fi

# Clean up
rm private_key.pem

echo "Done. Stored private key in $SECRET_JSON."
