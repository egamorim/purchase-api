package application_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/egamorim/purchase-api/integration"

	"github.com/egamorim/purchase-api/application"
)

var (
	port       = ":8000"
	dbUserName = "postgres"
	dbPassword = "123"
	dbName     = "postgres"
	a          = application.App{}
)

const tableCreationQuery = `
	create table IF NOT EXISTS purchase
	(id SERIAL,external_id TEXT NOT NULL, voucher_code TEXT NOT NULL, amount NUMERIC(10,2) 
	NOT NULL DEFAULT 0.00, CONSTRAINT p_pkey PRIMARY KEY (id))
`

func TestMain(m *testing.M) {

	a.Initialize(
		dbUserName,
		dbPassword,
		dbName)

	ensureTableExists()
	code := m.Run()
	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/purchase", nil)
	req.Header.Set("Bearer", "12345")
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreatePurchaseWithoutBearer(t *testing.T) {

	payload := []byte(`{"amount":201.90}`)
	req, _ := http.NewRequest("POST", "/purchase", bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestCreatePurchaseWithInvalidAmount(t *testing.T) {

	payload := []byte(`{"amount":"invalid"}`)
	req, _ := http.NewRequest("POST", "/purchase", bytes.NewBuffer(payload))
	req.Header.Set("Bearer", "12345")

	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreatePurchaseWithAmountLess20(t *testing.T) {

	payload := []byte(`{"amount":19.99}`)
	req, _ := http.NewRequest("POST", "/purchase", bytes.NewBuffer(payload))
	req.Header.Set("Bearer", "12345")
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePurchaseWithValidAmount(t *testing.T) {

	payload := []byte(`{"amount":210.54}`)
	req, _ := http.NewRequest("POST", "/purchase", bytes.NewBuffer(payload))
	req.Header.Set("Bearer", "12345")

	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestUUIDValidate(t *testing.T) {

	invalidResponse := integration.ExternalResponse{}
	invalidResponse.ID = "INVALID_ID001"
	isValid := invalidResponse.IsIDValid()

	if isValid {
		t.Error(invalidResponse.ID, " should be considered invalid. Validation failed")
	}

	validResponse := integration.ExternalResponse{}
	validResponse.ID = "6af5fa8f-e3e9-46d5-ad3c-239a36f4a395"

	if !validResponse.IsIDValid() {
		t.Error(validResponse.ID, " should be considered valid. Validation failed")
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM purchase")
	a.DB.Exec("ALTER SEQUENCE purchase_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
