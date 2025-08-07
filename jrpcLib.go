package jrpcLib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Destination struct {
	Client   *http.Client
	Method   string `json:"method"`
	Protocol string `json:"protocol"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
}

type JRPC struct {
	Version string                 `json:"jsonrpc"`
	ID      string                 `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type JRPCResult struct {
	Version string                 `json:"jsonrpc"`
	ID      string                 `json:"id"`
	Result  map[string]interface{} `json:"result"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []byte `json:"data"`
	} `json:"error"`
}

// Wrapper creates a json byte array from the 'JRPC' struct to be used in the 'Call' function
// It checks for an empty param field and if true will remove that from the array
func (jrpc *JRPC) Wrapper() ([]byte, error) {
	if data, err := json.Marshal(jrpc); err != nil {
		return nil, err
	} else {
		if bytes.Contains(data, []byte(",\"params\":{}")) {
			data = bytes.ReplaceAll(data, []byte(",\"params\":{}"), []byte(""))
		}
		return data, nil
	}
}

// Call sends an HTTP request to the destination using JRPC parameters and returns the response or an error.
func (dest *Destination) Call(jrpc *JRPC) (*http.Response, error) {
	if dest.Method != "" && dest.Protocol != "" && dest.IP != "" {
		if reqBody, err := jrpc.Wrapper(); err != nil {
			return nil, err
		} else {
			url := fmt.Sprintf("%s://%s:%d", dest.Protocol, dest.IP, dest.Port)
			if dest.Path != "" {
				url = fmt.Sprintf("%s/%s", url, dest.Path)
			}
			if req, err := http.NewRequest(dest.Method, url, bytes.NewBuffer(reqBody)); err != nil {
				return nil, err
			} else {
				req.Header.Set("Content-Type", "application/json")
				if res, err := dest.Client.Do(req); err != nil {
					return nil, err
				} else {
					return res, nil
				}
			}
		}
	} else {
		return nil, fmt.Errorf("jrpcLib-Call: invalid destination")
	}
}
