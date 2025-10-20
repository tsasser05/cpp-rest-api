package cpprestapitests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type apiTest struct {
	client           *http.Client
	baseURL          string
	lastResponse     *http.Response
	lastResponseBody []byte
	lastCreatedID    int
}

func (a *apiTest) theAPIIsRunning() error {
	return nil
}

func (a *apiTest) iSendAPOSTRequestToWithContactDetails(path, details string) error {
	req, err := http.NewRequest("POST", a.baseURL+path, bytes.NewBuffer([]byte(details)))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	a.lastResponse = resp
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read POST response: %v", err)
	}
	a.lastResponseBody = body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err == nil {
		if id, ok := responseBody["id"].(float64); ok {
			a.lastCreatedID = int(id)
		}
	}
	return nil
}

func (a *apiTest) iSendAGETRequestTo(path string) error {
	path = strings.ReplaceAll(path, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	req, err := http.NewRequest("GET", a.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %v", err)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %v", err)
	}
	a.lastResponse = resp
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read GET response: %v", err)
	}
	a.lastResponseBody = body
	return nil
}

func (a *apiTest) iSendAPUTRequestToWithUpdatedDetails(path, details string) error {
	path = strings.ReplaceAll(path, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	req, err := http.NewRequest("PUT", a.baseURL+path, bytes.NewBuffer([]byte(details)))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PUT request: %v", err)
	}
	a.lastResponse = resp
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read PUT response: %v", err)
	}
	a.lastResponseBody = body
	return nil
}

func (a *apiTest) iSendADELETERequestTo(path string) error {
	path = strings.ReplaceAll(path, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	req, err := http.NewRequest("DELETE", a.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request: %v", err)
	}
	a.lastResponse = resp
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read DELETE response: %v", err)
	}
	a.lastResponseBody = body
	return nil
}

func (a *apiTest) aContactExistsWithID(id string) error {
	contact := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"phone":      "1234567890",
		"email":      "john@example.com",
		"street":     "123 Main St",
		"city":       "Anytown",
		"state":      "CA",
		"zip":        "12345",
	}
	body, err := json.Marshal(contact)
	if err != nil {
		return fmt.Errorf("failed to marshal contact: %v", err)
	}
	req, err := http.NewRequest("POST", a.baseURL+"/records", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create contact, got status %d", resp.StatusCode)
	}
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	a.lastCreatedID = int(responseBody["id"].(float64))
	return nil
}

func (a *apiTest) aContactExistsWithFirstNameAndPhone(firstName, phone string) error {
	contact := map[string]interface{}{
		"first_name": firstName,
		"last_name":  "Doe",
		"phone":      phone,
		"email":      "john@example.com",
		"street":     "123 Main St",
		"city":       "Anytown",
		"state":      "CA",
		"zip":        "12345",
	}
	body, err := json.Marshal(contact)
	if err != nil {
		return fmt.Errorf("failed to marshal contact: %v", err)
	}
	req, err := http.NewRequest("POST", a.baseURL+"/records", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create contact, got status %d", resp.StatusCode)
	}
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	a.lastCreatedID = int(responseBody["id"].(float64))
	return nil
}

func (a *apiTest) theResponseStatusCodeShouldBe(statusCode int) error {
	if a.lastResponse == nil {
		return fmt.Errorf("no response received")
	}
	if a.lastResponse.StatusCode != statusCode {
		return fmt.Errorf("expected status code %d, got %d, response body: %s", statusCode, a.lastResponse.StatusCode, string(a.lastResponseBody))
	}
	return nil
}

func (a *apiTest) theResponseShouldContainTheContactID() error {
	var responseBody map[string]interface{}
	if err := json.Unmarshal(a.lastResponseBody, &responseBody); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(a.lastResponseBody))
	}
	if _, ok := responseBody["id"]; !ok {
		return fmt.Errorf("response does not contain an ID, body: %s", string(a.lastResponseBody))
	}
	return nil
}

func (a *apiTest) theResponseShouldContainTheContactDetails(docString string) error {
	docString = strings.ReplaceAll(docString, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	var expected map[string]interface{}
	if err := json.Unmarshal([]byte(docString), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal expected: %v", err)
	}
	var actual map[string]interface{}
	if err := json.Unmarshal(a.lastResponseBody, &actual); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(a.lastResponseBody))
	}
	return compareJSON(actual, expected)
}

func (a *apiTest) theResponseShouldContainTheUpdatedContactDetails(docString string) error {
	docString = strings.ReplaceAll(docString, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	var expected map[string]interface{}
	if err := json.Unmarshal([]byte(docString), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal expected: %v", err)
	}
	var actual map[string]interface{}
	if err := json.Unmarshal(a.lastResponseBody, &actual); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(a.lastResponseBody))
	}
	return compareJSON(actual, expected)
}

func (a *apiTest) theResponseShouldContainAListWithTheContact(docString string) error {
	var expected []map[string]interface{}
	docString = strings.ReplaceAll(docString, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	if err := json.Unmarshal([]byte(docString), &expected); err != nil {
		return fmt.Errorf("failed to unmarshal expected: %v", err)
	}
	var actual []map[string]interface{}
	if err := json.Unmarshal(a.lastResponseBody, &actual); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v, body: %s", err, string(a.lastResponseBody))
	}
	if len(actual) != len(expected) {
		return fmt.Errorf("expected %d contact(s), got %d, response body: %s", len(expected), len(actual), string(a.lastResponseBody))
	}
	return compareJSON(actual, expected)
}

func (a *apiTest) aSubsequentGETRequestToShouldReturn(path string, statusCode int) error {
	path = strings.ReplaceAll(path, "{lastCreatedID}", fmt.Sprintf("%d", a.lastCreatedID))
	req, err := http.NewRequest("GET", a.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %v", err)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != statusCode {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status code %d, got %d, response body: %s", statusCode, resp.StatusCode, string(body))
	}
	return nil
}

func compareJSON(actual, expected interface{}) error {
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		return fmt.Errorf("failed to marshal actual: %v", err)
	}
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("failed to marshal expected: %v", err)
	}
	if string(actualJSON) != string(expectedJSON) {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expectedJSON), string(actualJSON), false)
		return fmt.Errorf("response mismatch (-expected +actual):\n%s", dmp.DiffPrettyText(diffs))
	}
	return nil
}

func (a *apiTest) InitializeScenario(s *godog.ScenarioContext) {
	s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if a.lastCreatedID != 0 {
			req, _ := http.NewRequest("DELETE", a.baseURL+fmt.Sprintf("/records/%d", a.lastCreatedID), nil)
			resp, _ := a.client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}
		return ctx, nil
	})
	s.Step(`^the API is running$`, a.theAPIIsRunning)
	s.Step(`^I send a POST request to "([^"]*)" with contact details$`, a.iSendAPOSTRequestToWithContactDetails)
	s.Step(`^I send a GET request to "([^"]*)"$`, a.iSendAGETRequestTo)
	s.Step(`^I send a PUT request to "([^"]*)" with updated details$`, a.iSendAPUTRequestToWithUpdatedDetails)
	s.Step(`^I send a DELETE request to "([^"]*)"$`, a.iSendADELETERequestTo)
	s.Step(`^a contact exists with ID \{([^}]*)\}$`, a.aContactExistsWithID)
	s.Step(`^a contact exists with first name "([^"]*)" and phone "([^"]*)"$`, a.aContactExistsWithFirstNameAndPhone)
	s.Step(`^the response status code should be (\d+)$`, a.theResponseStatusCodeShouldBe)
	s.Step(`^the response should contain the contact ID$`, a.theResponseShouldContainTheContactID)
	s.Step(`^the response should contain the contact details$`, a.theResponseShouldContainTheContactDetails)
	s.Step(`^the response should contain a list with the contact$`, a.theResponseShouldContainAListWithTheContact)
	s.Step(`^the response should contain the updated contact details$`, a.theResponseShouldContainTheUpdatedContactDetails)
	s.Step(`^a subsequent GET request to "([^"]*)" should return (\d+)$`, a.aSubsequentGETRequestToShouldReturn)
}

func TestMain(m *testing.M) {
	apiURL := "http://localhost:8080"
	client := &http.Client{}
	ctx := &apiTest{
		client:  client,
		baseURL: apiURL,
	}
	status := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx.InitializeScenario(s)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/contacts.feature"},
		},
	}.Run()
	os.Exit(status)
}

