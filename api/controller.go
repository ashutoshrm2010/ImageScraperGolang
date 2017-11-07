package api

import (
	"github.com/sourcecode/ImageScrapGolang/system"
	"net/http"
	"encoding/json"
	"github.com/sourcecode/ImageScrapGolang/service"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) ImageScrapfromGoogle(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	response, err :=service.ImageScrapfromGoogleService(data["userInput"])

	if err != nil {
		return nil, err
	}
	return response, err
}

func (controller *Controller) ListUserSearchInputs(w http.ResponseWriter, r *http.Request) ([]byte, error) {

	response, err :=service.ListUserSearchInputsService()

	if err != nil {
		return nil, err
	}
	return response, err
}

func (controller *Controller) GetSearchedImageUrlsFromDB(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	response, err :=service.GetSearchedImageUrlsFromDBService(data["searchkeyId"])

	if err != nil {
		return nil, err
	}
	return response, err
}