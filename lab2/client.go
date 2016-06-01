//Lab3-CMPE-273
package main

import (
	"os"
	"sort"
	"log"
	"strings"
	"strconv"
	"net/http"
	"fmt"
	 "hash/crc32"
)

type Shard []uint32

//ConsistentHashing Structure
type  ChClient struct {
    serverMap   map[uint32]string
    hashCircle Shard
    isPresent map[int]bool
}

//Initialize ConsistentHashing Structure
func  ChClientCurrent() * ChClient {
    return & ChClient{
        serverMap:     make(map[uint32]string),
        hashCircle:    Shard{},
        isPresent: make(map[int]bool),
    }  
}  

//Add new server instances
func (ch * ChClient) addServerInstance(Id int,IP string) bool {
    if _, ok := ch.isPresent[Id]; ok {
        return false  
    }  
    ch.serverMap[ch.hashValueGeneration(IP)] = IP
    ch.isPresent[Id] = true
    ch. hashCircle = Shard{}
    for k := range ch.serverMap {
        ch. hashCircle = append(ch. hashCircle, k)
    }
    sort.Sort(ch. hashCircle)
    return true  
}

//Sort interfaces method, need to implement them
func (ch Shard) Len() int {
    return len(ch)
}

func (ch Shard) Less(i, j int) bool {
    return ch[i] < ch[j]
}

func (ch Shard) Swap(i, j int) {
    ch[i], ch[j] = ch[j], ch[i]
}

func main() {
	hashCircle :=  ChClientCurrent()
	params := strings.Split(os.Args[1],"-")
	start,_ := strconv.Atoi(params[0])
	finish,_ := strconv.Atoi(params[1])
	counter := start
	for i:=start; i<=finish ;i++ {
		hashCircle.addServerInstance(counter,strconv.Itoa(counter))
		counter++
	}
	mapValues := strings.Split(os.Args[2],",")

	// Putting the key
	for i:=0; i<len(mapValues) ;i++ {
		pair := strings.Split(mapValues[i],"->")
		key := pair[0]
		ip:= hashCircle.getIP(key)
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
func (ch * ChClient) getIP(key string) string {
    hash := ch.hashValueGeneration(key)  
    i := ch.ShardingSearch(hash)
    return ch.serverMap[ch. hashCircle[i]]
}

//Get the hash value
func (ch * ChClient) hashValueGeneration(key string) uint32 {
    	return crc32.ChecksumIEEE([]byte(key))
}

func (ch * ChClient) ShardingSearch(hash uint32) int {
    i := sort.Search(len(ch. hashCircle), func(i int) bool {return ch. hashCircle[i] >= hash })  
    if i < len(ch. hashCircle) {  
        if i == len(ch. hashCircle)-1 {  
            return 0  
        } else {  
            return i  
        }  
    } else {  
        return len(ch. hashCircle) - 1  
    }  
}

