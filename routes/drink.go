package routes

import (
	"drinkBack/models"
	"drinkBack/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Router) CreateDrinkHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "An error occurred while trying to read the body", http.StatusBadRequest)
		return
	}

	var bdJn models.CreateDrinkRequest
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}
	if ok := utils.ValidateBody(bdJn); !ok {
		http.Error(w, "Information missing in body", http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(bdJn.UsrId)
	if err != nil {
		http.Error(w, "Invalid information in body", http.StatusBadRequest)
		return
	}

	inserted, err := r.Client.CreateNewDrink(models.Drink{
		UsrId: id,
		Name:  bdJn.Name,
		Done:  false,
	})
	if err != nil {
		http.Error(w, "Error inserting drink to DB", http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(models.CreateDrinkResponse{
		Inserted: inserted,
	})
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

	var bdJn models.UpdateDrinkDoneRequest
	if err := json.Unmarshal(body, &bdJn); err != nil {
		http.Error(w, "Invalid JSON sent in body", http.StatusBadRequest)
		return
	}

	var ids []primitive.ObjectID
	for _, id := range bdJn.Ids {
		if parsed, err := primitive.ObjectIDFromHex(id); err != nil {
			http.Error(w, "Invalid information in body", http.StatusBadRequest)
			return
		} else {
			ids = append(ids, parsed)
		}
	}

	updated, err := r.Client.UpdateDrinksByIds(ids, bdJn.Done)
	if err != nil {
		http.Error(w, "Error while updating drinks", http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(models.UpdateDrinkDoneResponse{
		Updated: updated,
	})
	if err != nil {
		http.Error(w, "Error converting data to send back", http.StatusInternalServerError)
		return
	}
	w.Write(res)
}
