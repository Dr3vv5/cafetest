package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	totalMoscowCafes := len(cafeList["moscow"])
	totalTulaCafes := len(cafeList["tula"])
	requests := []struct { //Create a slice of structs to hold test cases
		city  string
		count int
		want  int
	}{
		{city: "moscow", count: 0, want: 0},
		{city: "moscow", count: 1, want: 1},
		{city: "moscow", count: 2, want: 2},
		{city: "moscow", count: 100, want: totalMoscowCafes},
		{city: "tula", count: 0, want: 0},
		{city: "tula", count: 1, want: 1},
		{city: "tula", count: 2, want: 2},
		{city: "tula", count: 100, want: totalTulaCafes},
	}
	for _, v := range requests {
		response := httptest.NewRecorder() // Create a new ResponseRecorder for each request
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		body := response.Body.String()
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
		city      string
		search    string
		wantCount int
	}{
		{city: "moscow", search: "фасоль", wantCount: 0},
		{city: "moscow", search: "кофе", wantCount: 2},
		{city: "moscow", search: "вилка", wantCount: 1},
	}

	for _, v := range requests { // Iterate over each test case
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&search="+v.search, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code, "для города %s и поиска %s", v.city, v.search)
		body := response.Body.String()
		rawItems := strings.Split(body, ",") // Split the response body into individual cafe names

		actualItems := make([]string, 0, len(rawItems))
		for _, item := range rawItems { // Trim whitespace and filter out empty strings
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				actualItems = append(actualItems, trimmed)
			}
		}
		assert.Equal(t, v.wantCount, len(actualItems), "неверное количество найденных кафе для поиска '%s'", v.search) // Check if the actual count of cafes matches the expected count
		searchLower := strings.ToLower(v.search)
		for _, cafeName := range actualItems { // Check if each cafe name contains the search term (case-insensitive)
			cafeLower := strings.ToLower(cafeName)
			if !strings.Contains(cafeLower, searchLower) { // If a cafe name does not contain the search term, report an error
				t.Errorf("Кафе '%s' не содержит искомое слово '%s' (поиск в городе %s)", cafeName, v.search, v.city)
			}
		}
	}
}
