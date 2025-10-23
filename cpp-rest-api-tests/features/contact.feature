Feature: Contact Management API
  As a developer
  I want to manage contacts via REST API
  So that I can create, read, update, and delete contact records

  Background:
    Given the API server is running on "http://localhost:8080"

  Scenario: Create a new contact
    When I POST a contact with:
      | first_name | John   |
      | last_name  | Doe    |
      | phone      | 1234567890 |
      | email      | john@example.com |
    Then the response status should be 201
    And the response should contain "John"
    And the response should contain "Doe"

  Scenario: Retrieve all contacts
    Given I have created 2 contacts
    When I GET "/records"
    Then the response status should be 200
    And the response should contain 2 contacts

  Scenario: Retrieve contact by ID
    Given I have created a contact with ID 1
    When I GET "/records/1"
    Then the response status should be 200
    And the response should contain "John"

  Scenario: Update existing contact
    Given I have created a contact with ID 1
    When I PUT contact ID 1 with first_name "Jane"
    Then the response status should be 200
    And the response should contain "Jane"

  Scenario: Delete existing contact
    Given I have created a contact with ID 1
    When I DELETE "/records/1"
    Then the response status should be 204

  Scenario: Query contacts by first name
    Given I have created 2 contacts with first_name "John"
    When I GET "/records?first_name=John"
    Then the response status should be 200
    And the response should contain 2 contacts

  Scenario: Query by phone area code
    Given I have created a contact with phone "1234567890"
    When I GET "/records?phone=123"
    Then the response status should be 200
    And the response should contain 1 contact
