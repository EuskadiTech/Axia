package config

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"r3/log"
	"r3/types"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Activation server configuration
const activationServerURL = "https://ta-licensing.tech.eus"

// API response structures for activation server
type PublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

type RevokedLicensesResponse struct {
	RevokedLicenses []string `json:"revokedLicenses"`
}

// parseLicenseKey parses the new base64 license format: client_id;registered_for;login_count_limit;valid_for_days;ext1,ext2,ext3
func parseLicenseKey(licenseKey string) (*types.License, error) {
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(licenseKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 license key: %v", err)
	}
	
	// Split by semicolon
	parts := strings.Split(string(decoded), ";")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid license key format, expected 5 parts, got %d", len(parts))
	}
	
	clientId := parts[0]
	registeredFor := parts[1]
	
	// Parse login count
	loginCount, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid login count: %v", err)
	}
	
	// Parse valid for days and calculate validUntil timestamp
	validForDays, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid valid for days: %v", err)
	}
	validUntil := time.Now().Add(time.Duration(validForDays) * 24 * time.Hour).Unix()
	
	// Parse extensions
	var extensions []string
	if parts[4] != "" {
		extensions = strings.Split(parts[4], ",")
		// Trim whitespace from each extension
		for i, ext := range extensions {
			extensions[i] = strings.TrimSpace(ext)
		}
	}
	
	// Generate a license ID based on client ID for consistency
	licenseId := fmt.Sprintf("LK_%s_%d", clientId, time.Now().Unix())
	
	return &types.License{
		LicenseId:     licenseId,
		ClientId:      clientId,
		Extensions:    extensions,
		LoginCount:    loginCount,
		RegisteredFor: registeredFor,
		ValidUntil:    validUntil,
	}, nil
}

// marshalLicenseJSON marshals license data in the same format as Python's
// json.dumps(separators=(',', ':'), sort_keys=True) to ensure signature compatibility
func marshalLicenseJSON(license types.License) ([]byte, error) {
	// First marshal to get the data, then unmarshal to generic map
	data, err := json.Marshal(license)
	if err != nil {
		return nil, err
	}
	
	var jsonObj map[string]interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, err
	}
	
	// Get sorted keys
	keys := make([]string, 0, len(jsonObj))
	for k := range jsonObj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// Build JSON with sorted keys and compact format (no spaces after separators)
	var parts []string
	for _, k := range keys {
		v := jsonObj[k]
		vBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		parts = append(parts, fmt.Sprintf(`"%s":%s`, k, strings.TrimSpace(string(vBytes))))
	}
	
	return []byte("{" + strings.Join(parts, ",") + "}"), nil
}

// fetchPublicKey retrieves the public key from the activation server
func fetchPublicKey() (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(activationServerURL + "/api/public-key")
	if err != nil {
		return "", fmt.Errorf("failed to fetch public key: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	
	var keyResp PublicKeyResponse
	if err := json.Unmarshal(body, &keyResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	
	return keyResp.PublicKey, nil
}

// fetchRevokedLicenses retrieves the revoked licenses list from the activation server
func fetchRevokedLicenses() ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(activationServerURL + "/api/revoked-licenses")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch revoked licenses: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	var revokedResp RevokedLicensesResponse
	if err := json.Unmarshal(body, &revokedResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return revokedResp.RevokedLicenses, nil
}

func ActivateLicense() {
	if GetString("licenseFile") == "" {
		log.Info(log.ContextServer, "skipping activation check, no license installed")

		// set empty in case license was removed
		SetLicense(types.License{})
		return
	}

	var licFile types.LicenseFile

	if err := json.Unmarshal([]byte(GetString("licenseFile")), &licFile); err != nil {
		log.Error(log.ContextServer, "could not unmarshal license from config", err)
		return
	}

	var messageToVerify []byte
	var license types.License

	// Check if we have the new licenseKey format
	if licFile.LicenseKey != "" {
		// Use new base64 license key format for signature verification
		messageToVerify = []byte(licFile.LicenseKey)
		
		// Parse the license key to get license details
		parsedLicense, err := parseLicenseKey(licFile.LicenseKey)
		if err != nil {
			log.Error(log.ContextServer, "could not parse license key", err)
			return
		}
		license = *parsedLicense
	} else {
		// Fall back to old JSON format for backward compatibility
		licenseJson, err := marshalLicenseJSON(licFile.License)
		if err != nil {
			log.Error(log.ContextServer, "could not marshal license data", err)
			return
		}
		messageToVerify = licenseJson
		license = licFile.License
	}


	// verify signature
	data, _ := pem.Decode([]byte(publicKeyPEM))
	if data == nil {
		log.Error(log.ContextServer, "could not decode public key", errors.New(""))
		return
	}
	key, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		log.Error(log.ContextServer, "could not parse public key", errors.New(""))
		return
	}

	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature); err != nil {
		log.Error(log.ContextServer, "failed to verify license", err)
		return
	}

	// fetch revoked licenses from activation server
	log.Info(log.ContextServer, "checking license revocation status")
	revocations, err := fetchRevokedLicenses()
	if err != nil {
		log.Error(log.ContextServer, "failed to fetch revoked licenses from activation server", err)
		return
	}

	// check if license has been revoked
	if slices.Contains(revocations, license.LicenseId) {
		log.Error(log.ContextServer, "failed to enable license", fmt.Errorf("license ID '%s' has been revoked", license.LicenseId))
		return
	}

	// set license
	log.Info(log.ContextServer, "setting license")
	SetLicense(license)
}
