package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchHandlerSuccess(t *testing.T) {
	loadData("data/test_data.txt")
	req := httptest.NewRequest("GET", "/search/{num}", nil)
	req.SetPathValue("num", "100")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(searchValueHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, "2", w.Body.String())
}

func TestSearchHandlerWithTolerationSuccess(t *testing.T) {
	loadData("data/test_data.txt")
	req := httptest.NewRequest("GET", "/search/{num}", nil)
	req.SetPathValue("num", "10010")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(searchValueHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, "4", w.Body.String())
}

func TestSearchHandlerNotFound(t *testing.T) {
	loadData("data/test_data.txt")
	req := httptest.NewRequest("GET", "/search/{num}", nil)
	req.SetPathValue("num", "5000")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(searchValueHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Value not found\n", w.Body.String())
}

func TestHandlebad(t *testing.T) {
	loadData("data/test_data.txt")
	req := httptest.NewRequest("GET", "/search/{num}", nil)
	req.SetPathValue("num", "h")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(searchValueHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoadNumbersFromFile(t *testing.T) {

	loadData("data/test_data.txt")
	assert.Equal(t, len(numberSlice), 5)

	expected := []int{1, 10, 100, 1000, 10000}
	assert.Equal(t, numberSlice, expected)
}

func TestLoadNumbersFromNonExistentFile(t *testing.T) {
	err := loadData("no_file.txt")
	assert.Equal(t, err.Error(), "open no_file.txt: no such file or directory")
}

func TestLoadNumbersFromInvalidData(t *testing.T) {
	err := loadData("data/invalid_data.txt")
	assert.NotEmpty(t, err.Error())
}

func TestFindValue(t *testing.T) {
	numbers := []int{1, 10, 100, 1000, 10000}
	index := findValue(numbers, 100)
	assert.Equal(t, index, 2)
}

func TestFindValueWithTolerance(t *testing.T) {
	numbers := []int{1, 10, 100, 1000, 10000}
	index := findValue(numbers, 109)
	assert.Equal(t, index, 2)
}

func TestFindValueNotFound(t *testing.T) {
	numbers := []int{1, 10, 100, 1000, 10000}
	index := findValue(numbers, 150)
	assert.Equal(t, index, -1)
}
