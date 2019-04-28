package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HealthEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	router().ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, `{ "status": "ok" }`, rec.Body.String())
}
