package itemController

import (
	"encoding/json"
	"fmt"
	"net/http"

	itemModel "test.com/events/model/itemModel"
)

func GetItems(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	items, err := itemModel.GetAllItem()

	if err != nil {
		http.Error(response, "Failed to fetch items from database", http.StatusInternalServerError)
	}

	json.NewEncoder(response).Encode(items)
}

func PostItem(response http.ResponseWriter, request *http.Request) {
	var newItem itemModel.Item
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	errJson := decoder.Decode(&newItem)

	if errJson != nil {
		http.Error(response, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err := itemModel.PostItem(newItem)

	if err != nil {
		fmt.Print("error")
	}
}
