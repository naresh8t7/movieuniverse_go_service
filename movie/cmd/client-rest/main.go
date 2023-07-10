package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// get configuration
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
	flag.Parse()

	pfx := "3"

	var body string

	// Call Create
	resp, err := http.Post(*address+"/v1/movie", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"movie": {
				"title":"Star Wars %s"
			}
		}
	`, pfx)))
	if err != nil {
		log.Fatalf("failed to call Create method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Create response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Add tags
	tagData := map[string]interface{}{
		"title": "Star Wars 3",
		"tags": []string{"Action"},
	}

	reqData, err := json.Marshal(tagData)
	if err != nil {
		log.Fatalf("failed to marshall data: %v", err)
	}
	tagReq, err := http.NewRequest("PUT", *address+"/v1/movie/tags",  bytes.NewBuffer(reqData))
	client := http.Client{
		Timeout:       time.Duration(60 * time.Second),
	}
	tagReq.Header.Set("Content-type", "application/json")
	resp, err = client.Do(tagReq)

	if err != nil {
		log.Fatalf("failed to call add tags method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read add tags response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Add tags response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/movie", "title (test)"))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call ReadAll
	resp, err = http.Get(*address + "/v1/movie/all")
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read ReadAll response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}

