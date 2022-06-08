package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type Article struct {
	Query string `json:"query"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Score float64 `json:"score"`
}

type SimilarWordSearchPlugin struct {
	plugin.MattermostPlugin
}

func (p *SimilarWordSearchPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "SimilarWordSearchPlugin\n")
	case http.MethodPost:
		// json request body from plugin webapp decode to golang struct 'term'
		var term struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&term); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// golang struct 'term' to json
		body, err := json.Marshal(term)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }

		// send term json (bytebuffer) to outside API (in this case, similarwordssearchapi:8000)
		req, err := http.NewRequest(
			http.MethodPost, 
			"http://similarwordssearchapi:8000/similarwords/",
			bytes.NewBuffer([]byte(body)),
		 )
    if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }

		// parse json response from outside API to golang struct '[]*Article'
		results := make([]*Article, 0)
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// golang struct '[]*Article' to json
		res, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write response (to plugin webapp)
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	default:
		fmt.Fprint(w, "Method not allowed.\n")
	}
}

func main() {
	plugin.ClientMain(&SimilarWordSearchPlugin{})
}
