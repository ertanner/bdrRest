package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"log"
	"strconv"
	"github.com/gorilla/sessions"

	//"github.com/satori/go.uuid"
)
var store = sessions.NewCookieStore([]byte("radar-super-secret"))

func GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println(body)
	log.Println(r)
	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	name :=  data["username"]
	pass := data["password"]
	email := data["accountname"]
	fmt.Println(name + " " + pass + " " + email)
	if name != "" && pass != "" {
		//var cred bool
		// check to be sure the user is in the db
		cred := getPwd(name, pass)
		log.Println("credentials: " + strconv.FormatBool(cred))
		if !cred {
			http.Error(w,"Invalid User", 418)
		}
		//get pwd hash and salt
		bpass := []byte(pass)
		hash := hashAndSalt(bpass)
		log.Println(hash)

		pwdMatch := comparePasswords(hash, bpass)
		if !pwdMatch {
			http.Error(w,"Invalid Password", 418)
		} else {
			//it passed the check.  Return an ok
			// log it and set up a session
			log.Println("pwdMatch: " + strconv.FormatBool(pwdMatch))
			// w.Header().Set("Access-Control-Allow-Origin", "*")
			//if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			//}
			setSession(email, w)
		}

	}
}
func SetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println(body)
	log.Println(r)
	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	name :=  data["username"]
	pass := data["password"]
	email := data["accountname"]
	fmt.Println(name + " " + pass + " " + email)

	//check if valid user
	acctChck := checkAccountName(email)
	log.Println("credentials: " + strconv.FormatBool(acctChck))
	if acctChck {
		log.Println("Account already used")
		http.
		http.Error(w,"Account already used", 418)
	}

	//return if invalid
	//write user to db
	// call hashAndSaltpwd
	//call setSession
	//return sessionCookie to user
}
func ChagePwd(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

}

func setSession(email string, response http.ResponseWriter) {
	value := map[string]string{
		"name": email,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}

		http.SetCookie(response, cookie)
	}
}

func validUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	body, err := ioutil.ReadAll(r.Body)
	log.Println(body)
	if err != nil {
		panic(err.Error())
	}
}
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
func submit(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	fmt.Println(r)

	//seasonality()
}

func getPwd(name string, pass string) bool {
	// Prompt the user to enter a password
	var ok bool = false
	var msg string
	err := db.QueryRow("SELECT userName FROM user WHERE userName=? and userPwd=?", name, pass).Scan(&msg)
	if err != nil {
		log.Println(err)
		ok = false
	}else {
		ok = true
	}
	return ok
}
func checkAccountName(accountName string) bool{
	// Check the account name
	var ok bool = false
	var msg string
	err := db.QueryRow("SELECT accountName FROM user WHERE accountName=?", accountName).Scan(&msg)
	if err != nil {
		log.Println(err)
		ok = false
	}else {
		ok = true
	}
	return ok
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it

	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
func logoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clearSession(w)
	http.Redirect(w, r, "/login", 302)
}