package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/egamorim/purchase-api/domain"
)

var (
	externalURL = "http://www.mocky.io/v2/5b664ef8330000d718f6aac8"
	//WRONG : externalURL = "http://www.mocky.io/v2/5b2a61d93000000e009cd364"
	contextType = "application/json"
)

// Get execute a GET request to external service
func Get(request domain.PurchaseRequest) (ExternalResponse, error) {

	jsonValue, _ := json.Marshal(request)

	log.Println("Request: ", request)
	response, err := http.Post(externalURL, contextType, bytes.NewBuffer(jsonValue))

	externalResponse := ExternalResponse{}

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return ExternalResponse{}, err

	}

	decodeErr := json.NewDecoder(response.Body).Decode(&externalResponse)

	if decodeErr != nil {
		log.Printf("Error parsing response: %s\n", decodeErr)
		return ExternalResponse{}, decodeErr
	}

	if externalResponse.IsIDValid() {
		return externalResponse, nil
	} else {
		return ExternalResponse{}, errors.New("Response ID is not a valid UUID v4")
	}

}
