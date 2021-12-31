package routes

import (
	"drinkBack/database"
	"drinkBack/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	//Users
	p := router.PathPrefix("/user").Subrouter()
	p.HandleFunc("/", r.CreateUserHandler).Methods("POST", "OPTIONS")
	p.HandleFunc("/{id}", r.GetUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/", r.GetUsersHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}/drinks", r.GetDrinksFromUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}/debts", r.GetDebtsFromUserHandler).Methods("GET", "OPTIONS")
	p.HandleFunc("/{id}", r.UpdateUserHandler).Methods("PUT", "OPTIONS")
	// p.HandleFunc("/", DeleteUserHandler).Methods("DELETE")
	//Debt
	de := router.PathPrefix("/debt").Subrouter()
	de.HandleFunc("/", r.CreateDebtHandler).Methods("POST", "OPTIONS")
	de.HandleFunc("/{id}", r.GetDebtHandler).Methods("GET", "OPTIONS")
	de.HandleFunc("/", r.GetDebtsHandler).Methods("GET", "OPTIONS")
	de.HandleFunc("/{id}/pay/{usrId}", r.PayDebtHandler).Methods("PUT", "OPTIONS")

	router.Use(cors())

	return router
}

func cors() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			if req.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func (r *Router) CreateDrinkHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn models.Drink
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	error_id, _ := primitive.ObjectIDFromHex("000000000000000000000000")

	if (bdJn.UsrId == error_id) || (bdJn.Name == "") {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	inserted, err := r.Client.CreateNewDrink(models.Drink{
		UsrId: bdJn.UsrId,
		Name:  bdJn.Name,
		Done:  false,
	})
	if err != nil {
		http.Error(w, "Error inserting drink to DB", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(inserted)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
func (r *Router) UpdateDrinkDoneHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn bson.M
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	var objId_ids []primitive.ObjectID
	for _, id := range bdJn["ids"].([]interface{}) {
		switch id.(type) {
		case string:
		default:
			http.Error(w, "Invalid ID in ids", http.StatusBadRequest)
			return
		}
		if objId, err := primitive.ObjectIDFromHex(id.(string)); err != nil {
			http.Error(w, "Invalid id in request", http.StatusBadRequest)
			return
		} else {
			objId_ids = append(objId_ids, objId)
		}
	}
	switch bdJn["done"].(type) {
	case bool:
	default:
		http.Error(w, "Invalid Done in body", http.StatusBadRequest)
		return
	}
	drinks, err := r.Client.UpdateDrinksByIds(objId_ids, bdJn["done"].(bool))
	if err != nil {
		http.Error(w, "Error while updating drinks", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(drinks)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) GetDrinksFromUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id in request", http.StatusBadRequest)
		return
	}
	usrId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "The id in the request is invalid", http.StatusBadRequest)
		return
	}
	drinks, err := r.Client.FindDrinksOfUser(usrId)
	if err != nil {
		http.Error(w, "Error gathering drinks from DB", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(drinks)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn models.User
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	if bdJn.Name == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	inserted, err := r.Client.CreateNewUser(models.User{
		Name: bdJn.Name,
		Path: bdJn.Path,
	})
	if err != nil {
		http.Error(w, "Error inserting user to DB", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(inserted)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
func (r *Router) GetUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id in request", http.StatusBadRequest)
		return
	}
	usrId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "The id in the request is invalid", http.StatusBadRequest)
		return
	}
	user, err := r.Client.FindUserById(usrId)
	if err != nil {
		http.Error(w, "Error getting user from DB", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) GetUsersHandler(w http.ResponseWriter, req *http.Request) {
	users, err := r.Client.FindAllUsers()
	if err != nil {
		http.Error(w, "Error getting users from DB", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	usrId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id sent in request", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn bson.M
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	users, err := r.Client.UpdateUserById(usrId, bdJn)
	if err != nil {
		http.Error(w, "Invalid update query", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) CreateDebtHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn bson.M
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	if bdJn["amount"] == nil {
		http.Error(w, "Missing amount in debtor", http.StatusBadRequest)
		return
	}
	switch bdJn["amount"].(type) {
	case float64:
	default:
		http.Error(w, "Invalid amount in debt", http.StatusBadRequest)
		return
	}
	var debtors []models.Debtor
	debt_remaining := float32(bdJn["amount"].(float64))
	for _, debtor_orig := range bdJn["debtors"].([]interface{}) {
		debtor := debtor_orig.(map[string]interface{})
		id, err := primitive.ObjectIDFromHex(debtor["_id"].(string))
		if err != nil {
			http.Error(w, "Invalid id in debtor", http.StatusBadRequest)
			return
		}
		if debtor["amount"] == nil {
			http.Error(w, "Invalid amount in debtor", http.StatusBadRequest)
			return
		}
		amount := float32(debtor["amount"].(float64))
		debt_remaining -= amount
		if debt_remaining < float32(0) {
			http.Error(w, "The value paid by all debtors is higher than the debt value", http.StatusBadRequest)
			return
		}
		debtors = append(debtors, models.Debtor{
			Id:     id,
			Amount: amount,
		})
	}
	if debt_remaining > float32(0.1) {
		http.Error(w, "The value paid by all debtors is lower than the debt value", http.StatusBadRequest)
		return
	}
	var debt models.Debt
	if err := json.Unmarshal(body, &debt); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	debt.Debtors = debtors
	debt, err = r.Client.CreateNewDebt(debt)
	if err != nil {
		http.Error(w, "Error creating new debt", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(debt)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
func (r *Router) PayDebtHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	str_id := vars["id"]
	if str_id == "" {
		http.Error(w, "Missing id in request", http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(str_id)
	if err != nil {
		http.Error(w, "Invalid id in request", http.StatusBadRequest)
		return
	}
	str_usrId := vars["usrId"]
	if str_usrId == "" {
		http.Error(w, "Missing usrId in request", http.StatusBadRequest)
		return
	}
	usrId, err := primitive.ObjectIDFromHex(str_usrId)
	if err != nil {
		http.Error(w, "Invalid usrId in request", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn map[string]bool
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	debt, err := r.Client.PayDebt(bson.M{
		"_id":   id,
		"usrId": usrId,
		"paid":  bdJn["paid"],
	})
	if err != nil {
		http.Error(w, "Error updating debt", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(debt)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) GetDebtHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	str_id := vars["id"]
	id, err := primitive.ObjectIDFromHex(str_id)
	if err != nil {
		http.Error(w, "Invalid id in request", http.StatusBadRequest)
		return
	}
	debt, err := r.Client.FindDebtById(id)
	if err != nil {
		http.Error(w, "Error finding the debt", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(debt)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) GetDebtsHandler(w http.ResponseWriter, req *http.Request) {
	debt, err := r.Client.FindAllDebts()
	if err != nil {
		http.Error(w, "Error finding debts", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(debt)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
func (r *Router) GetDebtsFromUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	str_usrId := vars["id"]
	usrId, err := primitive.ObjectIDFromHex(str_usrId)
	if err != nil {
		http.Error(w, "Invalid usrId in request", http.StatusBadRequest)
		return
	}
	debts, err := r.Client.FindDebtsOfUser(usrId)
	if err != nil {
		http.Error(w, "Error finding debts", http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(debts)
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
