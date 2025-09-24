package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"r3/log"
	"r3/types"
	"sort"
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

// fetchLicenseFromServer retrieves a specific license from the activation server's folder
func fetchLicenseFromServer(licenseId string) (*types.License, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(activationServerURL + "/api/licenses/" + licenseId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch license: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	// The server returns a structure with "license" field containing the license data
	var response struct {
		License types.License `json:"license"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse license response: %v", err)
	}
	
	return &response.License, nil

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

	var license types.License

	// Extract license ID from the license file to fetch from server
	licenseId := licFile.LicenseId
	if licenseId == "" {
		log.Error(log.ContextServer, "license ID not found in license file", nil)
		return
	}

	// fetch <licenseserver>/api/licenses/<LicenseId> and set license var on it's json output
	log.Info(log.ContextServer, fmt.Sprintf("fetching license %s from activation server", licenseId))
	fetchedLicense, err := fetchLicenseFromServer(licenseId)
	if err != nil {
		log.Error(log.ContextServer, "failed to fetch license from server", err)
		return
	}
	
	// Set license from fetched data
	license = *fetchedLicense

	// fetch revoked licenses from activation server
	// set license
	log.Info(log.ContextServer, "setting license")
	SetLicense(license)
}
