package config

import "fmt"

type (
	TlsValidatorParams struct {
		ValidationCertPath string `json:"Validation_Certificates_Path"`
	}
)

func (rp TlsValidatorParams) Validate() error {
	if rp.ValidationCertPath == "" {
		return fmt.Errorf("path to trusted certificates cannot be empty, consider using path like: /usr/local/share/ca-certificates")
	}
	return nil
}
