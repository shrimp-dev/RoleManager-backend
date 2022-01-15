package routes

import (
	"drinkBack/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (r *Router) CreateDrinkHandler(w http.ResponseWriter, req *http.Request) {
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

	if bdJn.Name == "" {
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

	var bdJn models.Request
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	drinks, err := r.Client.UpdateDrinksByIds(bdJn.Ids, bdJn.Done)
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
