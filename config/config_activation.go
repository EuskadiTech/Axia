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

	licenseJson, err := json.Marshal(licFile.License)
	if err != nil {
		log.Error(log.ContextServer, "could not marshal license data", err)
		return
	}
	hashed := sha256.Sum256(licenseJson)

	// get license signature
	signature, err := base64.URLEncoding.DecodeString(licFile.Signature)
	if err != nil {
		log.Error(log.ContextServer, "could not decode license signature", err)
		return
	}

	// fetch public key from activation server
	log.Info(log.ContextServer, "fetching public key from activation server")
	publicKeyPEM, err := fetchPublicKey()
	if err != nil {
		log.Error(log.ContextServer, "failed to fetch public key from activation server", err)
		return
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
	if slices.Contains(revocations, licFile.License.LicenseId) {
		log.Error(log.ContextServer, "failed to enable license", fmt.Errorf("license ID '%s' has been revoked", licFile.License.LicenseId))
		return
	}

	// set license
	log.Info(log.ContextServer, "setting license")
	SetLicense(licFile.License)
}
