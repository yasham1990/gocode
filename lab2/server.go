//Lab3-CMPE-273
package main

import (
    "fmt"
    "os"
    "net/http"
    "encoding/json"
    "github.com/julienschmidt/httprouter"
    "strings"
    "strconv"
)

type KeyValuePair struct {
  Key int 		`json:"key"`
  Value string 		`json:"value"`
}

var  ServerAndDataStoreMap map[int]map[int]string

//Main function
func main() {
  ServerAndDataStoreMap = make(map[int]map[int]string)
	params := strings.Split(os.Args[1],"-")
	start,_ := strconv.Atoi(params[0])
	finish,_ := strconv.Atoi(params[1])
	//Make a channel
	channelRoutine := make(chan bool)
	for i:=start; i<=finish ;i++ {
		portnumber := ":"+strconv.Itoa(i)
		//Create a new HTTP router for serving requests
		router := httprouter.New()
		router.GET("/", getAllKeysFromDatastore)
		router.GET("/:key", getKeyFromDatastore)
		router.PUT("/:key/:value", putKeyValue)
		go func() {
			http.ListenAndServe("localhost"+portnumber, router)			
		}()
	}
	<-channelRoutine
}

// This function will get all the values for all the keys.
func getAllKeysFromDatastore(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
	//take values from request
	paramsInfo := strings.Split(request.Host,":")
	portnumber,_:= strconv.Atoi(paramsInfo[1])
	datastore :=  ServerAndDataStoreMap[portnumber];
	var jsonResponse string
	if len(datastore) > 0	{
		counter := 0
		for key, value := range datastore {
			jsonResponse+= `{"key":`+strconv.Itoa(key)+`,"value":"`+value+`"},`
			counter++
		}
		jsonResponse=jsonResponse[0:len(jsonResponse)-1]
		if counter > 1	{
			jsonResponse = `[`+jsonResponse+`]`
		}
	}
	// If content not found give No content
	if jsonResponse !="" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK) // Status code for success
			fmt.Fprintf(rw, "%s", jsonResponse)
	} else {
			rw.WriteHeader(http.StatusNoContent)
	}
}

// This function will get value for a particular key.
func getKeyFromDatastore(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
	//take values from request
	paramsInfo := strings.Split(request.Host,":")
	portnumber,_:= strconv.Atoi(paramsInfo[1])
 	keyId,_:= strconv.Atoi(p.ByName("key"))
	datastore :=  ServerAndDataStoreMap[portnumber];
 	var keyvalue KeyValuePair
	value := datastore[keyId]
	keyvalue.Key = keyId
	keyvalue.Value = value
	// If content not found give No content
	if keyvalue.Value !="" {
			jsonResponse, _ := json.Marshal(&keyvalue)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK) // Status code for success
			fmt.Fprintf(rw, "%s", jsonResponse)
	} else {
			rw.WriteHeader(http.StatusNoContent)
	}
	
}

// This function will put values for the corresponding keys.
func putKeyValue(rw http.ResponseWriter, request *http.Request, p httprouter.Params) {
	//take values from request
	paramsInfo := strings.Split(request.Host,":")
	portnumber,_:= strconv.Atoi(paramsInfo[1])
	keyId,_:= strconv.Atoi(p.ByName("key"))
	var datastore map[int]string
	//Check if the datastore is present or not
	if( ServerAndDataStoreMap[portnumber]!=nil){
		datastore =  ServerAndDataStoreMap[portnumber]
	}else{
		datastore = make(map[int]string)
	}
	var keyvalue KeyValuePair
  	keyvalue.Key = keyId
  	keyvalue.Value = p.ByName("value")
	datastore[keyvalue.Key] = keyvalue.Value
	ServerAndDataStoreMap[portnumber] = datastore;
  	rw.WriteHeader(http.StatusOK) // Status code for success
  	fmt.Fprint(rw)
}
