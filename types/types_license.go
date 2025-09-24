package types

type License struct {
	LicenseId     string   `json:"licenseId"`
	ClientId      string   `json:"clientId"`
	Extensions    []string `json:"extensions"`
	LoginCount    int64    `json:"loginCount"`
	RegisteredFor string   `json:"registeredFor"`
	ValidUntil    int64    `json:"validUntil"`
}

type LicenseFile struct {
	License    License `json:"license"`
	LicenseKey string  `json:"licenseKey,omitempty"` // New base64 format
	Signature  string  `json:"signature"`
}
