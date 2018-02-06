package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/securecookie"
	"os"
	//"github.com/gorilla/sessions"
	//"github.com/nu7hatch/gouuid"
	"github.com/rs/cors"
	"encoding/json"
)
var db *sql.DB
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var r = httprouter.New()
var err error
type Configuration struct {
		HttpPort string
		ConnectionString string
		Appname string
		Runmode string
		Mysqluser string
		Mysqlpass string
		Mysqldb string
		SessionName string
}

func main() {

	configuration := Configuration{}
	filename := "app.json"
	log.Println(filename)
	pwd, _ := os.Getwd()

	//filename is the path to the json config file
	fto := pwd+"\\conf\\"+filename
	log.Println(fto)
	file, err := os.Open(fto)
	if err != nil {
		log.Println("File Open error")
		os.Exit(500) //return err
		}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("json erorr")
		os.Exit(500)//return err
		}

	// Create the database handle, confirm driver is present
	connectString := configuration.Mysqluser + ":" + configuration.Mysqlpass + configuration.Mysqldb
	log.Println(connectString)
	db, err = sql.Open("mysql", connectString  ) // "root:Bambie69@/test")
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	fmt.Println("db opened at root:****@/test")
	db.SetMaxIdleConns(100)
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}else {fmt.Println("verified db is open")}


	// Authentication and Authorization and Sessions
	r.GET("/", HomeHandler)
	r.GET("/getUser", GetUser)
	r.POST("/getUser", GetUser)
	r.GET("/index", indexPageHandler)
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)

	// products and articles
	r.GET("/internal", internalPageHandler)
	r.GET("/products", ProductsHandler)
	r.GET("/articles", SitesHandler)
	r.GET("/getUserName", validUser)
	// open for business
	fmt.Println("Router is open for business on port " + configuration.HttpPort)
	port := ":"+ configuration.HttpPort
	log.Println(port)
	handler := cors.Default().Handler(r)
	http.ListenAndServe(port, handler)
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		var msg string
		err := db.QueryRow("SELECT userName FROM user WHERE userName=? and userPwd=?", name, pass).Scan(&msg)
		if err != nil {
			log.Println(err)
			redirectTarget = "/index"
			http.Redirect(w, r, redirectTarget, 302)

		}else {
			setSession(name, w)
			redirectTarget = "/internal"
			http.Redirect(w, r, redirectTarget, 302)
		}
	}
}

func indexPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "")
}
func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hi")
}
func ProductsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
func SitesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
func internalPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userName := getUserName(r)

	//var cat = [12]map[string]string{}
	//var store = [6]map[string]string{}
	//var dim = make(map[int][]string)

	//var categories = ([]prod, []store,)
	//cat[0] = getCat("pCat1")
	//cat[1] = getCat("pCat2")
	//cat[2] = getCat("pCat3")
	//cat[3] = getCat("pCat4")
	//cat[4] = getCat("pCat5")
	//cat[5] = getCat("pCat6")
	//cat[6] = getCat("sCat1")
	//cat[7] = getCat("sCat2")
	//cat[8] = getCat("sCat3")
	//cat[9] = getCat("sCat4")
	//cat[10] = getCat("sCat5")
	//cat[11] = getCat("sCat6")

	if userName != "" {
		if err != nil{log.Fatalln(err)}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
// Code has not been tested.
