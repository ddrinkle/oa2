package handlers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/oa2/encoder"
	"github.com/ddrinkle/oa2/internal/models"
	"github.com/ddrinkle/oa2/internal/service"
	"github.com/ddrinkle/oa2/internal/token"
	"github.com/ddrinkle/platform/api"
	"github.com/ddrinkle/platform/query"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

func (h *App) LoginHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	p, err := h.Storage.GetAuthProvider(uuid.FromStringOrNil(vars["id"]))

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	app := r.Form.Get("app")
	returnURL := r.Form.Get("redirect_url")

	if app == "" || returnURL == "" {
		api.WriteError(w, errors.New("Missing App or RedirectURL"), http.StatusBadRequest)
		return
	}

	s := models.AuthState{
		AuthProviderID: p.UUID.String(),
		AppID:          app,
		Random:         rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
		ReturnURL:      returnURL,
	}

	oa2Factory := service.NewOA2Factory(h.Config.GetEnvironment().String())
	service, err := oa2Factory.GetServiceFromAuthProvider(p)

	if err != nil {
		api.WriteError(w, err, http.StatusBadRequest)
		return
	}

	stateEnc, err := s.Encode()
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	serviceCfg := service.GetOAuth2Config()
	serviceCfg.RedirectURL = h.App.Config.GetApp().BaseURL + "/oauth/token"

	url := serviceCfg.AuthCodeURL(stateEnc, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusFound)

}

// CodeHandler is used when an authorization request is completed by a third party, and they are sending us a code to get a token
func (h *App) CodeHandler(w http.ResponseWriter, r *http.Request) {
	//env := environment.GetInstance()

	err := r.ParseForm()
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	code := r.Form.Get("code")
	state := r.Form.Get("state")

	// Validate the State, which stores the AuthProvider of the original request
	s := models.AuthState{}
	err = s.Decode(state)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// Get the AuthProvider
	p, err := h.Storage.GetAuthProvider(uuid.FromStringOrNil(s.AuthProviderID))

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	oa2Factory := service.NewOA2Factory(h.Config.GetEnvironment().String())
	service, err := oa2Factory.GetServiceFromAuthProvider(p)

	parent := oauth2.NoContext
	ctx := context.WithValue(parent, oauth2.HTTPClient, http.Client{})
	service.SetContext(ctx)

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	// With the Code, get a Token Collection (access, and refresh tokens usually)
	tkn, err := service.RequestAccessTokenCollection(code)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	fmt.Println(tkn.Valid())
	fmt.Println(tkn.AccessToken)
	h.Env.Log.Warn("Got Token")

	response, err := service.RequestUserInfo(service.GetAuthorizationMethod())
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	authProviderUser, err := service.ParseUserInfoResponse(response)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{
		{
			Field: "a.uuid",
			Value: s.AppID,
		},
		{
			Field: "ap.uuid",
			Value: s.AuthProviderID,
		},
		{
			Field: "provider_user_id",
			Value: authProviderUser.Id,
		},
	})

	existingUsers, err := h.Storage.GetAppAuthProviderUsers(condition, 1, 1)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	app, err := uuid.FromString(s.AppID)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	if len(existingUsers) == 0 {

		authProvider, err := uuid.FromString(s.AuthProviderID)

		appUser := oa2.AppAuthProviderUser{
			UserID:       authProviderUser.Id,
			UserName:     authProviderUser.UserName,
			AccessToken:  tkn.AccessToken,
			RefreshToken: &tkn.RefreshToken,

			AuthProvider: &oa2.AuthProvider{
				UUID: authProvider,
			},
			App: &oa2.App{
				UUID: app,
			},
		}

		newUUID, err := h.Storage.AddAppAuthProviderUser(&appUser)
		if err != nil {
			api.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		existingUsers = append(existingUsers, oa2.AppAuthProviderUser{UUID: newUUID, AuthProvider: &oa2.AuthProvider{UUID: authProvider}})
	} else {
		for _, u := range existingUsers {
			u.AccessToken = tkn.AccessToken
			u.RefreshToken = &tkn.RefreshToken
			h.Storage.UpdateAppAuthProviderUser(u)
		}
	}

	encAuthUser, err := encoder.Encode(authProviderUser)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	collection, err := token.CreateJWT(h.Config.GetApp().AuthInfo.SigningKey, existingUsers[0].UUID, existingUsers[0].AuthProvider.UUID, app)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	if s.ReturnURL != "" {
		fmt.Println(tkn.AccessToken)
		v := url.Values{
			"user":     {encAuthUser},
			"ip_token": {tkn.AccessToken},
			"token":    {*collection.AccessToken},
			"refresh":  {*collection.RefreshToken},
		}
		if strings.Contains(s.ReturnURL, "?") {
			s.ReturnURL += "&"
		} else {
			s.ReturnURL += "?"
		}

		s.ReturnURL += v.Encode()
		h.Env.Log.Warn("encoding return url:" + s.ReturnURL)
		http.Redirect(w, r, s.ReturnURL, http.StatusFound)
	} else {
		res := oa2.TokenResponse{
			AccessToken:                 *collection.AccessToken,
			RefreshToken:                *collection.RefreshToken,
			UserInfo:                    encAuthUser,
			IdentityProviderAccessToken: tkn.AccessToken,
		}

		api.WriteJSON(w, res)
	}

}

func (h *App) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	refreshToken := r.Form.Get("refresh_token")
	tokenCollection, err := token.NewJWTCollection(h.Config.GetApp().AuthInfo.SigningKey, nil, &refreshToken)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	identityProvider, err := h.Storage.GetAuthProvider(tokenCollection.RefreshTokenClaims.IdentityProvider)

	if err != nil {
		fmt.Println("Can't find Identity Provider: " + tokenCollection.RefreshTokenClaims.IdentityProvider.String())
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	oa2Factory := service.NewOA2Factory(h.Config.GetEnvironment().String())
	service, err := oa2Factory.GetServiceFromAuthProvider(identityProvider)

	service.SetContext(r.Context())

	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	// With the Code, get a Token Collection (access, and refresh tokens usually)

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{
		{
			Field: "a.uuid",
			Value: tokenCollection.RefreshTokenClaims.App,
		},
		{
			Field: "ap.uuid",
			Value: tokenCollection.RefreshTokenClaims.IdentityProvider,
		},
		{
			Field: "aapu.uuid",
			Value: tokenCollection.RefreshTokenClaims.Subject,
		},
	})

	existingUsers, err := h.Storage.GetAppAuthProviderUsers(condition, 1, 1)
	if err != nil {
		fmt.Println("Can't find existing User: " + tokenCollection.RefreshTokenClaims.IdentityProvider.String())
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	if len(existingUsers) > 1 {
		api.WriteError(w, errors.New("Found too many users"), http.StatusInternalServerError)
		return
	}

	oAuthTkn := oauth2.Token{
		AccessToken:  existingUsers[0].AccessToken,
		RefreshToken: *existingUsers[0].RefreshToken,
	}
	newTokenCollection, err := service.RefreshAccessTokenInCollection(&oAuthTkn)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	existingUsers[0].AccessToken = newTokenCollection.AccessToken
	existingUsers[0].RefreshToken = &newTokenCollection.RefreshToken

	h.Storage.UpdateAppAuthProviderUser(existingUsers[0])

	encAuthUser, err := encoder.Encode(existingUsers[0])
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	collection, err := token.CreateJWT(h.Config.GetApp().AuthInfo.SigningKey, existingUsers[0].UUID, existingUsers[0].AuthProvider.UUID, tokenCollection.RefreshTokenClaims.App)
	if err != nil {
		api.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	res := oa2.TokenResponse{
		AccessToken:                 *collection.AccessToken,
		RefreshToken:                *collection.RefreshToken,
		UserInfo:                    encAuthUser,
		IdentityProviderAccessToken: newTokenCollection.AccessToken,
	}

	api.WriteJSON(w, res)
}
