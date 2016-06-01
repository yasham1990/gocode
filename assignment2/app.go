// Assignment2 persist the personal preference data that you get from the assignment 1 into a database
package main

import (
	"encoding/json"
	"github.com/drone/routes"
	"github.com/mkilling/goejdb"
	"github.com/naoina/toml"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Person struct {
	Emailid        string `json:"email"`
	Zip            string `json:"zip"`
	Country        string `json:"country"`
	Profession     string `json:"profession"`
	Favorite_color string `json:"favorite_color"`
	Is_smoking     string `json:"is_smoking"`
	Favorite_sport string `json:"favorite_sport"`
	Food           Food   `json:"food"`
	Music          Music  `json:"music"`
	Movie          Movie  `json:"movie"`
	Travel         Travel `json:"travel"`
}
type Food struct {
	Type          string `json:"type"`
	Drink_alcohol string `json:"drink_alcohol"`
}
type Music struct {
	Spotify_user_id string `json:"spotify_user_id"`
}
type Movie struct {
	Tv_shows [3]string `json:"tv_shows"`
	Movies   [3]string `json:"movies"`
}
type Travel struct {
	Flight Flight `json:"flight"`
}
type Flight struct {
	Seat string `json:"seat"`
}

type tomlConfiguration struct {
	Database struct {
		File_name string
		Port_num  int
	}
	Replication struct {
		Rpc_server_port_num int
		Replica             []string
	}
}

type Listener int

var collection *goejdb.EjColl
var configuration tomlConfiguration

func main() {

	//Create configuration file
	ConfigurationFileReading()
	// Create a new database file and open it
	jb, err := goejdb.Open(configuration.Database.File_name, goejdb.JBOWRITER|goejdb.JBOCREAT)
	if err != nil {
		os.Exit(1)
	}
	//Create configuration file
	go RPCServerInitialization()

	//Create person profile
	collection, _ = jb.CreateColl("personprofile", nil)
	mux := routes.New()
	mux.Post("/profile", CreateNewProfile)
	mux.Get("/profile/:email", GetProfile)
	mux.Put("/profile/:email", UpdateProfile)
	mux.Del("/profile/:email", DeleteProfile)
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":"+strconv.Itoa(configuration.Database.Port_num), nil)
	// Don't forget to close the database connection
	jb.Close()

}

// This function will create a new person as per the value passed in the jSON object from request.
func CreateNewProfile(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	person := &Person{}
	err := json.Unmarshal([]byte(body), &person)
	var recordBson []byte
	if err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else if err := request.Body.Close(); err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else {
		recordBson, _ = bson.Marshal(person)
		collection.SaveBson(recordBson)
		writer.WriteHeader(http.StatusCreated)
	}
	for _, each := range configuration.Replication.Replica {
		client, err := rpc.Dial("tcp", each)
		if err != nil {
			log.Fatal(err)
		}
		var acknowledge bool
		err = client.Call("Listener.CreateReplication", recordBson, &acknowledge)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// This function will get person details from the storage map if a person object is found as per the requested email id.
func GetProfile(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	email := params.Get(":email")
	databaseQuery := `{"emailid": "` + email + `"}`
	res, _ := collection.Find(databaseQuery)
	//Check if the person is there or not in the DB
	if len(res) > 0 {
		p := &Person{}
		bson.Unmarshal(res[0], &p)
		profiles, err := json.Marshal(p)
		if err != nil {
			panic(err)
			writer.Write([]byte(err.Error()))
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write(profiles)
		}
	} else {
		log.Println("Person not found...")
	}

}

// This function will delete person details from the storage map if a person object is found as per the requested emailid.
func DeleteProfile(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	email := params.Get(":email")
	databaseQuery := `{"emailid": "` + email + `"}`
	res, _ := collection.Find(databaseQuery)
	//Check if the person is there or not in the DB
	if len(res) > 0 {
		collection.Update(`{"emailid": "` + email + `", "$dropall" : true}`)
		writer.WriteHeader(http.StatusNoContent)
		for _, each := range configuration.Replication.Replica {
			client, err := rpc.Dial("tcp", each)
			if err != nil {
				log.Fatal(err)
			}
			var acknowledge bool
			err = client.Call("Listener.DeleteReplication", []byte(email), &acknowledge)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Println("Person not found...")
	}

}

// This function will update person details from the storage map if a person object is found as per the requested emailid.
func UpdateProfile(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	personNew := &Person{}
	err := json.Unmarshal([]byte(body), &personNew)
	if err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else {
		params := request.URL.Query()
		email := params.Get(":email")
		databaseQuery := `{"emailid": "` + email + `"}`
		res, _ := collection.Find(databaseQuery)
		person := &Person{}
		if len(res) > 0 {
			bson.Unmarshal(res[0], &person)
			if len(personNew.Emailid) == 0 {
				personNew.Emailid = person.Emailid
			}
			if len(personNew.Zip) == 0 {
				personNew.Zip = person.Zip
			}
			if len(personNew.Country) == 0 {
				personNew.Country = person.Country
			}
			if len(personNew.Profession) == 0 {
				personNew.Profession = person.Profession
			}
			if len(personNew.Favorite_color) == 0 {
				personNew.Favorite_color = person.Favorite_color
			}
			if len(personNew.Is_smoking) == 0 {
				personNew.Is_smoking = person.Is_smoking
			}
			if len(personNew.Favorite_sport) == 0 {
				personNew.Favorite_sport = person.Favorite_sport
			}
			if len(personNew.Food.Type) == 0 {
				personNew.Food.Type = person.Food.Type
			}
			if len(personNew.Food.Drink_alcohol) == 0 {
				personNew.Food.Drink_alcohol = person.Food.Drink_alcohol
			}
			if len(personNew.Music.Spotify_user_id) == 0 {
				personNew.Music.Spotify_user_id = person.Music.Spotify_user_id
			}
			if len(personNew.Movie.Tv_shows[0]) == 0 {
				personNew.Movie.Tv_shows[0] = person.Movie.Tv_shows[0]
			}
			if len(personNew.Movie.Tv_shows[1]) == 0 {
				personNew.Movie.Tv_shows[1] = person.Movie.Tv_shows[1]
			}
			if len(personNew.Movie.Tv_shows[2]) == 0 {
				personNew.Movie.Tv_shows[2] = person.Movie.Tv_shows[2]
			}
			if len(personNew.Movie.Movies[0]) == 0 {
				personNew.Movie.Movies[0] = person.Movie.Movies[0]
			}
			if len(personNew.Movie.Movies[1]) == 0 {
				personNew.Movie.Movies[1] = person.Movie.Movies[1]
			}
			if len(personNew.Movie.Movies[2]) == 0 {
				personNew.Movie.Movies[2] = person.Movie.Movies[2]
			}
			if len(personNew.Travel.Flight.Seat) == 0 {
				personNew.Travel.Flight.Seat = person.Travel.Flight.Seat
			}
		}
		updatedProfile, err := json.Marshal(personNew)
		if err != nil {
			panic(err)
			writer.Write([]byte(err.Error()))
		} else {
			var databaseQuery = `{"emailid": "` + email + `", "$set":` + string(updatedProfile) + `}`
			collection.Update(databaseQuery)
		}
		writer.WriteHeader(http.StatusNoContent)
		for _, each := range configuration.Replication.Replica {
			client, err := rpc.Dial("tcp", each)
			if err != nil {
				log.Fatal(err)
			}
			var acknowledge bool
			err = client.Call("Listener.UpdationReplication", updatedProfile, &acknowledge)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (l *Listener) CreateReplication(line []byte, ack *bool) error {
	collection.SaveBson(line)
	return nil
}
func (l *Listener) DeleteReplication(line []byte, ack *bool) error {
	email := string(line)
	collection.Update(`{"emailid": "` + email + `", "$dropall" : true}`)
	return nil
}
func (l *Listener) UpdationReplication(line []byte, ack *bool) error {
	person := &Person{}
	err := json.Unmarshal([]byte(line), &person)
	if err == nil {
		var databaseQuery = `{"emailid": "` + person.Emailid + `", "$set":` + string(line) + `}`
		collection.Update(databaseQuery)
	}
	return nil
}
func ConfigurationFileReading() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	//Dont forget to close.
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if err := toml.Unmarshal(buf, &configuration); err != nil {
		panic(err)
	}
	for i := 0; i < len(configuration.Replication.Replica); i++ {
		configuration.Replication.Replica[i] = strings.TrimLeft(configuration.Replication.Replica[i], "http://")
	}
}
func RPCServerInitialization() {
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(configuration.Replication.Rpc_server_port_num))
	if err != nil {
		log.Fatal(err)
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}
	listener := new(Listener)
	//Register the port
	rpc.Register(listener)
	//Accept the inbound
	rpc.Accept(inbound)
}
