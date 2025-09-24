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

    // FIXME: fetch <licenseserver>/api/licenses/<LicenseId> and set license var on it's json output.

	// fetch revoked licenses from activation server
	// set license
	log.Info(log.ContextServer, "setting license")
	SetLicense(license)
}
