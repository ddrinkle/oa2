package handlers

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/platform/api"
	"github.com/ddrinkle/platform/query"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
	uuid "github.com/satori/go.uuid"
)

func (h *App) GetAuthProvidersHandler(w http.ResponseWriter, r *http.Request) {

	a, err := h.Storage.GetAuthProviders(query.SQLCondition{}, 0, 0)

	if err != nil {
		h.Env.Log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	api.WriteJSONAPI(w, a, http.StatusOK)

}
func (h *App) GetAuthProviderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a, err := h.Storage.GetAuthProvider(uuid.FromStringOrNil(vars["id"]))

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	api.WriteJSONAPI(w, a, http.StatusOK)

}

func (h *App) CreateAuthProviderHandler(w http.ResponseWriter, r *http.Request) {

	provider := oa2.AuthProvider{}
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	err = jsonapi.Unmarshal(jsonBytes, &provider)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	providerID, err := h.Storage.AddAuthProvider(&provider)

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	a, err := h.Storage.GetAuthProvider(providerID)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	api.WriteJSONAPI(w, a, http.StatusOK)
}
