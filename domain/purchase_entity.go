package domain

import (
	"database/sql"
)

// Purchase is the structure that represents the Purchase data table on our code
type Purchase struct {
	ID          string  `json:"id,omitempty"`
	ExternalID  string  `json:"external_id,omitempty"`
	VoucherCode string  `json:"voucher-code,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
}

// SavePurchase should be used to save a Purchase on database
func (p *Purchase) SavePurchase(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO purchase(external_id, voucher_code, amount) VALUES($1, $2, $3) RETURNING id",
		p.ExternalID, p.VoucherCode, p.Amount).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetAllPurchases get all records on Purchase data table and return a parsed list
func (p *Purchase) GetAllPurchases(db *sql.DB) ([]Purchase, error) {
	rows, err := db.Query("SELECT id, external_id, voucher_code, amount FROM purchase")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	purchases := []Purchase{}

	for rows.Next() {
		var p Purchase
		if err := rows.Scan(&p.ID, &p.ExternalID, &p.VoucherCode, &p.Amount); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}

	return purchases, nil
}
