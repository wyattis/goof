package mock

import (
	"io"
	"net/http"
	"strings"
	"testing"

	mhttp "github.com/wyattis/goof/http"
)

var routes = Routes{
	"GET /": Json(
		map[string]string{"hello": "world"},
		ExpectHeader("Authorization", "Token test-token"),
		ExpectHeader("Accept", "application/json"),
		ExpectContentType("application/json"),
	),
}

func TestGetJson(t *testing.T) {
	s := NewServer(routes)
	defer s.Close()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token test-token")
	req.Header.Set("Accept", "application/json")
	client := mhttp.NewBaseClient(s.URL, s.Client())
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json; got %v", res.Header.Get("Content-Type"))
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	expecting := `{"hello":"world"}`
	if strings.TrimSpace(string(data)) != expecting {
		t.Errorf("expected %s; got %s", expecting, string(data))
	}
}
