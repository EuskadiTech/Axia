# Tallarin Activation Server

A Python Flask-based activation server for generating, validating, and managing licenses for the Tallarin application.

## Features

- **License Generation**: Create signed licenses with customizable parameters
- **New License Format**: Uses base64-encoded license keys with semicolon-separated format
- **Backward Compatibility**: Maintains support for legacy JSON-based licenses
- **License Validation**: Verify license signatures and check revocation status
- **License Revocation**: Maintain a revocation list for invalidated licenses
- **Web Interface**: Simple HTML interface for license management
- **REST API**: RESTful endpoints for programmatic access

## License Format

### New Format (Default)
The server now generates licenses using a new base64-encoded format:

**Format**: `client_id;registered_for;login_count_limit;valid_for_days;ext1,ext2,ext3`

**Example**:
```
Original: "TEST_CLIENT_001;Test Company;100;365;premium,api_access"
Base64: "VEVTVF9DTElFTlRfMDAxO1Rlc3QgQ29tcGFueTsxMDA7MzY1O3ByZW1pdW0sYXBpX2FjY2Vzcw=="
```

### Backward Compatibility
The system maintains full backward compatibility with existing JSON-based licenses.

## Installation

1. **Install Python dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

2. **Run the activation server:**
   ```bash
   python app.py
   ```

3. **Access the web interface:**
   Open your browser to `http://localhost:5000`

## API Endpoints

### Generate License (New Format)
- **URL**: `POST /api/generate-license`
- **Content-Type**: `application/json`
- **Body**:
  ```json
  {
    "clientId": "CUSTOMER_001",
    "registeredFor": "Company Name",
    "loginCount": 100,
    "validDays": 365,
    "extensions": ["extension1", "extension2"]
  }
  ```
- **Response**: License file with signature and new `licenseKey` field
  ```json
  {
    "license": {
      "licenseId": "LI12345678",
      "clientId": "CUSTOMER_001",
      "extensions": ["extension1", "extension2"],
      "loginCount": 100,
      "registeredFor": "Company Name",
      "validUntil": 1735689600
    },
    "licenseKey": "Q1VTVE9NRVJfMDAxO0NvbXBhbnkgTmFtZTsxMDA7MzY1O2V4dGVuc2lvbjEsZXh0ZW5zaW9uMg==",
    "signature": "base64_encoded_signature_here"
  }
  ```

### Generate License (Old Format - for compatibility testing)
- **URL**: `POST /api/generate-license-old`
- **Content-Type**: `application/json`
- **Body**: Same as new format
- **Response**: License file without `licenseKey` field (JSON-based signature)

### Validate License
- **URL**: `POST /api/validate-license`
- **Content-Type**: `application/json`
- **Body**: Complete license file JSON (supports both old and new formats)
- **Response**: Validation result

### Revoke License
- **URL**: `POST /api/revoke-license`
- **Content-Type**: `application/json`
- **Body**:
  ```json
  {
    "licenseId": "LI12345678"
  }
  ```

### Get Public Key
- **URL**: `GET /api/public-key`
- **Response**: Public key in PEM format

### Get Revoked Licenses
- **URL**: `GET /api/revoked-licenses`
- **Response**: List of revoked license IDs

### List Licenses (New)
- **URL**: `GET /api/licenses`
- **Response**: List of all licenses saved in the folder
  ```json
  {
    "licenses": [
      {
        "licenseId": "LI12345678",
        "clientId": "CUSTOMER_001",
        "registeredFor": "Company Name",
        "validUntil": 1735689600,
        "filename": "LI12345678.json"
      }
    ]
  }
  ```

### Get Specific License (New)
- **URL**: `GET /api/licenses/{license_id}`
- **Response**: Complete license file for the specified license ID
  ```json
  {
    "license": { ... },
    "licenseKey": "base64_encoded_key",
    "signature": "base64_encoded_signature"
  }
  ```

## Signature Compatibility

**⚠️ Important Change**: The system now primarily uses base64-encoded license keys for signatures instead of JSON.

### New Format Signatures
The new format signs the base64-encoded license key directly:
- **Input**: `"CLIENT_001;Company Name;100;365;premium,api"`
- **Base64**: `"Q0xJRU5UXzAwMTtDb21wYW55IE5hbWU7MTAwOzM2NTtwcmVtaXVtLGFwaQ=="`
- **Signature**: Generated from the base64 string

### Backward Compatibility
For old format licenses, the system still uses Python's `json.dumps(separators=(',', ':'), sort_keys=True)`:

**Example JSON format**:
```json
{"clientId":"CUSTOMER_001","extensions":["premium","api"],"licenseId":"LI12345678","loginCount":100,"registeredFor":"Company Name","validUntil":1735689600}
```

The Go client automatically detects and handles both formats.

## Key Management

The server automatically generates RSA key pairs on first run:
- `private_key.pem`: Used for signing licenses (keep secure!)
- `public_key.pem`: Can be shared, used for verification

## Integration with Tallarin

1. **Update the public key in the Go application** by replacing the `publicKey` variable in `config/config_activation.go` with the generated public key from this server.

2. **Update the revocation list** by replacing the `revocations` slice in `config/config_activation.go` with license IDs from `/api/revoked-licenses`.

3. **Generate licenses** using the web interface or API and upload them to the Tallarin application.

## Security Notes

- Keep the `private_key.pem` file secure and never share it
- Consider implementing authentication for the API endpoints in production
- Use HTTPS in production environments
- Regularly backup the revocation list and keys

## Files Generated

- `private_key.pem`: Private RSA key (keep secure)
- `public_key.pem`: Public RSA key (for integration)
- `revoked_licenses.json`: List of revoked license IDs
- `licenses/`: Folder containing all generated license files

## Go Application Integration

### Using License ID (New Method)
Configure the Go application to fetch licenses directly from the activation server:
```json
{
  "licenseId": "LI12345678"
}
```

The Go application will automatically:
1. Fetch the license from the activation server's `/api/licenses/{license_id}` endpoint
2. Validate the license signature and check revocation status
3. Apply the license settings

### Using License File (Legacy Method)
Upload the complete license JSON file to the Go application configuration:
```json
{
  "licenseFile": "{ complete license JSON here }"
}
```

## Example License Files

### New Format (with licenseKey)
```json
{
  "license": {
    "licenseId": "LI8A3B9C5F",
    "clientId": "CUSTOMER_001", 
    "extensions": ["premium", "api"],
    "loginCount": 100,
    "registeredFor": "Company Name",
    "validUntil": 1735689600
  },
  "licenseKey": "Q1VTVE9NRVJfMDAxO0NvbXBhbnkgTmFtZTsxMDA7MzY1O3ByZW1pdW0sYXBp",
  "signature": "base64_encoded_signature_here"
}
```

### Old Format (backward compatibility)
```json
{
  "license": {
    "licenseId": "LI8A3B9C5F",
    "clientId": "CUSTOMER_001", 
    "extensions": ["premium", "api"],
    "loginCount": 100,
    "registeredFor": "Company Name",
    "validUntil": 1735689600
  },
  "signature": "base64_encoded_signature_here"
}
```

## Troubleshooting

1. **Key generation issues**: Ensure cryptography package is properly installed
2. **Signature verification failures**: Check that the public key in the Go app matches the one generated by this server
3. **Port conflicts**: Change the port in `app.py` if 5000 is already in use
4. **License format issues**: The Go client automatically detects old vs new format licenses