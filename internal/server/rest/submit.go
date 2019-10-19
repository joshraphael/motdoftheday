package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

func (r Rest) SubmitHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			msg := "Error reading submit request data: " + err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		post := post.New(apierror.MethodHTTP)
		if err := json.Unmarshal(data, &post); err != nil {
			msg := "Error marshalling sumbit json data: " + err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if err = r.processor.SubmitForm(post); err != nil {
			msg := "Error processing submit request: " + err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		log.Println("Submitted post")
		return
	}
}
