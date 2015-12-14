// mongo - CRUD package for the Profile API micro-service
package mongo

import (
	"time"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/base64"

	"ProfileService-IdiomaticGo/crypto"
)

// Profile - is the memory representation of one user profile
type Profile struct {
	Name 			string			`json: "username"`
	Password 		string			`json: "password"`
	Age 			int				`json: "age"`
	LastUpdated     time.Time
}

func init() {
	tmp, _ := crypto.Encrypt([]byte( "test1234"))
	p1 := Profile{Name: "john@aol.com", Password: base64.StdEncoding.EncodeToString(tmp), Age: 23, LastUpdated: time.Now()}
	p1.CreateOrUpdateProfile()

	tmp, _ = crypto.Encrypt([]byte( "testing1"))
	p1 = Profile{Name: "harry@comcast.com", Password: base64.StdEncoding.EncodeToString(tmp), Age: 44, LastUpdated: time.Now()}
	p1.CreateOrUpdateProfile()

	tmp, _ = crypto.Encrypt([]byte( "pass4321"))
	p1 = Profile{Name: "sally@microsoft.com", Password: base64.StdEncoding.EncodeToString(tmp), Age: 63, LastUpdated: time.Now()}
	p1.CreateOrUpdateProfile()
}

// GetProfiles - Returns all the profile in the Profiles Collection
func GetProfiles() []Profile {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return nil
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	var profiles []Profile
	err = c.Find(bson.M{}).All(&profiles)

	return profiles
}


// ShowProfile - Returns the profile in the Profiles Collection with name equal to the id parameter (id == name)
func ShowProfile(id string) Profile {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return Profile{}
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	profile := Profile{}
	err = c.Find(bson.M{"name": id}).One(&profile)

	return profile
}


// DeleteProfile - Deletes the profile in the Profiles Collection with name equal to the id parameter (id == name)
func DeleteProfile(id string) bool {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return false
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	err = c.RemoveId(id)

	return true
}


// CreateOrUpdateProfile - Creates or Updates (Upsert) the profile in the Profiles Collection with id parameter
func (p *Profile) CreateOrUpdateProfile() bool {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return false
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("ProfileService").C("Profiles")
	_ , err = c.UpsertId( p.Name, p )
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
		return false
	}
	return true
}


