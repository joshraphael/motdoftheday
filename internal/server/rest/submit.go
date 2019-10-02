package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/joshraphael/diary/pkg/processors"
)

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "Error reading request data: " + err.Error()
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		var submit processors.Submit
		if err := json.Unmarshal(data, &submit); err != nil {
			msg := "Error marshalling json data: " + err.Error()
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if err = processors.SubmitForm(submit); err != nil {
			msg := "Error processing submit request: " + err.Error()
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
}
