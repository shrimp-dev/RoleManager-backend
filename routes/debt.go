package routes

import (
	"drinkBack/models"
	"drinkBack/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Router) CreateDebtHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}

	var bdJn models.CreateDebtRequest
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	if ok := utils.ValidateBody(bdJn); !ok {
		http.Error(w, "Missing information", http.StatusBadRequest)
		return
	}

	var debtors []models.Debtor
	debt_remaining := bdJn.Amount
	for _, debtor := range bdJn.Debtors {
		if debtor.Amount == 0 {
			http.Error(w, "Invalid amount in debtor", http.StatusBadRequest)
			return
		}

		amount := debtor.Amount
		debt_remaining -= amount
		if debt_remaining < float32(0) {
			http.Error(w, "The value paid by all debtors is higher than the debt value", http.StatusBadRequest)
			return
		}

		debtors = append(debtors, models.Debtor{
			Id:     debtor.Id,
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
	creditor, _ := primitive.ObjectIDFromHex(req.Context().Value("usrToken").(models.AccessTokenClaims).Id)
	debt.Creditor = creditor
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
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}

	var bdJn models.PayDebtRequest
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	bdJn.Creditor = req.Context().Value("usrToken").(models.AccessTokenClaims).Id
	bdJn.Id = vars["id"]
	debt, err := r.Client.PayDebt(bdJn)
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

	if len(debt) == 0 {
		res = []byte("[]")
	}
	w.Write(res)
}
