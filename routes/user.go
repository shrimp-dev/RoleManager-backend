package routes

import (
	"drinkBack/models"
	"drinkBack/utils"
	"encoding/json"
	"fmt"
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

	type request struct {
		Name      string             `bson:"name,omitempty" json:"name,omitempty"`
		Email     string             `bson:"email" json:"email"`
		Path      string             `bson:"path,omitempty" json:"path,omitempty"`
		Password  string             `bson:"password" json:"password"`
		CreatedBy primitive.ObjectID `bson:"createdby" json:"createdby"`
	}

	var bdJn request
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	if createdBy, err := primitive.ObjectIDFromHex(req.Context().Value("creator").(models.AccessTokenClaims).Id); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	} else {
		bdJn.CreatedBy = createdBy
	}

	if ok := utils.ValidateBody(bdJn); !ok {
		http.Error(w, "Missing information", http.StatusBadRequest)
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

	inserted, err := r.Client.CreateNewUser(models.User{
		Name:      bdJn.Name,
		Email:     bdJn.Email,
		Path:      bdJn.Path,
		Password:  hash,
		Salt:      salt,
		CreatedBy: bdJn.CreatedBy,
	})
	if err != nil {
		http.Error(w, "Error inserting user to DB", http.StatusBadRequest)
		return
	}

	token, _ := utils.GenerateAuthenticationToken(inserted.Id.Hex(), utils.AUTH)

	res, err := json.Marshal(models.LoginResponse{
		//Todo: Put Login Response
		Token: token,
	})

	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (r *Router) CreateUserInviteHandler(w http.ResponseWriter, req *http.Request) {
	a_token := req.Context().Value("usrToken").(models.AccessTokenClaims)
	createdBy := a_token.Id
	token, err := utils.GenerateAuthenticationToken(createdBy, utils.INVITE)
	if err != nil {
		http.Error(w, "Error authenticating your invite", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("http://localhost:3000/user/%s", token)))
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

	var bdJn models.UserUpdate
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
