package routes

import (
	"drinkBack/models"
	"drinkBack/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	if _, err := mail.ParseAddress(bdJn.Email); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = utils.ValidateUserPath(bdJn.Path)
	if err != nil {
		http.Error(w, "Invalid user path", http.StatusBadRequest)
		return
	}
	// hash password
	//regex = (?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?!.*?[=?<>()'"\/\&]).{8,20}
	ok, err := utils.ValidatePassword([]byte(bdJn.Password))
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	salt, err := utils.GenerateRandomSalt(10)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	hash, err := utils.HashPassword(bdJn.Password, salt)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}

	bdJn.Password = hash
	bdJn.Salt = salt

	inserted, err := r.Client.CreateNewUser(bdJn)
	if err != nil {
		http.Error(w, "Error inserting user to DB", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateAuthenticationToken(inserted.Id.Hex())
	if err != nil {
		// Deletar usuário criado se der merda no generate token
		http.Error(w, "Error generating authentication to user", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(models.LoginResponse{
		UserData: inserted,
		Token:    token,
	})

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

	if len(users) == 0 {
		res = []byte("[]")
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

	var bdJn models.User
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

func (r *Router) GetDrinksFromUserHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
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

	if len(drinks) == 0 {
		res = []byte("[]")
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

	if len(debts) == 0 {
		res = []byte("[]")
	}
	w.Write(res)
}
