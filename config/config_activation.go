package config

import (
	"context"
	"encoding/json"
	"fmt"

	"axia4/log"
	"axia4/types"

	keygen "github.com/keygen-sh/keygen-go/v3"
)

const (
	KeygenAccountID = "ec4dbc50-c77f-4b20-b974-c16d8c22bb87"
	KeygenProductID = "76559033-0137-4ee3-bb1d-4184b0564cf5"
)

func ActivateLicense() {
	if GetString("licenseFile") == "" {
		log.Info(
			log.ContextServer,
			"skipping activation check, no license installed",
		)

		SetLicense(types.License{})
		return
	}

	var licFile types.LicenseFile

	if err := json.Unmarshal(
		[]byte(GetString("licenseFile")),
		&licFile,
	); err != nil {
		log.Error(
			log.ContextServer,
			"could not unmarshal license from config",
			err,
		)

		return
	}

	if licFile.LicenseKey == "" {
		log.Error(
			log.ContextServer,
			"license key not found in license file",
			nil,
		)

		return
	}

	keygen.Account = KeygenAccountID
	keygen.Product = KeygenProductID
	keygen.LicenseKey = licFile.LicenseKey

	kgLicense, err := keygen.Validate(context.Background())
	if err != nil {
		log.Error(
			log.ContextServer,
			"license validation failed",
			err,
		)

		return
	}

	license := convertKeygenLicense(kgLicense)

	log.Info(
		log.ContextServer,
		fmt.Sprintf(
			"license %s validated successfully",
			license.LicenseId,
		),
	)

	SetLicense(license)
}

func convertKeygenLicense(
	kg *keygen.License,
) types.License {
	var validUntil int64

	if !kg.Expiry.IsZero() {
		validUntil = kg.Expiry.Unix()
	}

	var (
		clientID      string
		registeredFor string
		loginCount    int64
		extensions    []string
	)

	if kg.Metadata != nil {
		if v, ok := kg.Metadata["clientId"].(string); ok {
			clientID = v
		}

		if v, ok := kg.Metadata["registeredFor"].(string); ok {
			registeredFor = v
		}

		switch v := kg.Metadata["loginCount"].(type) {
		case int:
			loginCount = int64(v)
		case int64:
			loginCount = v
		case float64:
			loginCount = int64(v)
		}

		if raw, ok := kg.Metadata["extensions"].([]interface{}); ok {
			for _, ext := range raw {
				if s, ok := ext.(string); ok {
					extensions = append(extensions, s)
				}
			}
		}
	}

	return types.License{
		LicenseId:     kg.ID,
		ClientId:      clientID,
		Extensions:    extensions,
		LoginCount:    loginCount,
		RegisteredFor: registeredFor,
		ValidUntil:    validUntil,
	}
}
