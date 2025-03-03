package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"concurrent_money_transfer_system/internals/server"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Setup() {
	router = server.SetupRouter()
}

type Request struct {
	URL     string                 `json:"url"`
	Method  string                 `json:"method"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers,omitempty"`
}

type Response struct {
	Status  int                    `json:"status"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers,omitempty"`
}

type TestData struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

func MakeRequestAndValidateResponse(t *testing.T, testData TestData) *httptest.ResponseRecorder {
	// Compare response
	responseBodyMap, recorder := MakeRequestAndGetResponse(t, testData)

	expectedBody := testData.Response.Body
	if !SelectiveEqual(expectedBody, responseBodyMap) {
		t.Fatalf("Expected response body %v, got %v", testData.Response.Body, responseBodyMap)
	}

	// Compare headers
	expectedHeaders := testData.Response.Headers
	for key, value := range expectedHeaders {
		if recorder.Header().Get(key) != value {
			t.Fatalf("Expected header %s to be %s, got %s", key, value, recorder.Header().Get(key))
		}
	}

	return recorder
}

func MakeRequestAndGetResponse(t *testing.T, testData TestData) (map[string]interface{}, *httptest.ResponseRecorder) {
	jsonBody, err := json.Marshal(testData.Request.Body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create a test request
	request := httptest.NewRequest(testData.Request.Method, "/"+testData.Request.URL, bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")

	// Add headers if present
	if testData.Request.Headers != nil {
		for key, value := range testData.Request.Headers {
			request.Header.Set(key, value)
		}
	}

	// Create a response recorder
	recorder := httptest.NewRecorder()

	// Get the handler from your app and serve the request
	router.ServeHTTP(recorder, request)

	// Check status code
	if recorder.Code != testData.Response.Status {
		t.Fatalf("Expected status code %d, got %d. Response body: %s",
			testData.Response.Status, recorder.Code, recorder.Body.String())
	}

	// Parse response body
	var responseBodyMap map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&responseBodyMap); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	return responseBodyMap, recorder
}

func SelectiveEqual(expectedBody map[string]interface{}, actualBody map[string]interface{}) bool {
	for key, expectedValue := range expectedBody {
		actualValue, ok := actualBody[key]
		if !ok {
			return false
		}
		if key == "created_at" || key == "updated_at" || key == "deleted_at" {
			continue
		}
		if expectedValue != actualValue {
			return false
		}
	}
	return len(expectedBody) == len(actualBody)
}

func ReadTestData(path string) map[string]TestData {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to open test data file: %v", err))
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var x interface{}
	json.Unmarshal(byteValue, &x)

	var testData map[string]TestData
	err = json.Unmarshal(byteValue, &testData)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal test data: %v", err))
	}
	return testData
}
