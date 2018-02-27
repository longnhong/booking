package web

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ResParamArrUrlClient(url string, objArray interface{}) interface{} {
	result, err := json.Marshal(objArray)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var objres interface{}
	json.NewDecoder(resp.Body).Decode(&objres)
	return &objres
}

func ResUrlClientGet(url string) interface{} {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var objres interface{}
	json.NewDecoder(resp.Body).Decode(&objres)
	return &objres
}
