#!/bin/bash
# Integration helper script for Tallarin Activation Server
# This script helps extract the public key and update the Go application

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ACTIVATION_SERVER_DIR="$SCRIPT_DIR"
GO_APP_DIR="$(dirname "$SCRIPT_DIR")"

echo "Tallarin Activation Server Integration Helper"
echo "=============================================="

# Check if we're in the right directory
if [ ! -f "$ACTIVATION_SERVER_DIR/app.py" ]; then
    echo "Error: app.py not found. Please run this script from the activation_server directory."
    exit 1
fi

# Check if public key exists
if [ ! -f "$ACTIVATION_SERVER_DIR/public_key.pem" ]; then
    echo "Error: public_key.pem not found. Please run the activation server first to generate keys."
    echo "Run: python app.py"
    exit 1
fi

echo "1. Extracting public key from activation server..."
PUBLIC_KEY=$(cat "$ACTIVATION_SERVER_DIR/public_key.pem")

echo "2. Creating Go code snippet for config_activation.go..."
cat > "$ACTIVATION_SERVER_DIR/go_integration_snippet.txt" << EOF
// Updated public key from activation server
// Replace the publicKey variable in config/config_activation.go with this:
var publicKey = \`$PUBLIC_KEY\`
EOF

echo "3. Checking for revoked licenses..."
if [ -f "$ACTIVATION_SERVER_DIR/revoked_licenses.json" ]; then
    REVOKED_LICENSES=$(python3 -c "
import json
try:
    with open('$ACTIVATION_SERVER_DIR/revoked_licenses.json', 'r') as f:
        data = json.load(f)
    if data:
        licenses = ', '.join(['\"{}\"'.format(license_id) for license_id in data])
        print('var revocations = []string{{{}}}'.format(licenses))
    else:
        print('var revocations = []string{}')
except:
    print('var revocations = []string{}')
")
    
    cat >> "$ACTIVATION_SERVER_DIR/go_integration_snippet.txt" << EOF

// Updated revocation list from activation server
// Replace the revocations variable in config/config_activation.go with this:
$REVOKED_LICENSES
EOF
else
    cat >> "$ACTIVATION_SERVER_DIR/go_integration_snippet.txt" << EOF

// No revoked licenses yet
var revocations = []string{}
EOF
fi

echo "4. Integration instructions created!"
echo ""
echo "Next steps:"
echo "1. Copy the content from go_integration_snippet.txt"
echo "2. Edit $GO_APP_DIR/config/config_activation.go"
echo "3. Replace the publicKey and revocations variables with the new values"
echo "4. Rebuild the Go application: go build -o r3"
echo ""
echo "Files created:"
echo "- go_integration_snippet.txt (contains Go code to copy)"
echo ""

echo "Integration helper completed!"
echo "View the integration snippet:"
echo "cat $ACTIVATION_SERVER_DIR/go_integration_snippet.txt"