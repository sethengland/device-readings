package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func dataMapAndRouterSetup() *gin.Engine {
	router := gin.Default()
	dataMap = make(map[string]map[Reading]struct{})
	router.POST("/post-readings", postReadings)
	router.GET("/get-reading/:id", getReading)
	return router
}

func TestPost(t *testing.T) {
	router := dataMapAndRouterSetup()
	ts, _ := time.Parse(time.RFC3339, "2021-09-29T16:08:15+01:00")
	tests := []struct {
		name         string
		input        string
		dataMapState map[string]map[Reading]struct{}
		statuscode   int
	}{
		{
			"basic case",
			`{"id":"adieu", "readings":[{"timestamp":"2021-09-29T16:08:15+01:00", "count":5}]}`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}}},
			http.StatusOK,
		},
		{
			"try to add duplicate reading but datamap state should be unchanged",
			`{"id":"adieu", "readings":[{"timestamp":"2021-09-29T16:08:15+01:00", "count":5}]}`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}}},
			http.StatusOK,
		},
		{
			"case with multiple readings",
			`{"id":"adieu", "readings":[{"timestamp":"2021-09-29T16:08:15+01:00", "count":5}, {"timestamp":"2021-09-29T16:08:15+01:00", "count":2}, {"timestamp":"2021-09-29T16:08:15+01:00", "count":1}]}`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}, Reading{ts, 2}: struct{}{}, Reading{ts, 1}: struct{}{}}},
			http.StatusOK,
		},
		{
			"case with no new readings",
			`{"id":"adieu", "readings":[]}`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}, Reading{ts, 2}: struct{}{}, Reading{ts, 1}: struct{}{}}},
			http.StatusOK,
		},
		{
			"malformed case to simulate incomplete reading",
			`{"id":"adieu", "readings":[{"timestamp":"2021-09`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}, Reading{ts, 2}: struct{}{}, Reading{ts, 1}: struct{}{}}},
			http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest("POST", "/post-readings", strings.NewReader(tc.input))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, tc.statuscode, resp.Code)
		assert.True(t, reflect.DeepEqual(tc.dataMapState, dataMap))

	}
}

func TestGet(t *testing.T) {
	router := dataMapAndRouterSetup()
	ts, _ := time.Parse(time.RFC3339, "2021-09-29T16:08:15+01:00")
	tests := []struct {
		name             string
		inputID          string
		dataMapState     map[string]map[Reading]struct{}
		expectedReadings []Reading
		statuscode       int
	}{
		{
			"id doesn't exist",
			`homerun`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}}},
			[]Reading{{ts, 5}},
			http.StatusNotFound,
		},
		{
			"basic case",
			`adieu`,
			map[string]map[Reading]struct{}{"adieu": {Reading{ts, 5}: struct{}{}}},
			[]Reading{{ts, 5}},
			http.StatusOK,
		},
		{
			"case with multiple readings",
			`wingo`,
			map[string]map[Reading]struct{}{"wingo": {Reading{ts, 3}: struct{}{}, Reading{ts, 4}: struct{}{}, Reading{ts, 5}: struct{}{}}},
			[]Reading{{ts, 5}, {ts, 4}, {ts, 3}},
			http.StatusOK,
		},
		{
			"case with valid id but no current readings",
			`wingo`,
			map[string]map[Reading]struct{}{"wingo": {}},
			[]Reading{},
			http.StatusOK,
		},
	}

	for _, tc := range tests {
		dataMap = tc.dataMapState
		req, _ := http.NewRequest("GET", "/get-reading/"+tc.inputID, nil)
		req.Header.Set("Content-Type", "application/json")

		var respBody []Reading
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, tc.statuscode, resp.Code)
		if tc.statuscode == http.StatusOK {
			err := json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tc.dataMapState, dataMap))
			assert.ElementsMatch(t, tc.expectedReadings, respBody)
		}
	}
}
