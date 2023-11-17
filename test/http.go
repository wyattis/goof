package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// DoJson performs a request to the provided url with the given method while encoding the payload as JSON and decoding
// the response as JSON
func DoJson[Req any, Res any](s *httptest.Server, method string, relUrl string, payload Req, res Res) (err error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}
	url := s.URL + relUrl
	req, err := http.NewRequest(method, url, bytes.NewBuffer(encoded))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := s.Client().Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		err = fmt.Errorf("Unexpected status code %d", r.StatusCode)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&res)
	return
}

// GetJson performs a GET request to the provided url and decodes the response into res.
func GetJson[Res any](s *httptest.Server, relUrl string, res Res) (err error) {
	url := s.URL + relUrl
	r, err := s.Client().Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		err = fmt.Errorf("Unexpected status code %d", r.StatusCode)
		return
	}
	err = json.NewDecoder(r.Body).Decode(res)
	return
}

// PutJson is an alias for DoJson with method PUT
func PutJson[Req any, Res any](s *httptest.Server, relUrl string, payload Req, res Res) (err error) {
	return DoJson(s, http.MethodPut, relUrl, payload, res)
}

// PostJson is an alias for DoJson with method POST
func PostJson[Req any, Res any](s *httptest.Server, relUrl string, payload Req, res Res) (err error) {
	return DoJson[Req, Res](s, http.MethodPost, relUrl, payload, res)
}

// PatchJson is an alias for DoJson with method PATCH
func PatchJson[Req any, Res any](s *httptest.Server, relUrl string, payload Req, res Res) (err error) {
	return DoJson[Req, Res](s, http.MethodPatch, relUrl, payload, res)
}
