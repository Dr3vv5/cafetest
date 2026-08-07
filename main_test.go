package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct { //Create a slice of structs to hold test cases
		city  string
		count int
		want  int
	}{
		{city: "moscow", count: 0, want: 0},
		{city: "moscow", count: 1, want: 1},
		{city: "moscow", count: 2, want: 2},
		{city: "moscow", count: 100, want: 5},
		{city: "tula", count: 0, want: 0},
		{city: "tula", count: 1, want: 1},
		{city: "tula", count: 2, want: 2},
		{city: "tula", count: 100, want: 3},
	}
	for _, v := range requests {
		response := httptest.NewRecorder() // Create a new ResponseRecorder for each request
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		body := response.Body.String()
		if response.Code != http.StatusOK { // Check if the response code is not 200 OK
			t.Fatalf("For city %s and count %d expected status %d but got %d", v.city, v.count, http.StatusOK, response.Code)
		}
		var actualCount int
		if body != "" { // Check if the response body is not empty
			actualCount = len(strings.Split(body, ","))
		} else {
			actualCount = 0
		}
		assert.Equal(t, v.want, actualCount) // Check if the actual count matches the expected count
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		city   string
		search string
		want   []string
	}{
		{city: "moscow", search: "фасоль", want: []string{}},
		{city: "moscow", search: "кофе", want: []string{"Мир кофе", "Кофе и завтраки"}},
		{city: "moscow", search: "вилка", want: []string{"Ложка и вилка"}},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&search="+v.search, nil)
		handler.ServeHTTP(response, req)
		body := response.Body.String()
		if response.Code != http.StatusOK {
			t.Fatalf("For city %s and search %s expected status %d but got %d", v.city, v.search, http.StatusOK, response.Code)
		}
		var actualItems []string
		if body != "" {
			actualItems = strings.Split(strings.TrimSpace(body), ",")
		} else {
			actualItems = []string{}
		}
		assert.Equal(t, v.want, actualItems)
	}
}
