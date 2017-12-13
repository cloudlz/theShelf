//Negroni is an idiomatic approach to web middleware but not a framework in Go.
// It is tiny, non-intrusive, and encourages use of net/http Handlers.
package main

import (
	"net/http"
	"log"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
	"strings"
)

//the middleware used to "/api..."
func MiddlewareApi(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//Gets the "appid" and checks if it exists
	appid := r.Header.Get("appid")
	if appid == "" {
		rw.Write([]byte("you need login in\n"))
		return
	}
	rw.Write([]byte("welcome:" + appid + "\n"))
	next(rw, r)
}

//the middleware used to "/web..."
func MiddlewareWeb(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Write([]byte("this is a web\n"))
	next(rw, r)
}

//the function will be run after web's middleware
func finalWeb(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing webFinalHandler")
	parameter := strings.Split(r.URL.Path, "/")[2]
	//if visit /web/user
	if parameter == "user" {
		w.Write([]byte("user list is : 454646464"))
		//if visit /web/admin
	} else if parameter == "admin" {
		w.Write([]byte("this is your world"))
	}

}

//the function will be run after api's middleware
func finalApi(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing apiFinalHandler")
	userid := mux.Vars(r)["userid"]
	adminid := mux.Vars(r)["adminid"]
	//according to the url parameter
	if userid != "" {
		w.Write([]byte("api userid:" + userid))
	}
	if adminid != "" {
		w.Write([]byte("api adminid:" + adminid))
	}
}

//it is the 404 handle function
func notFindHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get the 404"))
}

func main() {
	//main router
	router := mux.NewRouter()
	//get the 404 to the function notFindHandle
	router.NotFoundHandler = http.HandlerFunc(notFindHandle)

	//apiRouter under router
	apiRoutes := mux.NewRouter()
	//get the 404 to the function notFindHandle
	apiRoutes.NotFoundHandler = http.HandlerFunc(notFindHandle)
	// 在此新增API路由

	//webRouter under router
	webRoutes := mux.NewRouter()
	//get the 404 to the function notFindHandle
	webRoutes.NotFoundHandler = http.HandlerFunc(notFindHandle)
	// 在此新增Web路由

	// the middlewares use to router
	apiCommon := negroni.New(
		negroni.HandlerFunc(MiddlewareApi),
		negroni.NewLogger(),
	)
	webCommon := negroni.New(
		negroni.HandlerFunc(MiddlewareWeb),
		negroni.NewLogger(),
	)

	//the handle used to the url that begin with "/api"
	router.PathPrefix("/api").Handler(apiCommon.With(
		//connect the apiRouter
		negroni.Wrap(apiRoutes),
	))
	//the handle used to the url that begin with "/web"
	router.PathPrefix("/web").Handler(webCommon.With(
		//connect the webRouter
		negroni.Wrap(webRoutes),
	))
	apiRoutes.Handle("/api/user/{userid}", http.HandlerFunc(finalApi))
	apiRoutes.Handle("/api/admin/{adminid}", http.HandlerFunc(finalApi))
	webRoutes.Handle("/web/user", http.HandlerFunc(finalWeb))
	webRoutes.Handle("/web/admin", http.HandlerFunc(finalWeb))

	http.ListenAndServe(":8080", router)
	/*
		/api/
			welcome:[http.header.appid]
			you need login in.
		/api/user/:userid
			userid:123123123
		/api/admin/:adminid
			adminid:123123123
		/web/
			this is web
		/web/user/
			user list 123
		/web/admin/
			this is your world.
	 */
}
