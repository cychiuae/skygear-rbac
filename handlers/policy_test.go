package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin"
)

func TestPolicyHandler(t *testing.T) {
	e := casbin.NewEnforcer("../model.conf")

	e.AddPolicy("root", "admin", "data1", "write")

	handler := &PolicyHandler{e}
	server := httptest.NewServer(handler)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	q := req.URL.Query()
	q.Add("domain", "root")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", res.StatusCode)
	}
	policy, _ := json.Marshal([]Policy{Policy{"root", "admin", "data1", "write"}})
	expected := string(policy)
	actual, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(actual) {
		t.Errorf("Expected the message '%s'\n", expected)
		t.Errorf("Received '%s'\n", actual)
	}
}
