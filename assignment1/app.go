// Building a simple RESTful service to manage personal perferences that can be shared with service providers like travel agencies so that they can provide personalization in their services.
package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"net/http"
	"github.com/drone/routes"
)

type Person struct {
	Emailid        string `json:"email"`
	Zip            string `json:"zip"`
	Country        string `json:"country"`
	Profession     string `json:"profession"`
	Favorite_color string `json:"favorite_color"`
	Is_smoking     string `json:"is_smoking"`
	Favorite_sport string `json:"favorite_sport"`
	Food Food 	      `json:"food"`
	Music Music           `json:"music"`
	Movie Movie	      `json:"movie"`
	Travel Travel         `json:"travel"`
	}
type Food struct {
		Type     string `json:"type"`
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
		Flight Flight      `json:"flight"`
	} 	
type Flight struct {
		Seat string 	   `json:"seat"`
	} 
//Map to store person details
var storageMap = make(map[string]Person)

func main() {
	mux := routes.New()
	mux.Post("/profile", CreateNewProfile)
	mux.Get("/profile/:email", GetProfile)
	mux.Put("/profile/:email", UpdateProfile)
	mux.Del("/profile/:email", DeleteProfile)
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

// This function will create a new person as per the value passed in the jSON object from request.
func CreateNewProfile(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	person := &Person{}
	err := json.Unmarshal([]byte(body), &person)
	if err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else if err := request.Body.Close(); err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else {
		storageMap[person.Emailid] = *person
		writer.WriteHeader(http.StatusCreated)
	}
}

// This function will get person details from the storage map if a person object is found as per the requested email id.
func GetProfile(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	email := params.Get(":email")
	if storageMap[email].Emailid == email {
		person, err := json.Marshal(storageMap[email])
		if err != nil {
			panic(err)
			writer.Write([]byte(err.Error()))
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write(person)
		}
	} else {
		log.Println("Person not found...")
	}
}

// This function will delete person details from the storage map if a person object is found as per the requested emailid.
func DeleteProfile(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	email := params.Get(":email")
	if storageMap[email].Emailid == email {
		delete(storageMap, email)
		writer.WriteHeader(http.StatusNoContent)
	} else {
		log.Println("Person not found...")
	}
}

// This function will update person details from the storage map if a person object is found as per the requested emailid.
func UpdateProfile(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	person := &Person{}
	err := json.Unmarshal([]byte(body), &person)
	if err != nil {
		panic(err)
		writer.Write([]byte(err.Error()))
	} else {
		params := request.URL.Query()
		email := params.Get(":email")
		//check which field is updated by the length method
		if storageMap[email].Emailid == email {
			if len(person.Emailid) == 0 {
				person.Emailid = storageMap[email].Emailid
			}
			if len(person.Zip) == 0 {
				person.Zip = storageMap[email].Zip
			}
			if len(person.Country) == 0 {
				person.Country = storageMap[email].Country
			}
			if len(person.Profession) == 0 {
				person.Profession = storageMap[email].Profession
			}
			if len(person.Favorite_color) == 0 {
				person.Favorite_color = storageMap[email].Favorite_color
			}
			if len(person.Is_smoking) == 0 {
				person.Is_smoking = storageMap[email].Is_smoking
			}
			if len(person.Favorite_sport) == 0 {
				person.Favorite_sport = storageMap[email].Favorite_sport
			}
			if len(person.Food.Type) == 0 {
				person.Food.Type = storageMap[email].Food.Type
			}
			if len(person.Food.Drink_alcohol) == 0 {
				person.Food.Drink_alcohol = storageMap[email].Food.Drink_alcohol
			}
			if len(person.Music.Spotify_user_id) == 0 {
				person.Music.Spotify_user_id = storageMap[email].Music.Spotify_user_id
			}
			if len(person.Movie.Tv_shows[0]) == 0 {
				person.Movie.Tv_shows[0] = storageMap[email].Movie.Tv_shows[0]
			}
			if len(person.Movie.Tv_shows[1]) == 0 {
				person.Movie.Tv_shows[1] = storageMap[email].Movie.Tv_shows[1]
			}
			if len(person.Movie.Tv_shows[2]) == 0 {
				person.Movie.Tv_shows[2] = storageMap[email].Movie.Tv_shows[2]
			}
			if len(person.Movie.Movies[0]) == 0 {
				person.Movie.Movies[0] = storageMap[email].Movie.Movies[0]
			}
			if len(person.Movie.Movies[1]) == 0 {
				person.Movie.Movies[1] = storageMap[email].Movie.Movies[1]
			}
			if len(person.Movie.Movies[2]) == 0 {
				person.Movie.Movies[2] = storageMap[email].Movie.Movies[2]
			}
			if len(person.Travel.Flight.Seat) == 0 {
				person.Travel.Flight.Seat = storageMap[email].Travel.Flight.Seat
			}
			storageMap[person.Emailid] = *person
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
