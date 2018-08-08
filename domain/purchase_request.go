package domain

// PurchaseRequest is the structure to be used to receive the request from client
type PurchaseRequest struct {
	Amount float32 `json:"amount"`
}
