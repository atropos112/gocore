package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	sloghttp "github.com/samber/slog-http"
)

// APIError is an error type that is returned when an API request fails.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return "API request failed with status code " + strconv.Itoa(e.StatusCode) + ": " + e.Message
}

// MakeAPIRequest is a generic function to make an API request. It supports GET, POST, PUT, and DELETE requests.
func MakeAPIRequest(client *http.Client, kind, apiBaseURL, endpoint, token string, request, response interface{}) error {
	l := slog.Default().With("kind", kind, "apiBaseURL", apiBaseURL, "endpoint", endpoint)
	// If response is not nil, check its a pointer (easy dev mistake to make).
	if reflect.ValueOf(response).Kind() != reflect.Ptr {
		return &DeveloperError{"response provided must be a pointer"}
	}

	var jsonData []byte
	var err error

	if request != nil {
		if kind == "GET" || kind == "DELETE" {
			l.Error("GET and DELETE requests do not support request bodies", "kind", kind)
			return &DeveloperError{"GET and DELETE requests do not support request bodies"}
		}

		jsonData, err = json.Marshal(request)
		if err != nil {
			return err
		}
	}

	var req *http.Request

	if request != nil {
		buf := bytes.NewBuffer(jsonData)
		req, err = http.NewRequest(kind, apiBaseURL+endpoint, buf)
	} else {
		req, err = http.NewRequest(kind, apiBaseURL+endpoint, nil)
	}

	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if request != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		l.Error("Failed to unmarshal response", "error", err, "body", string(body))
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Message: string(body)}
	}

	return nil
}

// MakeDeleteRequest is a helper function to make a DELETE request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.
func MakeDeleteRequest(client *http.Client, apiBaseURL, endpoint, token string, response interface{}) error {
	return MakeAPIRequest(client, "DELETE", apiBaseURL, endpoint, token, nil, response)
}

// MakeGetRequest is a helper function to make a GET request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.
func MakeGetRequest(client *http.Client, apiBaseURL, endpoint, token string, response interface{}) error {
	return MakeAPIRequest(client, "GET", apiBaseURL, endpoint, token, nil, response)
}

// MakePostRequest is a helper function to make a POST request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.
func MakePostRequest(client *http.Client, apiBaseURL, endpoint, token string, request, response interface{}) error {
	return MakeAPIRequest(client, "POST", apiBaseURL, endpoint, token, request, response)
}

// MakePutRequest is a helper function to make a PUT request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.
func MakePutRequest(client *http.Client, apiBaseURL, endpoint, token string, request, response interface{}) error {
	return MakeAPIRequest(client, "PUT", apiBaseURL, endpoint, token, request, response)
}

// RunAPIServer attaches logging middleware to the default http server and starts it on the specified port.
func RunAPIServer(port int) error {
	l := slog.Default()

	// mux router
	mux := http.NewServeMux()

	// Include the default mux in your custom mux
	mux.Handle("/", http.DefaultServeMux)

	// Middleware
	l.Info("Setting up middleware for API server")
	handler := sloghttp.Recovery(mux)
	handler = sloghttp.New(l)(handler)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	l.Info("Starting API server", "port", port)
	return http.ListenAndServe(addr, handler)
}
