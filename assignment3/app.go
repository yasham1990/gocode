package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Id                         int64   `bson:"id"`
	Success_http_response_code int     `bson:"success_http_response_code"`
	Max_retries                int     `bson:"max_retries"`
	Callback_webhook_url       string  `bson:"callback_webhook_url"`
	Request                    Request `bson:"request"`
}

type Request struct {
	Url          string       `bson:"url"`
	Method       string       `bson:"method"`
	Http_headers Http_headers `bson:"http_headers"`
	Body         Body         `json:"body" bson:"body"`
}
type Body struct {
	Foo string `json:"foo" bson:"foo"`
}

type Http_headers struct {
	Content_Type string `bson:"Content-Type"`
	Accept       string `bson:"Accept"`
}

type OutputData struct {
	Job                    Job `bson:"job"`
	Person                 `bson:"input"`
	Output                 Output `bson:"output"`
	Callback_response_code int    `bson:"callback_response_code"`
}

type Job struct {
	Status      string `bson:"status"`
	Num_retries int    `bson:"num_retries"`
}
type Output struct {
	Response struct {
		Http_response_code int `bson:"http_response_code"`
		Http_headers       struct {
			Date           string `bson:"Date"`
			Content_Type   string `bson:"Content-Type"`
			Content_Length int64  `bson:"Content-Length"`
		} `bson:"http_headers"`
		Body struct {
			Hello string `json:"hello" bson:"hello"`
		} `json:"body" bson:"body"`
	} `bson:"response"`
}

type JsonResponseBody struct {
	Hello string `json:"hello" bson:"hello"`
}

var readInput = false

//Process all incoming requests
func RequestProcessing(inputFileObject *Person, numberOfTries int) *OutputData {
	outputData := new(OutputData)
	outputData.Person = *inputFileObject
	outputData.Person.Max_retries = inputFileObject.Max_retries
	jsonReqBody, _ := json.Marshal(inputFileObject.Request.Body)
	request, _ := http.NewRequest(inputFileObject.Request.Method, inputFileObject.Request.Url, bytes.NewReader(jsonReqBody))
	request.Header.Set("Content-Type", inputFileObject.Request.Http_headers.Content_Type)
	request.Header.Add("Accept", inputFileObject.Request.Http_headers.Accept)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err == nil {
		if resp.StatusCode == inputFileObject.Success_http_response_code {
			outputData.Job.Status = "COMPLETED"
			outputData.Job.Num_retries = numberOfTries
			outputData.Output.Response.Http_response_code = resp.StatusCode
			outputData.Output.Response.Http_headers.Content_Type = resp.Header.Get("Content-Type")
			fmt.Println(resp.Header.Get("Content-Type"))
			outputData.Output.Response.Http_headers.Date = resp.Header.Get("Date")
			outputData.Output.Response.Http_headers.Content_Length = resp.ContentLength
			respBody, _ := ioutil.ReadAll(resp.Body)
			jsonResponseBody := &JsonResponseBody{}
			jsonResponseBodyerr := json.Unmarshal([]byte(respBody), &jsonResponseBody)
			if jsonResponseBodyerr != nil {
				panic(jsonResponseBodyerr)
			}
			outputData.Output.Response.Body.Hello = jsonResponseBody.Hello
			callbackRequest, _ := http.NewRequest("POST", outputData.Person.Callback_webhook_url, bytes.NewReader(respBody))
			callbackResponse, err1 := client.Do(callbackRequest)
			if err1 == nil {
				outputData.Callback_response_code = callbackResponse.StatusCode
			} else {
				fmt.Print("Callback Sever Down........")
			}
		} else {
			if numberOfTries >= outputData.Person.Max_retries {
				outputData.Job.Status = "FAILED"
				outputData.Job.Num_retries = numberOfTries
				callbackRequest, _ := http.NewRequest("POST", outputData.Person.Callback_webhook_url, bytes.NewReader([]byte("howdy\n")))
				client.Do(callbackRequest)
			} else {
				outputData.Job.Status = "STILL_TRYING"
				outputData.Job.Num_retries = numberOfTries
			}
		}
	} else {
		if numberOfTries >= outputData.Person.Max_retries {
			outputData.Job.Status = "FAILED"
			outputData.Job.Num_retries = numberOfTries
			callback_req, _ := http.NewRequest("POST", outputData.Person.Callback_webhook_url, bytes.NewReader([]byte("FAILED")))
			client.Do(callback_req)
		} else {
			outputData.Job.Status = "STILL_TRYING"
			outputData.Job.Num_retries = numberOfTries
		}
	}
	return outputData
}

//write to file
func writeFile(output *OutputData) {
	fmt.Println("Saving=> Status: " + output.Job.Status + ", Retries: " + strconv.Itoa(output.Job.Num_retries))
	bsonFormatOutput := createBSON(*output)
	file, err := os.Create("output.bson")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	numBytes, err := file.Write(bsonFormatOutput)
	fmt.Printf("\nwrote %d bytes to %s\n", numBytes, "output.bson")
	file.Sync()
}

//read to file
func readFile(fileName string) *Person {
	inputFromFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	sz := len(inputFromFile)
	if sz == 0 {
		return nil
	}
	dataInput := new(Person)
	bson.Unmarshal(inputFromFile, &dataInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read from BSON file: Person.Request.Body=%s\n", dataInput.Request.Body.Foo)
	fmt.Printf("Read from BSON file: Person.Request.Http_headers.Content_Type=%s\n", dataInput.Request.Http_headers.Content_Type)
	fmt.Printf("Read from BSON file: Person.Request.Request.Http_headers.Accept=%s\n", dataInput.Request.Http_headers.Accept)
	return dataInput
}

func OutputFileReading(outFile string) *OutputData {
	in, err := ioutil.ReadFile(outFile)
	if err != nil {
		panic(err)
	}
	dataOutput := new(OutputData)
	err = bson.Unmarshal(in, &dataOutput)
	if err != nil {
		panic(err)
	}
	return dataOutput
}

func createBSON(output OutputData) []byte {
	data, err := bson.Marshal(&output)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	fmt.Printf("%q", data)
	jsonop, _ := json.Marshal(output)
	fmt.Println()
	fmt.Println(string(jsonop[:]))
	return data
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func task() {
	t := makeTimestamp()
	fmt.Printf("Running task...@%d\n", t)
	if readInput == false {
		var inputFileObject *Person = readFile("input.bson")
		output := RequestProcessing(inputFileObject, 0)
		writeFile(output)
		readInput = true
	}
	var output *OutputData = OutputFileReading("output.bson")
	if output.Job.Status == "STILL_TRYING" {
		outputProcessed := RequestProcessing(&output.Person, output.Job.Num_retries+1)
		writeFile(outputProcessed)
	}
}

func main() {
	period:=os.Args[1]
    	duration,_:=strconv.Atoi(period)
	s := gocron.NewScheduler()
	s.Every(uint64(duration)).Seconds().Do(task)
	<-s.Start()
	task()
}
