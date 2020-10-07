// Package aggregation_plugin a aggregation plugin.
package aggregation_plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type responseWriter struct {
	buffer       bytes.Buffer
	lastModified bool
	wroteHeader  bool

	http.ResponseWriter
}

// Config the plugin configuration.
type Config struct {
	Server string `json:"server,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Aggregation a Aggregation plugin.
type Aggregation struct {
	next   http.Handler
	server string
	name   string
}

// New created a new Aggregation plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Server) == 0 {
		return nil, fmt.Errorf("server cannot be empty")
	}

	return &Aggregation{
		server: config.Server,
		next:   next,
		name:   name,
	}, nil
}

func (a *Aggregation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	// Unmarshal
	var request map[string]interface{}

	err = json.Unmarshal([]byte(b), &request)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	response := make(map[string]*json.RawMessage)

	for key, value := range request {
		str := fmt.Sprintf("%v", value)

		resp, err := http.Get(a.server + "/" + str)

		if err != nil {
			stringRes := "{\"HTTPStatusCode\": \"500\"}"
			var partialResponse *json.RawMessage
			json.Unmarshal([]byte(stringRes), &partialResponse)
			response[key] = partialResponse
			continue
		}
		if resp.StatusCode != 200 {
			stringRes := "{\"HTTPStatusCode\": " + fmt.Sprint(resp.StatusCode) + "}"
			var partialResponse *json.RawMessage
			json.Unmarshal([]byte(stringRes), &partialResponse)
			response[key] = partialResponse
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		var partialResponse *json.RawMessage
		json.Unmarshal(b, &partialResponse)
		response[key] = partialResponse
	}

	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(response)

}
