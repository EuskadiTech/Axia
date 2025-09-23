#!/usr/bin/env python3
"""
Test script for the Tallarin Activation Server
"""
import requests
import json

# Server URL
BASE_URL = "http://localhost:5000"

def test_generate_license():
    """Test license generation"""
    print("Testing license generation...")
    data = {
        "clientId": "TEST_CLIENT_001",
        "registeredFor": "Test Company Ltd.",
        "loginCount": 50,
        "validDays": 365,
        "extensions": ["premium", "api_access"]
    }
    
    response = requests.post(f"{BASE_URL}/api/generate-license", json=data)
    if response.status_code == 200:
        license_data = response.json()
        print("✓ License generated successfully!")
        print(f"  License ID: {license_data['license']['licenseId']}")
        print(f"  Valid Until: {license_data['license']['validUntil']}")
        return license_data
    else:
        print("✗ Failed to generate license")
        print(f"  Error: {response.text}")
        return None

def test_validate_license(license_data):
    """Test license validation"""
    print("\nTesting license validation...")
    response = requests.post(f"{BASE_URL}/api/validate-license", json=license_data)
    if response.status_code == 200:
        result = response.json()
        print(f"✓ License validation result: Valid={result['valid']}")
        print(f"  Message: {result['message']}")
        return result['valid']
    else:
        print("✗ Failed to validate license")
        print(f"  Error: {response.text}")
        return False

def test_revoke_license(license_id):
    """Test license revocation"""
    print("\nTesting license revocation...")
    data = {"licenseId": license_id}
    response = requests.post(f"{BASE_URL}/api/revoke-license", json=data)
    if response.status_code == 200:
        result = response.json()
        print(f"✓ License revoked: {result['message']}")
        return True
    else:
        print("✗ Failed to revoke license")
        print(f"  Error: {response.text}")
        return False

def test_get_public_key():
    """Test getting public key"""
    print("\nTesting public key retrieval...")
    response = requests.get(f"{BASE_URL}/api/public-key")
    if response.status_code == 200:
        result = response.json()
        print("✓ Public key retrieved successfully!")
        print("  Key preview:", result['publicKey'][:80] + "...")
        return result['publicKey']
    else:
        print("✗ Failed to get public key")
        print(f"  Error: {response.text}")
        return None

if __name__ == "__main__":
    print("Tallarin Activation Server Test")
    print("=" * 40)
    
    # Test 1: Generate a license
    license_data = test_generate_license()
    if not license_data:
        exit(1)
    
    license_id = license_data['license']['licenseId']
    
    # Test 2: Validate the generated license
    if not test_validate_license(license_data):
        exit(1)
    
    # Test 3: Get public key
    public_key = test_get_public_key()
    if not public_key:
        exit(1)
    
    # Test 4: Revoke the license
    if not test_revoke_license(license_id):
        exit(1)
    
    # Test 5: Validate revoked license (should fail)
    print("\nTesting validation of revoked license...")
    if test_validate_license(license_data):
        print("✗ Revoked license was incorrectly validated as valid")
    else:
        print("✓ Revoked license correctly rejected")
    
    print("\n" + "=" * 40)
    print("All tests completed successfully!")
    
    # Save a sample license file
    with open('sample_license.json', 'w') as f:
        json.dump(license_data, f, indent=2)
    print("Sample license saved to 'sample_license.json'")