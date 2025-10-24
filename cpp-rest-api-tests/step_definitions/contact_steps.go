package step_definitions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

type ContactTest struct {
	baseURL       string
	httpClient    *http.Client
	lastResponse  string
	lastStatus    int
	contacts      []map[string]interface{}
	lastID        int
	lastDocString *godog.DocString // Store the last DocString for PUT
}

func (c *ContactTest) initializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the API is running$`, c.theAPIIsRunning)
	ctx.Step(`^the database should be empty$`, c.theDatabaseShouldBeEmpty)
	ctx.Step(`^I send a POST request to "([^"]*)" with contact details:$`, c.iSendAPOSTRequestToWithContactDetails)
	ctx.Step(`^I send a PUT request to "([^"]*)" with updated details:$`, c.iSendAPUTRequestToWithUpdatedDetails)
	ctx.Step(`^I send a GET request to "([^"]*)"$`, c.iSendAGETRequestTo)
	ctx.Step(`^I send a DELETE request to "([^"]*)"$`, c.iSendADELETERequestTo)
	ctx.Step(`^I have created (\d+) contacts?$`, c.iHaveCreatedContacts)
	ctx.Step(`^I have created a contact with ID (\d+)$`, c.iHaveCreatedAContactWithID)
	ctx.Step(`^I have created a contact with phone "([^"]*)"$`, c.iHaveCreatedAContactWithPhone)
	ctx.Step(`^the response status code should be (\d+)$`, c.theResponseStatusCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)"$`, c.theResponseShouldContain)
	ctx.Step(`^the response should contain (\d+) contacts?$`, c.theResponseShouldContainContacts)
	ctx.Step(`^I send a (GET|PUT|DELETE) request to "/records/\{lastCreatedID\}"(.*)$`, c.iSendRequestToLastCreatedID)
}

var test *ContactTest

func InitializeScenario(ctx *godog.ScenarioContext) {
	test = &ContactTest{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	test.initializeScenario(ctx)
}

func (c *ContactTest) theAPIIsRunning() error {
	c.baseURL = "http://localhost:8080"
	return nil
}

func (c *ContactTest) theDatabaseShouldBeEmpty() error {
	resp, _ := c.httpClient.Get(c.baseURL + "/records")
	if resp != nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var contacts []interface{}
		json.Unmarshal(body, &contacts)
		for _, contact := range contacts {
			id := int(contact.(map[string]interface{})["id"].(float64))
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/records/%d", c.baseURL, id), nil)
			c.httpClient.Do(req)
		}
	}
	c.contacts = nil
	c.lastID = 0
	c.lastDocString = nil
	return nil
}

func (c *ContactTest) iSendAPOSTRequestToWithContactDetails(path string, docString *godog.DocString) error {
	var data map[string]interface{}
	json.Unmarshal([]byte(docString.Content), &data)

	jsonData, _ := json.Marshal(data)
	resp, err := c.httpClient.Post(c.baseURL+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.lastResponse = string(body)
	c.lastStatus = resp.StatusCode

	var result map[string]interface{}
	json.Unmarshal([]byte(c.lastResponse), &result)
	c.lastID = int(result["id"].(float64))
	c.contacts = append(c.contacts, data)
	c.lastDocString = docString // Store for PUT
	return nil
}

func (c *ContactTest) iSendAPUTRequestToWithUpdatedDetails(path string, docString *godog.DocString) error {
	var data map[string]interface{}
	json.Unmarshal([]byte(docString.Content), &data)

	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", c.baseURL+path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.lastResponse = string(body)
	c.lastStatus = resp.StatusCode
	c.lastDocString = docString // Store for reuse
	return nil
}

func (c *ContactTest) iSendAGETRequestTo(path string) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	c.lastResponse = string(body)
	c.lastStatus = resp.StatusCode
	return nil
}

func (c *ContactTest) iSendADELETERequestTo(path string) error {
	req, _ := http.NewRequest("DELETE", c.baseURL+path, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.lastStatus = resp.StatusCode
	return nil
}

func (c *ContactTest) iHaveCreatedContacts(count int) error {
	for i := 0; i < count; i++ {
		data := map[string]interface{}{
			"first_name": fmt.Sprintf("User%d", i+1),
			"last_name":  fmt.Sprintf("Last%d", i+1),
			"phone":      fmt.Sprintf("123456789%d", i),
		}
		jsonData, _ := json.Marshal(data)
		resp, err := c.httpClient.Post(c.baseURL+"/records", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		c.lastID = int(result["id"].(float64))
		c.contacts = append(c.contacts, data)
	}
	return nil
}

func (c *ContactTest) iHaveCreatedAContactWithID(id int) error {
	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"phone":      "1234567890",
		"email":      "john@example.com",
	}
	jsonData, _ := json.Marshal(data)
	resp, err := c.httpClient.Post(c.baseURL+"/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	c.lastID = int(result["id"].(float64))
	c.contacts = append(c.contacts, data)
	return nil
}

func (c *ContactTest) iHaveCreatedAContactWithPhone(phone string) error {
	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"phone":      phone,
	}
	jsonData, _ := json.Marshal(data)
	resp, err := c.httpClient.Post(c.baseURL+"/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	c.lastID = int(result["id"].(float64))
	c.contacts = append(c.contacts, data)
	return nil
}

func (c *ContactTest) theResponseStatusCodeShouldBe(status int) error {
	if c.lastStatus != status {
		return fmt.Errorf("expected status %d, got %d", status, c.lastStatus)
	}
	return nil
}

func (c *ContactTest) theResponseShouldContain(text string) error {
	if !strings.Contains(c.lastResponse, text) {
		return fmt.Errorf("expected response to contain %q, got:\n%s", text, c.lastResponse)
	}
	return nil
}

func (c *ContactTest) theResponseShouldContainContacts(count int) error {
	var contacts []interface{}
	err := json.Unmarshal([]byte(c.lastResponse), &contacts)
	if err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	if len(contacts) != count {
		return fmt.Errorf("expected %d contacts, got %d", count, len(contacts))
	}
	return nil
}

func (c *ContactTest) iSendRequestToLastCreatedID(method string, docString *godog.DocString) error {
	path := fmt.Sprintf("/records/%d", c.lastID)
	if method == "GET" {
		return c.iSendAGETRequestTo(path)
	} else if method == "PUT" {
		if docString == nil {
			return fmt.Errorf("no DocString provided for PUT request")
		}
		return c.iSendAPUTRequestToWithUpdatedDetails(path, docString)
	} else if method == "DELETE" {
		return c.iSendADELETERequestTo(path)
	}
	return fmt.Errorf("unsupported method: %s", method)
}
