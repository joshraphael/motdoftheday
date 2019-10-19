package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
	"gitlab.com/joshraphael/motdoftheday/pkg/post"
)

func (r Rest) SaveHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			msg := "Error reading save request data: " + err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		post := post.New(apierror.MethodHTTP)
		if err := json.Unmarshal(data, &post); err != nil {
			msg := "Error marshalling save json data: " + err.Error()
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if apiErr := r.processor.SaveForm(post); apiErr != nil {
			msg := "Error processing save request: " + apiErr.Error()
			log.Println(msg)
			http.Error(w, msg, apiErr.Code())
			return
		}
		w.WriteHeader(http.StatusCreated)
		log.Println("Saved post")
		return
	}
}
