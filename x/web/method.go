package web

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ResParamArrUrlClient(url string, objArray interface{}, objRes interface{}) error {
	result, err := json.Marshal(objArray)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&objRes)
}

func ResUrlClientGet(url string, objRes interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var objres interface{}
	return json.NewDecoder(resp.Body).Decode(&objres)
}
