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
	d.HandleFunc("/", r.CreateDrinkHandler).Methods("POST")
	// d.HandleFunc("/{id}", r.GetDrinkHandler).Methods("GET")
	// d.HandleFunc("/", GetDrinksHandler).Methods("GET")
	// d.HandleFunc("/", UpdateDrinkHandler).Methods("PUT")
	d.HandleFunc("/done", r.UpdateDrinkDoneHandler).Methods("PUT")
	// d.HandleFunc("/{id}", DeleteDrinkhandler).Methods("DELETE")
	p := router.PathPrefix("/user").Subrouter()
	p.HandleFunc("/", r.CreateUserHandler).Methods("POST")
	p.HandleFunc("/{id}", r.GetUserHandler).Methods("GET")
	p.HandleFunc("/", r.GetUsersHandler).Methods("GET")
	p.HandleFunc("/{id}/drinks", r.GetDrinksFromUserHandler).Methods("GET")
	p.HandleFunc("/{id}", r.UpdateUserHandler).Methods("PUT")
	// p.HandleFunc("/", DeleteUserHandler).Methods("DELETE")
	return router
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

type updateDrinkDoneRequest struct {
	Ids  []string `bson:"ids" json:"ids"`
	Done bool     `bson:"done" json:"done"`
}

func (r *Router) UpdateDrinkDoneHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}
	var bdJn updateDrinkDoneRequest
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	var objId_ids []primitive.ObjectID
	for _, id := range bdJn.Ids {
		if objId, err := primitive.ObjectIDFromHex(id); err != nil {
			http.Error(w, "Invalid id in request", http.StatusBadRequest)
			return
		} else {
			objId_ids = append(objId_ids, objId)
		}
	}
	drinks, err := r.Client.UpdateDrinksByIds(objId_ids, bdJn.Done)
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
	drinks, err := r.Client.FindDrinksByUser(usrId)
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
