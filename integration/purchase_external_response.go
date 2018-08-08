package integration

import (
	"regexp"
)

// ExternalResponse is a structure that represents the response from external service
type ExternalResponse struct {
	ID          string `json:"id,omitempty"`
	VoucherCode string `json:"voucher-code,omitempty"`
}

// IsIDValid validate if the ID is a valid UUID
func (externalResponse ExternalResponse) IsIDValid() bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(externalResponse.ID)
}
