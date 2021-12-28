package routes

import (
	"drinkBack/database"
	"drinkBack/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
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
	// d.HandleFunc("/{id}", DeleteDrinkhandler).Methods("DELETE")
	p := router.PathPrefix("/user").Subrouter()
	p.HandleFunc("/", r.CreateUserHandler).Methods("POST")
	p.HandleFunc("/{id}", r.GetUserHandler).Methods("GET")
	p.HandleFunc("/", r.GetUsersHandler).Methods("GET")
	p.HandleFunc("/{id}/drinks", r.GetDrinksFromUserHandler).Methods("GET")
	// p.HandleFunc("/", UpdateUserHandler).Methods("PUT")
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

// func (r *Router) GetDrinkHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Get One Drink"))
// }
// func GetDrinksHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Get Multiple Drinks"))
// }
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

// func UpdateDrinkHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Update One Drinks"))
// }
// func DeleteDrinkhandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Delete One Drinks"))
// }
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

// func UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Create One Drink"))
// }
// func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Write([]byte("Create One Drink"))
// }
