//Lab3-CMPE-273
package main

import (
	"os"
	"log"
	"strings"
	"strconv"
	"net/http"
	"fmt"
	"hash"
	 "hash/crc32"
)

// Server node structs 
type serverNode struct {
	server  []byte
}

//Rendezvous structure created
type Rendezvous struct {
	serverNodes  []serverNode
	rHashing hash.Hash32
}

// Initialize Rendezvous.
func Initialize(serverNodes ...string) *Rendezvous {
	hash := &Rendezvous{}
	hash.rHashing = crc32.New(hash32Table)
	hash.addServerInstance(serverNodes...)
	return hash
}

// addServerInstance will take server node and add them to the hash.
func (r *Rendezvous) addServerInstance(serverNodes ...string) {
	for _, server := range serverNodes {
		r.serverNodes = append(r.serverNodes, serverNode{[]byte(server)})
	}
}

// main method for execution
func main() {
	initR :=  Initialize()
	params := strings.Split(os.Args[1],"-")
	start,_ := strconv.Atoi(params[0])
	finish,_ := strconv.Atoi(params[1])
	counter := start
	for i:=start; i<=finish ;i++ {
		initR.addServerInstance(strconv.Itoa(counter))
		counter++
	}
	mapValues := strings.Split(os.Args[2],",")

	// Putting the key
	for i:=0; i<len(mapValues) ;i++ {
		pair := strings.Split(mapValues[i],"->")
		key := pair[0]
		ip:= initR.getIP(key)
		url :="http://127.0.0.1:"+ip+"/"+key+"/"+pair[1]
		client := &http.Client{}
		request, err := http.NewRequest("PUT", url, nil)
		response,err:= client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(response)
	}
}

//Get IP address to which key will be mapped
func (r * Rendezvous) getIP(key string) string {
    	keyInBytes := []byte(key)
	var maximumCal uint32
	var highNode []byte
	var hasCal uint32

	for _, ser := range r.serverNodes {
		hasCal = r.hashValueGeneration(ser.server, keyInBytes)
		if hasCal > maximumCal {
			maximumCal = hasCal
			highNode = ser.server
		}
	}

	return string(highNode)
}

//Get the hash value
func (r *Rendezvous) hashValueGeneration(server, key []byte) uint32 {
	r.rHashing.Reset()
	r.rHashing.Write(key)
	r.rHashing.Write(server)
	return r.rHashing.Sum32()
}

var (
	hash32Table = crc32.MakeTable(crc32.Castagnoli)
)
