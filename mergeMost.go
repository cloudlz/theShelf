//the Go file combined with global_error_handle,httpmock,log,middleware
//it used middleware and can handle error global output info by logrus
// and also can mock the http request
package main

import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
	"strings"
	"github.com/pkg/errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"io/ioutil"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the the severity or above that you choose .
	log.SetLevel(log.InfoLevel)
}

//it will mock all http request if head has mock also its attribute is true
func mockMiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	mockStr := r.Header.Get("mock")
	if mockStr == "true" {
		rw.Write([]byte("begin mock"))
		createMock() 
	} else {
		next(rw, r)
	}
}

//the middleware used to "/api..."
func apiMiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//Gets the "appid" and checks if it exists
	//if has "appid" return welcome "appid"
	//else return "you need login in"
	appid := r.Header.Get("appid")
	if appid == "" {
		rw.Write([]byte("you need login in\n"))
		return
	}
	rw.Write([]byte("welcome:" + appid + "\n"))
	next(rw, r)
}

//the middleware used to "/web..."
func webMiddleWare(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Write([]byte("this is a web\n"))
	next(rw, r)
}

//the function will be run after web's middleware
func webFinal(w http.ResponseWriter, r *http.Request) {
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
func apiFinal(w http.ResponseWriter, r *http.Request) {
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
func handleNotFind(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get the 404"))
}

//create a panic when http visit to it
func occurPanic(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("panicing from error"))
}

//It is a wrap function to catch panic in the  wrapped  function
func wrapRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				log.Warn("occur a panic")
				log.Errorf("location %+v", r)
				switch t := r.(type) {
				case string:
					//if it is string,
					//we will wrap it to a error to output
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				//if get the panic
				//we will output the error info
				//and return the http errorStatus
				log.Info("get the panic")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func createMock() {
	//Send the HTTP request to the target URL
	//It will print true response without mock
	resp, err := http.Get("http://localhost:8080/web/user")
	fmt.Println(resp, err)

	//Begin mock
	//It will intercept all http request
	//before execute httpmock.DeactivateAndReset
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//Register the url and method to be mocked
	httpmock.RegisterResponder("GET", "http://localhost:8080/web/user",
		httpmock.NewStringResponder(200, "I mocked it"))

	//Verification mock effect by request the url again
	mockResp, mockErr := http.Get("http://localhost:8080/web/user")
	body, _ := ioutil.ReadAll(mockResp.Body)
	fmt.Println(string(body), mockErr)
}

func main() {

	//main router
	router := mux.NewRouter()
	//get the 404 to the function notFindHandle
	router.NotFoundHandler = http.HandlerFunc(handleNotFind)

	//apiRouter under router
	apiRoutes := mux.NewRouter()
	//get the 404 to the function notFindHandle
	apiRoutes.NotFoundHandler = http.HandlerFunc(handleNotFind)
	// 在此新增API路由

	//webRouter under router
	webRoutes := mux.NewRouter()
	//get the 404 to the function notFindHandle
	webRoutes.NotFoundHandler = http.HandlerFunc(handleNotFind)
	// 在此新增Web路由

	// the middlewares use to router
	Common := negroni.New(
		negroni.HandlerFunc(mockMiddleWare),
	)

	//the handle used to the url that begin with "/api"
	router.PathPrefix("/api").Handler(Common.With(
		//connect the apiRouter
		negroni.HandlerFunc(apiMiddleWare),
		negroni.NewLogger(),
		negroni.Wrap(apiRoutes),
	))
	//the handle used to the url that begin with "/web"
	router.PathPrefix("/web").Handler(Common.With(
		//connect the webRouter
		negroni.HandlerFunc(webMiddleWare),
		negroni.NewLogger(),
		negroni.Wrap(webRoutes),
	))
	apiRoutes.Handle("/api/user/{userid}", wrapRecover(http.HandlerFunc(apiFinal)))
	apiRoutes.Handle("/api/panic", wrapRecover(http.HandlerFunc(occurPanic)))
	apiRoutes.Handle("/api/admin/{adminid}", wrapRecover(http.HandlerFunc(apiFinal)))
	webRoutes.Handle("/web/user", wrapRecover(http.HandlerFunc(webFinal)))
	webRoutes.Handle("/web/admin", wrapRecover(http.HandlerFunc(webFinal)))

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
