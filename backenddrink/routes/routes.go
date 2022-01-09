package routes

import (
	"drinkBack/database"
	"drinkBack/middlewares"
	"drinkBack/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Client database.DbClient
}

func (r *Router) GenerateHandler() *mux.Router {
	router := mux.NewRouter()
	d := router.PathPrefix("/drinks").Subrouter()
	//Drinks
	d.HandleFunc("/", r.CreateDrinkHandler).Methods("POST", "OPTIONS")
	// d.HandleFunc("/{id}", r.GetDrinkHandler).Methods("GET")
	// d.HandleFunc("/", GetDrinksHandler).Methods("GET")
	// d.HandleFunc("/", UpdateDrinkHandler).Methods("PUT")
	d.HandleFunc("/done", r.UpdateDrinkDoneHandler).Methods("PUT", "OPTIONS")
	// d.HandleFunc("/{id}", DeleteDrinkhandler).Methods("DELETE")
	d.Use(middlewares.Authenticate())
	//Users
	p_open := router.PathPrefix("/user").Subrouter()
	p_open.HandleFunc("/", r.CreateUserHandler).Methods("POST", "OPTIONS")
	// For closed requests
	p := p_open.PathPrefix("").Subrouter()
	p.HandleFunc("/{id}", r.GetUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/", r.GetUsersHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}/drinks", r.GetDrinksFromUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}/debts", r.GetDebtsFromUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}", r.UpdateUserHandler).Methods("PUT", "OPTIONS")
	// p.HandleFunc("/", DeleteUserHandler).Methods("DELETE")
	p.Use(middlewares.Authenticate())
	//Debt
	de := router.PathPrefix("/debt").Subrouter()
	de.HandleFunc("/", r.CreateDebtHandler).Methods("POST", "OPTIONS")
	de.HandleFunc("/{id}", r.GetDebtHandler).Methods("GET", "OPTIONS")
	de.HandleFunc("/", r.GetDebtsHandler).Methods("GET", "OPTIONS")
	de.HandleFunc("/{id}/pay/{usrId}", r.PayDebtHandler).Methods("PUT", "OPTIONS")
	de.Use(middlewares.Authenticate())

	a := router.PathPrefix("/auth").Subrouter()
	a.HandleFunc("/", r.AuthenticateHandler).Methods("POST", "OPTIONS")

	router.Use(middlewares.Cors())

	return router
}

func (r *Router) AuthenticateHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn models.Request
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	var data models.LoginResponse
	ok, err := r.Client.VerifyUserPassword(bdJn.Email, bdJn.Password, &data)
	if err != nil {
		http.Error(w, "Error validating user password", http.StatusBadRequest)
		return
	}
	if !ok {
		http.Error(w, "Password in request does not match", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error converting data to response", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
