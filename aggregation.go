// Package aggregation_plugin a aggregation plugin.
package aggregation_plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Header string `json:"header,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Aggregation a Aggregation plugin.
type Aggregation struct {
	next   http.Handler
	server string
	header string
	name   string
}

// New created a new Aggregation plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Server) == 0 {
		return nil, fmt.Errorf("server cannot be empty")
	}
	return &Aggregation{
		server: config.Server,
		header: config.Header,
		next:   next,
		name:   name,
	}, nil
}

func (a *Aggregation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	// Unmarshal
	var request map[string]interface{}

	err := json.Unmarshal([]byte(b), &request)

	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
		return
	}
	response := make(map[string]*json.RawMessage)

	if a.header != "" {
		fmt.Println("Adding header " + a.header)
	}
	for key, value := range request {
		str := fmt.Sprintf("%v", value)
		client := &http.Client{}
		subReq, _ := http.NewRequest("GET", a.server+"/"+str, nil)

		headerValue := req.Header[a.header]
		if len(headerValue) != 0 {
			subReq.Header.Set(a.header, headerValue[0])
		}

		resp, _ := client.Do(subReq)

		if err != nil {
			fmt.Println(err)
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
		b, _ := ioutil.ReadAll(resp.Body)
		var partialResponse *json.RawMessage
		json.Unmarshal(b, &partialResponse)
		response[key] = partialResponse
	}

	rw.Header().Set("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(response)

}
