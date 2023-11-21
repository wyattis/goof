package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type empty struct{}

func Empty() *empty {
	return &empty{}
}

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
		content, ioErr := io.ReadAll(r.Body)
		if ioErr != nil {
			panic(ioErr)
		}
		err = fmt.Errorf("Unexpected status code %d:\n%s", r.StatusCode, content)
		return
	}

	// Check if we passed Empty() or Empty as the response type
	_, isEmpty := any(res).(*empty)
	_, isEmptyFn := any(res).(func() *empty)
	if isEmpty || isEmptyFn {
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
		content, ioErr := io.ReadAll(r.Body)
		if ioErr != nil {
			panic(ioErr)
		}
		err = fmt.Errorf("Unexpected status code %d:\n%s", r.StatusCode, content)
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

// Delete performs a DELETE request to
func Delete(s *httptest.Server, relUrl string) (err error) {
	url := s.URL + relUrl
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	r, err := s.Client().Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		err = fmt.Errorf("Unexpected status code %d", r.StatusCode)
		return
	}
	return
}

// PostToJson performs a POST request to the provided url without a payload and decodes the response into res.
func PostToJson[Res any](s *httptest.Server, relUrl string, res Res) (err error) {
	url := s.URL + relUrl
	r, err := s.Client().Post(url, "application/json", nil)
	if err != nil {
		return
	}
	defer r.Body.Close()
	if r.StatusCode > 299 {
		content, ioErr := io.ReadAll(r.Body)
		if ioErr != nil {
			panic(ioErr)
		}
		err = fmt.Errorf("Unexpected status code %d:\n%s", r.StatusCode, content)
		return
	}
	err = json.NewDecoder(r.Body).Decode(res)
	return
}
