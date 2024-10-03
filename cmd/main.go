package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"
)

func data(w http.ResponseWriter, req *http.Request) {

	dataTemplate, err := os.ReadFile("data.template")
	check(err)
	fmt.Print(string(dataTemplate))

	tmpl, err := template.New("data").Parse(string(dataTemplate))
	check(err)

	dataJson, err := os.ReadFile("data.json")
	check(err)
	fmt.Print(string(dataJson))
	var params map[string]interface{}
	err = json.Unmarshal(dataJson, &params)
	check(err)

	var e map[string]interface{}

	if req.Method == http.MethodPost {
		err := json.NewDecoder(req.Body).Decode(&e)
		if err != nil {
			// http.Error(w, err.Error(), http.StatusBadRequest)
			logrus.WithError(err).Error("Cannot parse")
			return
		}

		var results bytes.Buffer
		writer := io.Writer(&results)
		err = tmpl.Execute(writer, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r := results.Bytes()

		// Send modified webhook request
		url, ok := params["url"].(string)
		if !ok {
			http.Error(w, "could not convert url", http.StatusInternalServerError)
			return
		}
		if url == "" {
			http.Error(w, "url not found", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var headers map[string]interface{}
		headersJson := params["headers"]
		if headersJson != nil {
			logrus.Info(headersJson)
			headers, ok = headersJson.(map[string]interface{})
			logrus.Info(headers)
			if !ok {
				http.Error(w, "could not convert headers", http.StatusInternalServerError)
				return
			}
		}

		for name, value := range headers {
			hValue, ok := value.(string)
			if !ok {
				http.Error(w, "could not convert headers value", http.StatusInternalServerError)
				return
			}
			req.Header.Set(name, hValue)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			_, err = w.Write(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, string(body), resp.StatusCode)
		}
	}
}

func main() {

	http.HandleFunc("/webhooks/data", data)

	http.ListenAndServe(":8090", nil)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
