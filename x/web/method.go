package web

import (
	"bytes"
	"ehelp/x/mlog"
	"encoding/json"
	"fmt"
	"net/http"
)

var logMarshal = mlog.NewTagLog("rest_cetm")

func ResParamArrUrlClient(url string, objArray interface{}, objRes interface{}) error {
	result, err := json.Marshal(objArray)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(result))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logMarshal.Errorf(url, "error marshal post")
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&objRes)
}

func ResUrlClientGet(url string, objRes interface{}) error {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		logMarshal.Errorf(url, "error marshal post")
		panic(err)
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&objRes)
}
