#!/usr/bin/env python3
"""
Flask-based Activation Server for Tallarin
This server generates, validates, and manages licenses for the application.
"""

import json
import base64
import hashlib
from datetime import datetime, timedelta
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.backends import default_backend
import os
from flask import Flask, request, jsonify, render_template_string

app = Flask(__name__)

# Configuration
PRIVATE_KEY_PATH = "private_key.pem"
PUBLIC_KEY_PATH = "public_key.pem"
REVOKED_LICENSES_PATH = "revoked_licenses.json"

# Generate or load RSA keys
def generate_or_load_keys():
    """Generate RSA key pair if not exists, otherwise load existing keys"""
    if not os.path.exists(PRIVATE_KEY_PATH) or not os.path.exists(PUBLIC_KEY_PATH):
        print("Generating new RSA key pair...")
        # Generate private key
        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=4096,
            backend=default_backend()
        )
        
        # Get public key
        public_key = private_key.public_key()
        
        # Save private key
        with open(PRIVATE_KEY_PATH, 'wb') as f:
            f.write(private_key.private_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PrivateFormat.PKCS8,
                encryption_algorithm=serialization.NoEncryption()
            ))
        
        # Save public key in PKCS1 format (same as Go app expects)
        with open(PUBLIC_KEY_PATH, 'wb') as f:
            f.write(public_key.public_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PublicFormat.PKCS1
            ))
        
        print(f"Keys saved to {PRIVATE_KEY_PATH} and {PUBLIC_KEY_PATH}")
    else:
        print("Loading existing RSA keys...")
        with open(PRIVATE_KEY_PATH, 'rb') as f:
            private_key = serialization.load_pem_private_key(
                f.read(), password=None, backend=default_backend()
            )
        with open(PUBLIC_KEY_PATH, 'rb') as f:
            public_key = serialization.load_pem_public_key(
                f.read(), backend=default_backend()
            )
    
    return private_key, public_key

def load_revoked_licenses():
    """Load revoked license IDs from file"""
    if os.path.exists(REVOKED_LICENSES_PATH):
        with open(REVOKED_LICENSES_PATH, 'r') as f:
            return json.load(f)
    return []

def save_revoked_licenses(revoked_list):
    """Save revoked license IDs to file"""
    with open(REVOKED_LICENSES_PATH, 'w') as f:
        json.dump(revoked_list, f, indent=2)

def generate_license_id():
    """Generate a unique license ID"""
    import uuid
    return f"LI{uuid.uuid4().hex[:8].upper()}"

def create_license_string(license_data):
    """Convert license data to new format string: client_id;registered_for;login_count_limit;valid_for_days;ext1,ext2,ext3"""
    client_id = license_data['clientId']
    registered_for = license_data['registeredFor']
    login_count = license_data['loginCount']
    
    # Calculate valid_for_days from validUntil timestamp
    valid_until_timestamp = license_data['validUntil']
    now = datetime.now()
    valid_until_date = datetime.fromtimestamp(valid_until_timestamp)
    valid_for_days = max(0, (valid_until_date - now).days)
    
    # Join extensions with commas
    extensions = license_data.get('extensions', [])
    extensions_str = ','.join(extensions) if extensions else ''
    
    # Create the semicolon-separated string
    license_string = f"{client_id};{registered_for};{login_count};{valid_for_days};{extensions_str}"
    
    # Encode to base64
    return base64.b64encode(license_string.encode('utf-8')).decode('utf-8')

def sign_license(license_data, private_key):
    """Sign license data with private key using new base64 format"""
    # Create the new license string format
    license_base64 = create_license_string(license_data)
    
    # Hash the base64 string
    message = license_base64.encode('utf-8')
    digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest.update(message)
    hashed = digest.finalize()
    
    # Sign the hash
    signature = private_key.sign(
        hashed,
        padding.PKCS1v15(),
        hashes.SHA256()
    )
    
    return base64.urlsafe_b64encode(signature).decode('utf-8')

# Initialize keys
private_key, public_key = generate_or_load_keys()

@app.route('/')
def index():
    """Simple web interface for license management"""
    return render_template_string('''
<!DOCTYPE html>
<html>
<head>
    <title>Tallarin Activation Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, select { padding: 8px; width: 300px; }
        button { padding: 10px 20px; background: #007cba; color: white; border: none; cursor: pointer; }
        button:hover { background: #005a87; }
        .result { margin: 20px 0; padding: 15px; background: #f0f0f0; border-radius: 4px; }
        .error { background: #ffe6e6; }
        .success { background: #e6ffe6; }
        textarea { width: 600px; height: 200px; }
    </style>
</head>
<body>
    <h1>Tallarin Activation Server</h1>
    
    <h2>Generate License</h2>
    <form id="generateForm">
        <div class="form-group">
            <label for="clientId">Client ID:</label>
            <input type="text" id="clientId" name="clientId" required>
        </div>
        <div class="form-group">
            <label for="registeredFor">Registered For:</label>
            <input type="text" id="registeredFor" name="registeredFor" required>
        </div>
        <div class="form-group">
            <label for="loginCount">Login Count Limit:</label>
            <input type="number" id="loginCount" name="loginCount" value="100" min="1">
        </div>
        <div class="form-group">
            <label for="validDays">Valid for (days):</label>
            <input type="number" id="validDays" name="validDays" value="365" min="1">
        </div>
        <div class="form-group">
            <label for="extensions">Extensions (comma-separated):</label>
            <input type="text" id="extensions" name="extensions" placeholder="extension1,extension2">
        </div>
        <button type="submit">Generate License</button>
    </form>
    
    <div id="generateResult" class="result" style="display:none;"></div>
    
    <h2>Validate License</h2>
    <form id="validateForm">
        <div class="form-group">
            <label for="licenseData">License Data (JSON):</label>
            <textarea id="licenseData" name="licenseData" placeholder="Paste license file content here..."></textarea>
        </div>
        <button type="submit">Validate License</button>
    </form>
    
    <div id="validateResult" class="result" style="display:none;"></div>
    
    <h2>Revoke License</h2>
    <form id="revokeForm">
        <div class="form-group">
            <label for="licenseId">License ID:</label>
            <input type="text" id="licenseId" name="licenseId" required>
        </div>
        <button type="submit">Revoke License</button>
    </form>
    
    <div id="revokeResult" class="result" style="display:none;"></div>
    
    <script>
        function showResult(elementId, message, isError = false) {
            const element = document.getElementById(elementId);
            element.innerHTML = message;
            element.className = 'result ' + (isError ? 'error' : 'success');
            element.style.display = 'block';
        }
        
        document.getElementById('generateForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData);
            
            if (data.extensions) {
                data.extensions = data.extensions.split(',').map(s => s.trim()).filter(s => s);
            } else {
                data.extensions = [];
            }
            data.loginCount = parseInt(data.loginCount);
            data.validDays = parseInt(data.validDays);
            
            try {
                const response = await fetch('/api/generate-license', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });
                const result = await response.json();
                
                if (response.ok) {
                    showResult('generateResult', 
                        `<strong>License Generated Successfully!</strong><br><br>
                         <strong>License ID:</strong> ${result.license.licenseId}<br><br>
                         <strong>License File Content:</strong><br>
                         <textarea readonly style="width:100%; height:150px;">${JSON.stringify(result, null, 2)}</textarea>
                         <br><br>Save this content to a .json file and upload it to your Tallarin application.`);
                } else {
                    showResult('generateResult', `<strong>Error:</strong> ${result.error}`, true);
                }
            } catch (error) {
                showResult('generateResult', `<strong>Error:</strong> ${error.message}`, true);
            }
        });
        
        document.getElementById('validateForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const licenseData = formData.get('licenseData');
            
            try {
                const parsedData = JSON.parse(licenseData);
                const response = await fetch('/api/validate-license', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(parsedData)
                });
                const result = await response.json();
                
                if (response.ok) {
                    showResult('validateResult', 
                        `<strong>License Validation Result:</strong><br><br>
                         <strong>Valid:</strong> ${result.valid ? 'Yes' : 'No'}<br>
                         <strong>Message:</strong> ${result.message}`);
                } else {
                    showResult('validateResult', `<strong>Error:</strong> ${result.error}`, true);
                }
            } catch (error) {
                showResult('validateResult', `<strong>Error:</strong> Invalid JSON or ${error.message}`, true);
            }
        });
        
        document.getElementById('revokeForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData);
            
            try {
                const response = await fetch('/api/revoke-license', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data)
                });
                const result = await response.json();
                
                if (response.ok) {
                    showResult('revokeResult', `<strong>Success:</strong> ${result.message}`);
                } else {
                    showResult('revokeResult', `<strong>Error:</strong> ${result.error}`, true);
                }
            } catch (error) {
                showResult('revokeResult', `<strong>Error:</strong> ${error.message}`, true);
            }
        });
    </script>
</body>
</html>
    ''')

@app.route('/api/generate-license', methods=['POST'])
def generate_license():
    """Generate a new signed license"""
    try:
        data = request.get_json()
        
        # Validate required fields
        required_fields = ['clientId', 'registeredFor']
        for field in required_fields:
            if not data.get(field):
                return jsonify({'error': f'Missing required field: {field}'}), 400
        
        # Calculate valid until timestamp
        valid_days = data.get('validDays', 365)
        valid_until = datetime.now() + timedelta(days=valid_days)
        valid_until_timestamp = int(valid_until.timestamp())
        
        # Create license
        license_data = {
            'licenseId': generate_license_id(),
            'clientId': data['clientId'],
            'extensions': data.get('extensions', []),
            'loginCount': data.get('loginCount', 100),
            'registeredFor': data['registeredFor'],
            'validUntil': valid_until_timestamp
        }
        
        # Sign the license
        signature = sign_license(license_data, private_key)
        
        # Create the new license string format
        license_base64 = create_license_string(license_data)
        
        # Create license file with both old and new formats for compatibility
        license_file = {
            'license': license_data,  # Keep old format for compatibility
            'licenseKey': license_base64,  # New base64 format
            'signature': signature
        }
        
        return jsonify(license_file)
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/validate-license', methods=['POST'])
def validate_license():
    """Validate a license signature and check if it's revoked"""
    try:
        data = request.get_json()
        
        # Check if license file structure is valid
        if 'license' not in data or 'signature' not in data:
            return jsonify({'valid': False, 'message': 'Invalid license file structure'})
        
        license_data = data['license']
        signature = data['signature']
        
        # Check if license is revoked
        revoked_licenses = load_revoked_licenses()
        if license_data.get('licenseId') in revoked_licenses:
            return jsonify({'valid': False, 'message': 'License has been revoked'})
        
        # Verify signature using new base64 format
        try:
            # Try new format first (base64 license key)
            if 'licenseKey' in data:
                license_base64 = data['licenseKey']
                message = license_base64.encode('utf-8')
            else:
                # Fall back to old format for backward compatibility
                license_json = json.dumps(license_data, separators=(',', ':'), sort_keys=True)
                message = license_json.encode('utf-8')
            
            digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
            digest.update(message)
            hashed = digest.finalize()
            
            signature_bytes = base64.urlsafe_b64decode(signature.encode('utf-8'))
            
            public_key.verify(
                signature_bytes,
                hashed,
                padding.PKCS1v15(),
                hashes.SHA256()
            )
            
            # Check if license is expired
            if license_data.get('validUntil', 0) < int(datetime.now().timestamp()):
                return jsonify({'valid': False, 'message': 'License has expired'})
            
            return jsonify({'valid': True, 'message': 'License is valid'})
            
        except Exception:
            return jsonify({'valid': False, 'message': 'Invalid license signature'})
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/revoke-license', methods=['POST'])
def revoke_license():
    """Revoke a license by adding it to the revocation list"""
    try:
        data = request.get_json()
        license_id = data.get('licenseId')
        
        if not license_id:
            return jsonify({'error': 'Missing licenseId'}), 400
        
        # Load current revoked licenses
        revoked_licenses = load_revoked_licenses()
        
        if license_id in revoked_licenses:
            return jsonify({'message': 'License already revoked'})
        
        # Add to revoked list
        revoked_licenses.append(license_id)
        save_revoked_licenses(revoked_licenses)
        
        return jsonify({'message': f'License {license_id} has been revoked'})
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/public-key', methods=['GET'])
def get_public_key():
    """Get the public key for the Go application"""
    with open(PUBLIC_KEY_PATH, 'r') as f:
        public_key_pem = f.read()
    
    return jsonify({'publicKey': public_key_pem})

@app.route('/api/revoked-licenses', methods=['GET'])
def get_revoked_licenses():
    """Get list of revoked license IDs"""
    return jsonify({'revokedLicenses': load_revoked_licenses()})

if __name__ == '__main__':
    print("Tallarin Activation Server")
    print("=" * 40)
    print(f"Private key: {PRIVATE_KEY_PATH}")
    print(f"Public key: {PUBLIC_KEY_PATH}")
    print("Starting server on http://localhost:5000")
    print("=" * 40)
    app.run(debug=True, host='0.0.0.0', port=5000)