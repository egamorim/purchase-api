package domain

// PurchaseResponse is the structure to be used to respond to original client
type PurchaseResponse struct {
	VoucherCode string `json:"voucher-code,omitempty"`
}
