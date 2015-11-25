package main

import (
	"log"
	"net/http"
	"time"

	"cubar.com/controller"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		var handler http.Handler
		handler = Logger(route.HandlerFunc, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Logger(inner http.Handler, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

var routes = Routes{
	Route{
		"Register",
		"Post",
		"/register",
		controller.Register,
	},
	Route{
		"Login",
		"POST",
		"/token-auth",
		controller.Login,
	},
	Route{
		"RefreshToken",
		"GET",
		"/refresh-token-auth",
		controller.RefreshToken,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		controller.Logout,
	},
	Route{
		"GetCommunity",
		"GET",
		"/community/{community_id}",
		controller.GetCommunity,
	},
	Route{
		"AddCommunity",
		"POST",
		"/community",
		controller.AddCommunity,
	},
	Route{
		"UpdateCommunity",
		"PUT",
		"/community/{community_id}",
		controller.UpdateCommunity,
	},
	Route{
		"DeleteCommunity",
		"DELETE",
		"/community/{community_id}",
		controller.DeleteCommunity,
	},
	Route{
		"QueryCommunities",
		"POST",
		"/communities",
		controller.QueryCommunities,
	},
	Route{
		"QueryUserCommunities",
		"POST",
		"/user/communities/{status}",
		controller.QueryUserCommunities,
	},
}
