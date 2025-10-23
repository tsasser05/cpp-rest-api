package step_definitions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages/go/v21"
)

type ContactTest struct {
	baseURL    string
	httpClient *http.Client
	contacts   []map[string]interface{}
}

func (c *ContactTest) initializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the API server is running on "([^"]*)"$`, c.theAPIServerIsRunningOn)
	ctx.Step(`^I POST a contact with:$`, c.iPOSTAContactWith)
	ctx.Step(`^I have created (\d+) contacts? with first_name "([^"]*)"$`, c.iHaveCreatedContactsWithFirstName)
	ctx.Step(`^I have created (\d+) contacts?$`, c.iHaveCreatedContacts)
	ctx.Step(`^I have created a contact with ID (\d+)$`, c.iHaveCreatedAContactWithID)
	ctx.Step(`^I GET "([^"]*)"$`, c.iGET)
	ctx.Step(`^I PUT contact ID (\d+) with first_name "([^"]*)"$`, c.iPUTContactIDWithFirstName)
	ctx.Step(`^I DELETE "([^"]*)"$`, c.iDELETE)
	ctx.Step(`^the response status should be (\d+)$`, c.theResponseStatusShouldBe)
	ctx.Step(`^the response should contain "([^"]*)"$`, c.theResponseShouldContain)
	ctx.Step(`^the response should contain (\d+) contacts?$`, c.theResponseShouldContainContacts)
}

var test *ContactTest

func InitializeScenario(ctx *godog.ScenarioContext) {
	test = &ContactTest{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	test.initializeScenario(ctx)
}

func (c *ContactTest) theAPIServerIsRunningOn(url string) error {
	c.baseURL = url
	return nil
}

func (c *ContactTest) iPOSTAContactWith(table *godog.Table) error {
	data := make(map[string]interface{})
	
	for _, row := range table.Rows[1:] {
		key := strings.TrimSpace(row.Cells[0].Value)
		value := strings.TrimSpace(row.Cells[1].Value)
		data[key] = value
	}
	
	jsonData, _ := json.Marshal(data)
	resp, err := c.httpClient.Post(c.baseURL+"/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	c.contacts = append(c.contacts, data)
	return nil
}

func (c *ContactTest) iHaveCreatedContacts(count int, firstName string) error {
	for i := 0; i < count; i++ {
		data := map[string]interface{}{
			"first_name": firstName,
			"last_name":  fmt.Sprintf("User%d", i),
			"phone":      fmt.Sprintf("123456789%d", i),
		}
		jsonData, _ := json.Marshal(data)
		resp, _ := c.httpClient.Post(c.baseURL+"/records", "application/json", bytes.NewBuffer(jsonData))
		resp.Body.Close()
		c.contacts = append(c.contacts, data)
	}
	return nil
}

func (c *ContactTest) iHaveCreatedContactsWithFirstName(count int, firstName string) error {
	return c.iHaveCreatedContacts(count, firstName)
}

func (c *ContactTest) iHaveCreatedAContactWithID(id int) error {
	return c.iHaveCreatedContacts(1, "John")
}

func (c *ContactTest) iGET(path string) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return err
	}
	
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	
	// Store last response for verification
	c.lastResponse = string(body)
	c.lastStatus = resp.StatusCode
	return nil
}

func (c *ContactTest) iPUTContactIDWithFirstName(id int, firstName string) error {
	data := map[string]interface{}{"first_name": firstName}
	jsonData, _ := json.Marshal(data)
	
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/records/%d", c.baseURL, id), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	c.lastStatus = resp.StatusCode
	return nil
}

func (c *ContactTest) iDELETE(path string) error {
	req, _ := http.NewRequest("DELETE", c.baseURL+path, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	c.lastStatus = resp.StatusCode
	return nil
}

func (c *ContactTest) theResponseStatusShouldBe(status int) error {
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
	json.Unmarshal([]byte(c.lastResponse), &contacts)
	
	if len(contacts) != count {
		return fmt.Errorf("expected %d contacts, got %d", count, len(contacts))
	}
	return nil
}

// Response storage
var (
	lastResponse string
	lastStatus   int
)

func (c *ContactTest) setLastResponse(resp string, status int) {
	lastResponse = resp
	lastStatus = status
}

