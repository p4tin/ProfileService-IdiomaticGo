// routes - Routing package for the Profile API micro-service
package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"

	"encoding/base64"

	"ProfileService-IdiomaticGo/mongo"
	"ProfileService-IdiomaticGo/crypto"
)

func Start() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")

	//v1 API calls
	router.HandleFunc("/v1/profile", handleProfileNoParams).Methods("GET", "POST", "PUT")
	router.HandleFunc("/v1/profile/{profileId}", handleProfileWithParams).Methods("GET", "DELETE")

	//v2 API calls
	router.HandleFunc("/v2/profile", v2HandleProfileNoParams).Methods("GET")

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Index Request")
	fmt.Fprintf(w, "Hello World!!!")
}

func handleProfileNoParams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProfile(w, r)
	case "POST":
		createProfile(w, r)
	case "PUT":
		updateProfile(w, r)
	}
}

func handleProfileWithParams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		showProfile(w, r)
	case "DELETE":
		deleteProfile(w, r)
	}
}

//getProfile - Get all profiles from the DB
func getProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Get All Profiles")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(Profiles)
	json.NewEncoder(w).Encode(mongo.GetProfiles())
}


//createProfile - Create a profile in the DB
func createProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Profile")
	var resp struct{}

	p := mongo.Profile{}
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&p)
	if error != nil {
		log.Println(error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	tmp, _ := crypto.Encrypt([]byte(p.Password))
	p.Password = base64.StdEncoding.EncodeToString(tmp)
	p.LastUpdated = time.Now()
	p.CreateOrUpdateProfile()

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

//updateProfile - Update a profile in the DB
func updateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Profile")
	var resp struct{}

	p := mongo.Profile{}
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&p)
	if error != nil {
		log.Println(error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	tmp, _ := crypto.Encrypt([]byte(p.Password))
	p.Password = base64.StdEncoding.EncodeToString(tmp)
	p.LastUpdated = time.Now()
	p.CreateOrUpdateProfile()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

//showProfile - given the profile ID in the request this returns the profile to the caller
func showProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	profileID := vars["profileId"]
	log.Println("Show Profile: ", profileID)
	repl := mongo.ShowProfile(profileID)
	if repl.Name == "" {
		var repl struct{}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(repl)
		return
	}
	json.NewEncoder(w).Encode(repl)
}

// deleteProfile - Removes a profile from the DB
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var repl struct{}
	vars := mux.Vars(r)
	profileID := vars["profileId"]
	log.Println("Delete Profile: ", profileID)
	mongo.DeleteProfile(profileID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}



func v2HandleProfileNoParams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		v2GetProfile(w, r)
	}
}


func v2GetProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("Get All Profiles")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(Profiles)
	json.NewEncoder(w).Encode(mongo.GetProfiles())
}
