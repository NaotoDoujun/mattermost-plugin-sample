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
			p.API.LogError("request body to term error.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// golang struct 'term' to json
		body, err := json.Marshal(term)
		if err != nil {
			p.API.LogError("term to body create error.")
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
			p.API.LogError("http.NewRequest create error.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
			p.API.LogError("http.Client do error.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }

		if resp.StatusCode >= 500 {
			p.API.LogError(fmt.Sprintf("Outside API ERROR StatusCode: %d", resp.StatusCode))
			http.Error(w, "Outside API InternalServerError", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode == 404 {
			p.API.LogWarn(fmt.Sprintf("Similor words not found. StatusCode: %d", resp.StatusCode))
			http.Error(w, "Similor words not found", http.StatusNotFound)
			return
		}

		// parse json response from outside API to golang struct '[]*Article'
		results := make([]*Article, 0)
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			p.API.LogError("response to articles error.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// golang struct '[]*Article' to json
		res, err := json.Marshal(results)
		if err != nil {
			p.API.LogError("articles to json error")
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
