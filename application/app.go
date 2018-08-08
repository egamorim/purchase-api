package application

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	govalidator "gopkg.in/thedevsaddam/govalidator.v1"

	"github.com/egamorim/purchase-api/domain"
	"github.com/egamorim/purchase-api/integration"
	"github.com/egamorim/purchase-api/tools"

	// to be used by postgresq connection
	_ "github.com/lib/pq"
)

//App represents our application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize the application and its dependencies
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err, "Database connection attempt failed")
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run - it executes the web server
func (a *App) Run(port string) {
	log.Printf("Server is running at port %s\n", port)
	http.ListenAndServe(port, a.Router)

}

func (a *App) getPurchases(w http.ResponseWriter, r *http.Request) {

	var purchase domain.Purchase

	purchases, err := purchase.GetAllPurchases(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, purchases)
}

func (a *App) newPurchase(w http.ResponseWriter, r *http.Request) {
	p, err := validateBody(r)
	if err != nil {
		log.Println(err.Error(), " - The request payload is not valid.")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := integration.Get(p)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	purchase := domain.Purchase{}
	purchase.ExternalID = data.ID
	purchase.VoucherCode = data.VoucherCode
	purchase.Amount = p.Amount

	log.Println(data)

	if err := purchase.SavePurchase(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := domain.PurchaseResponse{}
	response.VoucherCode = data.VoucherCode

	respondWithJSON(w, http.StatusCreated, response)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/purchase", a.newPurchase).Methods("POST")
	a.Router.HandleFunc("/purchase", a.getPurchases).Methods("GET")
	a.Router.Use(tools.LoggingRequest)
	a.Router.Use(tools.ValidadeRequestHeader)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func validateBody(r *http.Request) (domain.PurchaseRequest, error) {

	var p domain.PurchaseRequest
	rules := govalidator.MapData{
		"amount": []string{"required", "numeric_between:20.00,500.00"},
	}

	options := govalidator.Options{
		Request: r,
		Rules:   rules,
		Data:    &p,
	}
	validator := govalidator.New(options)

	e := validator.ValidateJSON()

	log.Println(p)
	log.Printf("The payload validation failed with error %s\n", e)

	if len(e) == 0 {
		return p, nil
	} else {
		jsonString, _ := json.Marshal(e)
		message := string(jsonString)
		return domain.PurchaseRequest{}, errors.New(message)
	}

}
