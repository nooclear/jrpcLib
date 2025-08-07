package jrpcLib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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

type HttpResponse struct {
	Status           string               `json:"status"`
	StatusCode       int                  `json:"status_code"`
	Proto            string               `json:"proto"`
	ProtoMajor       int                  `json:"proto_major"`
	ProtoMinor       int                  `json:"proto_minor"`
	Header           http.Header          `json:"header"`
	Body             []byte               `json:"body"`
	ContentLength    int64                `json:"content_length"`
	TransferEncoding []string             `json:"transfer_encoding"`
	Close            bool                 `json:"close"`
	Uncompressed     bool                 `json:"uncompressed"`
	Trailer          http.Header          `json:"trailer"`
	Request          *http.Request        `json:"request"`
	TLS              *tls.ConnectionState `json:"tls"`
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
func (dest *Destination) Call(jrpc *JRPC) (httpRes HttpResponse, err error) {
	if dest.Method != "" && dest.Protocol != "" && dest.IP != "" {
		if reqBody, err := jrpc.Wrapper(); err != nil {
			return httpRes, err
		} else {
			url := fmt.Sprintf("%s://%s:%d", dest.Protocol, dest.IP, dest.Port)
			if dest.Path != "" {
				url = fmt.Sprintf("%s/%s", url, dest.Path)
			}
			if req, err := http.NewRequest(dest.Method, url, bytes.NewBuffer(reqBody)); err != nil {
				return httpRes, err
			} else {
				req.Header.Set("Content-Type", "application/json")
				if res, err := dest.Client.Do(req); err != nil {
					return httpRes, err
				} else {
					defer func() {
						err = res.Body.Close()
					}()

					httpRes = HttpResponse{
						Status:           res.Status,
						StatusCode:       res.StatusCode,
						Proto:            res.Proto,
						ProtoMajor:       res.ProtoMajor,
						ProtoMinor:       res.ProtoMinor,
						Header:           res.Header,
						ContentLength:    res.ContentLength,
						TransferEncoding: res.TransferEncoding,
						Close:            res.Close,
						Uncompressed:     res.Uncompressed,
						Trailer:          res.Trailer,
						Request:          res.Request,
						TLS:              res.TLS,
					}

					httpRes.Body, err = io.ReadAll(res.Body)
					return httpRes, err
				}
			}
		}
	} else {
		return httpRes, fmt.Errorf("jrpcLib-Call: invalid destination")
	}
}
