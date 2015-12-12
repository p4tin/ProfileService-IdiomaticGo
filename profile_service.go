package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"encoding/json"
)

type Profile struct {
	Name 			string			`json: "username"`
	Password 		string			`json: "password"`
	Age 			int				`json: "age"`
	LastUpdated     time.Time
}

var Profiles []Profile

func init() {
	Profiles = make([]Profile, 3)
	Profiles[0].Name = "john@aol.com"
	Profiles[0].Password = "test1234"
	Profiles[0].Age = 23
	Profiles[0].LastUpdated = time.Now()

	Profiles[1].Name = "harry@comcast.com"
	Profiles[1].Password = "testing1"
	Profiles[1].Age = 44
	Profiles[1].LastUpdated = time.Now()

	Profiles[2].Name = "sally@microsoft.com"
	Profiles[2].Password = "pass4321"
	Profiles[2].Age = 63
	Profiles[2].LastUpdated = time.Now()
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")

	//v1 API calls
	router.HandleFunc("/v1/profile", handle_profile_no_params).Methods("GET", "POST", "PUT")
	router.HandleFunc("/v1/profile/{profileId}", handle_profile_with_params).Methods("GET", "DELETE")

	//v2 API calls
	router.HandleFunc("/v2/profile", v2_handle_profile_no_params).Methods("GET")

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Index Request")
	fmt.Fprintf(w, "Hello World!!!")
}

func handle_profile_no_params(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			get_profile(w, r)
		case "POST":
			create_profile(w,r)
		case "PUT":
			update_profile(w, r)
	}
}

func handle_profile_with_params(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			show_profile(w, r)
		case "DELETE":
			delete_profile(w,r)
	}
}

func get_profile(w http.ResponseWriter, r *http.Request) {
	log.Println("Get All Profiles")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Profiles)
}

func show_profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	profileId := vars["profileId"]
	log.Println("Show Profile: ", profileId)
	for _, profile := range Profiles {
		if(profileId == profile.Name) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(profile)
			return
		}
	}
	var repl struct{}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(repl)
}

func create_profile(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Profile")
	var resp struct{}

	p := new(Profile)
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&p)
	if error != nil {
		log.Println(error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	p.LastUpdated = time.Now()
	Profiles = append(Profiles, *p)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func update_profile(w http.ResponseWriter, r *http.Request) {
	log.Println("Create Profile")
	var resp struct{}

	p := new(Profile)
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&p)
	if error != nil {
		log.Println(error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	for i, prof := range Profiles {
		if(p.Name == prof.Name) {
			Profiles[i].Password = p.Password
			Profiles[i].Age = p.Age
			p.LastUpdated = time.Now()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return;
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func delete_profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var repl struct{}
	vars := mux.Vars(r)
	profileId := vars["profileId"]
	log.Println("Delete Profile: ", profileId)
	for i, profile := range Profiles {
		if(profileId == profile.Name) {
			Profiles = append(Profiles[:i], Profiles[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(repl)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(repl)
}



func v2_handle_profile_no_params(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			v2_get_profile(w, r)
	}
}


func v2_get_profile(w http.ResponseWriter, r *http.Request) {
	log.Println("Get All Profiles")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Profiles)
}
