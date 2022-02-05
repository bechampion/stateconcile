#!/bin/bash
set -euo pipefail
key_json_file="$1"
scope="$2"
valid_for_sec="${3:-3600}"
IFS="," read -r private_key sa_email <<<$(cat $1 | jq '(.private_key|@base64)+","+.client_email' -r)
header='{"alg":"RS256","typ":"JWT"}'
claim=$(cat <<EOF | jq -c .
  {
    "iss": "$sa_email",
    "scope": "$scope",
    "aud": "https://www.googleapis.com/oauth2/v4/token",
    "exp": $(($(date +%s) + $valid_for_sec)),
    "iat": $(date +%s)
  }
EOF
     )
request_body="$(echo "$header" | base64 -w0).$(echo "$claim" | base64 -w0)"
signature=$(openssl dgst -sha256 -sign <(echo "$private_key"| base64 -d ) <(printf "$request_body") | base64 -w0)
printf "$request_body.$signature"
