package routes

import (
	"github.com/ddrinkle/oa2/internal/handlers"
	"github.com/ddrinkle/platform/routing"
)

//GetRoutes returns the available Routes
func GetRoutes(app handlers.App) routing.Routes {
	return routing.Routes{
		//OAuth2 Endpoints
		routing.Route{
			Name:          "Auth",
			Method:        "GET",
			Pattern:       "/oauth/login/{id:.+}",
			Authenticated: false,
			HandlerFunc:   app.LoginHandler,
		},
		routing.Route{
			Name:          "Auth-Token",
			Method:        "GET",
			Pattern:       "/oauth/token",
			Authenticated: false,
			HandlerFunc:   app.CodeHandler,
		},
		routing.Route{
			Name:          "Auth-Token",
			Method:        "GET",
			Pattern:       "/oauth/refresh",
			Authenticated: false,
			HandlerFunc:   app.RefreshHandler,
		},
		//App Endpoints
		routing.Route{
			Name:          "App",
			Method:        "POST",
			Pattern:       "/app",
			Authenticated: true,
			HandlerFunc:   app.CreateAppHandler,
		},
		routing.Route{
			Name:          "App",
			Method:        "GET",
			Pattern:       "/app/{id:.+}",
			Authenticated: true,
			HandlerFunc:   app.GetAppHandler,
		},
		// AuthProvider Endpoints
		routing.Route{
			Name:          "AuthProvider",
			Method:        "POST",
			Pattern:       "/auth_provider",
			Authenticated: true,
			HandlerFunc:   app.CreateAuthProviderHandler,
		},
		routing.Route{
			Name:          "GetAuthProvider",
			Method:        "GET",
			Pattern:       "/auth_provider/{id:.+}",
			Authenticated: true,
			HandlerFunc:   app.GetAuthProviderHandler,
		},
		routing.Route{
			Name:          "GetAuthProviders",
			Method:        "GET",
			Pattern:       "/auth_provider",
			Authenticated: true,
			HandlerFunc:   app.GetAuthProvidersHandler,
		},
	}
}
