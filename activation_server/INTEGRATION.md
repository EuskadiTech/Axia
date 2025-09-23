# Integration Guide: Tallarin Activation Server

This guide explains how to integrate the Python Flask activation server with the main Tallarin Go application.

## Overview

The activation server generates RSA-signed licenses that the main application validates. Integration requires updating the Go application with the public key and revocation list from the activation server.

## Step 1: Start the Activation Server

1. **Navigate to the activation server directory:**
   ```bash
   cd activation_server/
   ```

2. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

3. **Start the server:**
   ```bash
   python app.py
   ```

4. **Access the web interface:**
   Open http://localhost:5000 in your browser

## Step 2: Get the Public Key

1. **Via Web Interface:**
   - The public key will be displayed in the server startup logs
   - Or visit http://localhost:5000/api/public-key

2. **Via API:**
   ```bash
   curl http://localhost:5000/api/public-key
   ```

3. **From File:**
   ```bash
   cat activation_server/public_key.pem
   ```

## Step 3: Update the Go Application

### Update Public Key

Edit `config/config_activation.go` and replace the `publicKey` variable with your new public key:

```go
// Replace this variable with the public key from your activation server
var publicKey = `-----BEGIN RSA PUBLIC KEY-----
YOUR_NEW_PUBLIC_KEY_HERE
-----END RSA PUBLIC KEY-----`
```

### Update Revocation List

Replace the `revocations` slice with current revoked licenses:

```go
// Get revoked licenses from: http://localhost:5000/api/revoked-licenses
var revocations = []string{"LI00334231", "LI58AF82DC"} // Add actual revoked license IDs
```

## Step 4: Build and Test

1. **Rebuild the Go application:**
   ```bash
   go build -o r3
   ```

2. **Generate a test license:**
   - Use the web interface at http://localhost:5000
   - Or via API:
     ```bash
     curl -X POST -H "Content-Type: application/json" \
          -d '{"clientId":"TEST_001","registeredFor":"Test Company","validDays":30}' \
          http://localhost:5000/api/generate-license
     ```

3. **Test license upload:**
   - Save the license JSON to a file (e.g., `test_license.json`)
   - Upload it through the Tallarin admin interface
   - The application should accept and activate the license

## Security Considerations

1. **Secure the Private Key:**
   - Store `private_key.pem` securely
   - Set proper file permissions (600)
   - Consider using environment variables or key management systems

2. **Enable HTTPS:**
   - Use a production WSGI server (e.g., Gunicorn, uWSGI)
   - Configure SSL/TLS certificates
   - Use a reverse proxy (nginx, Apache)

3. **Add Authentication:**
   - Implement API authentication for production
   - Consider rate limiting
   - Add access logging

## License Workflow

1. **Generate License:** Create license via web interface or API
2. **Distribute License:** Send license file to customer
3. **Install License:** Customer uploads license file to Tallarin admin interface
4. **Validation:** Tallarin validates signature and checks revocation status
5. **Activation:** Valid license activates features in the application
6. **Revocation:** If needed, revoke license via activation server

## Troubleshooting

### Common Issues

1. **Signature Verification Failed:**
   - Ensure public key in Go app matches activation server
   - Check JSON formatting consistency
   - Verify RSA key format (PKCS1 for public key)

2. **License Rejected:**
   - Check if license is expired (`validUntil` timestamp)
   - Verify license ID is not in revocation list
   - Confirm license file structure is correct

3. **Server Connection Issues:**
   - Verify activation server is running
   - Check firewall settings
   - Confirm port availability