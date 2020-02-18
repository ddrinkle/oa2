package handlers

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/oa2/internal/storage"
	"github.com/ddrinkle/platform/api"
	"github.com/ddrinkle/platform/app"
	"github.com/ddrinkle/platform/crypt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
	uuid "github.com/satori/go.uuid"
)

type App struct {
	app.App
	Storage storage.Crate
}

func (h *App) GetAppHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a, err := h.Storage.GetApp(uuid.FromStringOrNil(vars["id"]))

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSONAPI(w, &a, http.StatusOK)

}

func (h *App) CreateAppHandler(w http.ResponseWriter, r *http.Request) {

	app := oa2.App{
		AppKey:    &crypt.String{},
		SecretKey: &crypt.String{},
	}

	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	err = jsonapi.Unmarshal(jsonBytes, &app)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	appKey, err := crypt.GenerateRandomString(32)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	secretKey, err := crypt.GenerateRandomString(32)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	app.AppKey.Set(appKey)
	app.SecretKey.Set(secretKey)

	appID, err := h.Storage.AddApp(&app)

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	newApp, err := h.Storage.GetApp(appID)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSONAPI(w, newApp, http.StatusCreated)

}
