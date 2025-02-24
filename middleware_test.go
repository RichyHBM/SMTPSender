package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func OkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestEmptyEnsureHeaderMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	reqRecorder := httptest.NewRecorder()

	ensureHeaderMiddleware := MakeEnsureHeaderMiddleware("", "")
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusOK, reqRecorder.Code)
}

func TestRequestEmptyHeaderEnsureHeaderMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	reqRecorder := httptest.NewRecorder()

	ensureHeaderMiddleware := MakeEnsureHeaderMiddleware("foo", "")
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusForbidden, reqRecorder.Code)

	req.Header.Add("foo", "")
	reqRecorder = httptest.NewRecorder()
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusForbidden, reqRecorder.Code)
}

func TestRequestHeaderEnsureHeaderMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Add("foo", "bar")

	reqRecorder := httptest.NewRecorder()

	ensureHeaderMiddleware := MakeEnsureHeaderMiddleware("foo", "")
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusOK, reqRecorder.Code)
}

func TestRequestWrongHeaderEnsureHeaderMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Add("foo", "bar")

	reqRecorder := httptest.NewRecorder()

	ensureHeaderMiddleware := MakeEnsureHeaderMiddleware("foo", "boo")
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusForbidden, reqRecorder.Code)
}

func TestRequestRightHeaderEnsureHeaderMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Add("foo", "bar")

	reqRecorder := httptest.NewRecorder()

	ensureHeaderMiddleware := MakeEnsureHeaderMiddleware("foo", "bar")
	ensureHeaderMiddleware.ServeHTTP(reqRecorder, req, OkHandler)

	assert.Equal(t, http.StatusOK, reqRecorder.Code)
}
