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

func (jrpc *JRPC) Wrapper() ([]byte, error) {
	if data, err := json.Marshal(jrpc); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

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
